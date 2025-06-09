package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/matodrobec/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword("secret")
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testStore.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		FullName:      pgtype.Text{String: newFullName, Valid: true},
		WhereUsername: oldUser.Username,
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.Email, oldUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Email:         pgtype.Text{String: newEmail, Valid: true},
		WhereUsername: oldUser.Username,
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
}

func TestUpdateUserOnlyHashPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newHashedPassword, err := util.HashPassword(util.RandomString(20))
	require.NoError(t, err)
	require.NotEmpty(t, newHashedPassword)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: pgtype.Text{String: newHashedPassword, Valid: true},
		WhereUsername:  oldUser.Username,
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
}

func TestUpdateUserOnlyAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()
	newHashedPassword, err := util.HashPassword(util.RandomString(20))
	require.NoError(t, err)
	require.NotEmpty(t, newHashedPassword)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: pgtype.Text{String: newHashedPassword, Valid: true},
		Email:          pgtype.Text{String: newEmail, Valid: true},
		FullName:       pgtype.Text{String: newFullName, Valid: true},
		WhereUsername:  oldUser.Username,
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.FullName, newFullName)

	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.NotEqual(t, updatedUser.Email, oldUser.Email)
	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
}
