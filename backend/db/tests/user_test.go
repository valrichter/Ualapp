package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

// createRandomUser creates a random user of database for tests
func createRandomUser(t *testing.T) db.User {
	password := util.RandomPassword(util.RandomInt(6, 20))
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQuery.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.CreatedAt)
	require.WithinDuration(t, user.CreatedAt.Time, time.Now(), 2*time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, time.Now(), 2*time.Second)

	return user
}

// TestCreateUser tests the CreateUser function of database
func TestCreateUser(t *testing.T) {
	user1 := createRandomUser(t)
	arg := db.CreateUserParams{
		Email:          user1.Email,
		HashedPassword: user1.HashedPassword,
	}

	user2, err := testQuery.CreateUser(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, user2)
}

// TestUpdateUser tests the UpdateUserPassword function of database
func TestUpdateUserPassword(t *testing.T) {
	user := createRandomUser(t)
	newPassword := util.RandomPassword(util.RandomInt(6, 20))
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := db.UpdateUserPasswordParams{
		HashedPassword: newHashedPassword,
		ID:             user.ID,
		UpdatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	updatedUser, err := testQuery.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, arg.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, arg.ID, updatedUser.ID)
	require.WithinDuration(t, user.UpdatedAt.Time, time.Now(), 2*time.Second)
}

// TestGetUserbyID tests the GetUserById function of database
func TestGetUserbyID(t *testing.T) {
	user := createRandomUser(t)

	newUser, err := testQuery.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.HashedPassword, newUser.HashedPassword)
	require.WithinDuration(t, user.UpdatedAt.Time, time.Now(), 2*time.Second)
	require.WithinDuration(t, user.CreatedAt.Time, newUser.CreatedAt.Time, time.Second)

}

// TestGetUserByEmail tests the GetUserByEmail function of database
func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	newUser, err := testQuery.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.HashedPassword, newUser.HashedPassword)
	require.WithinDuration(t, user.CreatedAt.Time, newUser.CreatedAt.Time, time.Second)
}

// TestDeleteUser tests the DeleteUser function of database
func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQuery.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	newUser, err := testQuery.GetUserById(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, newUser)
}

// TestListUser tests the ListUsers function
func TestListUser(t *testing.T) {
	for i := 0; i < 30; i++ {
		createRandomUser(t)
	}

	arg := db.ListUsersParams{
		Limit:  0,
		Offset: 30,
	}

	users, err := testQuery.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 30, len(users))

}
