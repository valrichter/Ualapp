package db_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

// cleanDB deletes all users from database
func cleanDB() {
	err := testStore.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// createRandomUser creates a random user of database for tests
func createRandomUser(t *testing.T) db.User {

	password := util.RandomPassword(util.RandomInt(6, 20))
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := db.CreateUserParams{
		ID:             uuid.New(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.CreatedAt)
	require.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)
	require.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestCreateUser tests the CreateUser function of database
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUserById(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

// TestUpdateUser tests the UpdateUserPassword function of database
func TestUpdateUserPassword(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)
	newPassword := util.RandomPassword(util.RandomInt(6, 20))
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := db.UpdateUserPasswordParams{
		HashedPassword: newHashedPassword,
		ID:             user.ID,
		UpdatedAt:      time.Now(),
	}

	updatedUser, err := testStore.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, arg.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, arg.ID, updatedUser.ID)
	require.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)
}

// TestGetUserbyID tests the GetUserById function of database
func TestGetUserbyID(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)

	newUser, err := testStore.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.HashedPassword, newUser.HashedPassword)
	require.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)
	require.WithinDuration(t, user.CreatedAt, newUser.CreatedAt, time.Second)

}

// TestGetUserByEmail tests the GetUserByEmail function of database
func TestGetUserByEmail(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)

	newUser, err := testStore.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.HashedPassword, newUser.HashedPassword)
	require.WithinDuration(t, user.CreatedAt, newUser.CreatedAt, time.Second)
}

// TestDeleteUser tests the DeleteUser function of database
func TestDeleteUser(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)

	err := testStore.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	newUser, err := testStore.GetUserById(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, newUser)
}

// TestListUser tests the ListUsers function
func TestListUser(t *testing.T) {
	defer cleanDB()

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomUser(t)
		}()
	}
	wg.Wait()

	arg := db.ListUsersParams{
		Limit:  30,
		Offset: 0,
	}

	users, err := testStore.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 30, len(users))

}
