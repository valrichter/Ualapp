package db_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

// Deletes all users from test database
func cleanDB() {
	err := testStore.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// Creates a random user in test database
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

func TestGetUserById(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)

	userFromDB, err := testStore.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userFromDB)
	require.Equal(t, user.ID, userFromDB.ID)
	require.Equal(t, user.Email, userFromDB.Email)
	require.Equal(t, user.HashedPassword, userFromDB.HashedPassword)
	require.WithinDuration(t, user.CreatedAt, userFromDB.CreatedAt, 2*time.Second)
	require.Equal(t, user.Username, userFromDB.Username)
}

func TestGetUserByEmail(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)

	userFromDB, err := testStore.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, userFromDB)
	require.Equal(t, user.ID, userFromDB.ID)
	require.Equal(t, user.Email, userFromDB.Email)
	require.Equal(t, user.HashedPassword, userFromDB.HashedPassword)
	require.WithinDuration(t, user.CreatedAt, userFromDB.CreatedAt, 2*time.Second)
	require.Equal(t, user.Username, userFromDB.Username)
}

func TestListUser(t *testing.T) {
	defer cleanDB()

	usersAmount := util.RandomInt(1, 30)
	var wg sync.WaitGroup
	for i := 0; i < usersAmount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomUser(t)
		}()
	}
	wg.Wait()

	arg := db.ListUsersParams{
		Limit:  int32(usersAmount),
		Offset: int32(util.RandomInt(0, usersAmount-1)),
	}

	users, err := testStore.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, arg.Limit-arg.Offset, int32(len(users)))
}

func TestUpdateUserPassword(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)
	newPassword := util.RandomPassword(util.RandomInt(6, 20))

	newPasswordHashed, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	require.NotEmpty(t, newPasswordHashed)

	arg := db.UpdateUserPasswordParams{
		HashedPassword: newPasswordHashed,
		UpdatedAt:      time.Now(),
		ID:             user.ID,
	}

	updatedUser, err := testStore.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.HashedPassword, updatedUser.HashedPassword)
	require.WithinDuration(t, arg.UpdatedAt, time.Now(), 2*time.Second)
	require.Equal(t, arg.ID, updatedUser.ID)
	require.Equal(t, user.ID, updatedUser.ID)
	require.Equal(t, user.Email, updatedUser.Email)
	require.NotEqual(t, user.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, user.UpdatedAt, updatedUser.UpdatedAt)
}

func TestUpdateUsername(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)

	arg := db.UpdateUsernameParams{
		Username: pgtype.Text{String: util.RandomUsername(), Valid: true},
		ID:       user.ID,
	}

	updatedUser, err := testStore.UpdateUsername(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.Username, updatedUser.Username)
	require.Equal(t, arg.ID, updatedUser.ID)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, user.UpdatedAt, updatedUser.UpdatedAt)
}

func TestDeleteUser(t *testing.T) {
	defer cleanDB()

	user := createRandomUser(t)

	err := testStore.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	userFromDB, err := testStore.GetUserById(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, userFromDB)
}
