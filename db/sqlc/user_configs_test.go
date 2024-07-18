package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomUserConfig(t *testing.T, configType ConfigTypes, value string, userID *int64) UserConfig {
	var user User
	if userID == nil {
		user = createRandomUser(t)
	} else {
		user.ID = *userID
	}

	arg := CreateUserConfigParams{
		UserID:     user.ID,
		ConfigType: configType,
		Value:      value,
	}

	userConfig, err := testStore.CreateUserConfig(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userConfig)

	require.Equal(t, arg.ConfigType, userConfig.ConfigType)
	require.Equal(t, arg.Value, userConfig.Value)

	return userConfig
}

func TestCreateUserConfig(t *testing.T) {
	createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
}

func TestGetUserConfig(t *testing.T) {
	userConfig1 := createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
	userConfig2, err := testStore.GetUserConfig(context.Background(), GetUserConfigParams{
		UserID:     userConfig1.UserID,
		ConfigType: ConfigTypesAutoShoutoutActivation,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userConfig2)

	require.Equal(t, userConfig1.Value, userConfig2.Value)
}

func TestUpdateUserConfig(t *testing.T) {
	userConfig1 := createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
	arg := UpdateUserConfigParams{
		ID:    userConfig1.ID,
		Value: "false",
	}
	userConfig2, err := testStore.UpdateUserConfig(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userConfig2)

	require.Equal(t, arg.ID, userConfig2.ID)
	require.Equal(t, arg.Value, userConfig2.Value)
	require.NotZero(t, userConfig2.UpdatedAt)
}

func TestDeleteUserConfig(t *testing.T) {
	userConfig1 := createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
	err := testStore.DeleteUserConfig(context.Background(), userConfig1.ID)
	require.NoError(t, err)

	userConfig2, err := testStore.GetUserConfig(context.Background(), GetUserConfigParams{
		UserID:     userConfig1.UserID,
		ConfigType: ConfigTypesAutoShoutoutActivation,
	})
	require.Error(t, err)
	require.Equal(t, err, pgx.ErrNoRows)
	require.Empty(t, userConfig2)
}
