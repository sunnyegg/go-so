package db

import (
	"context"
	"testing"
	"time"

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

	userConfig, err := testQueries.CreateUserConfig(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userConfig)

	require.Equal(t, arg.UserID, userConfig.UserID)
	require.Equal(t, arg.ConfigType, userConfig.ConfigType)
	require.Equal(t, arg.Value, userConfig.Value)

	require.NotZero(t, userConfig.ID)
	require.NotZero(t, userConfig.CreatedAt)

	return userConfig
}

func TestCreateUserConfig(t *testing.T) {
	createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
}

func TestGetUserConfig(t *testing.T) {
	userConfig1 := createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
	userConfig2, err := testQueries.GetUserConfig(context.Background(), userConfig1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userConfig2)

	require.Equal(t, userConfig1.UserID, userConfig2.UserID)
	require.Equal(t, userConfig1.ConfigType, userConfig2.ConfigType)
	require.Equal(t, userConfig1.Value, userConfig2.Value)
	require.WithinDuration(t, userConfig1.CreatedAt.Time, userConfig2.CreatedAt.Time, time.Second)
}

func TestUpdateUserConfig(t *testing.T) {
	userConfig1 := createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
	arg := UpdateUserConfigParams{
		ID:    userConfig1.ID,
		Value: "false",
	}
	userConfig2, err := testQueries.UpdateUserConfig(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userConfig2)

	require.Equal(t, arg.ID, userConfig2.ID)
	require.Equal(t, arg.Value, userConfig2.Value)
	require.NotZero(t, userConfig2.UpdatedAt)
}

func TestDeleteUserConfig(t *testing.T) {
	userConfig1 := createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
	err := testQueries.DeleteUserConfig(context.Background(), userConfig1.ID)
	require.NoError(t, err)

	userConfig2, err := testQueries.GetUserConfig(context.Background(), userConfig1.ID)
	require.Error(t, err)
	require.Equal(t, err, pgx.ErrNoRows)
	require.Empty(t, userConfig2)
}

func TestListUserConfigs(t *testing.T) {
	userConfig1 := createRandomUserConfig(t, ConfigTypesAutoShoutoutActivation, "true", nil)
	createRandomUserConfig(t, ConfigTypesAutoShoutoutDelay, "5", &userConfig1.UserID)
	createRandomUserConfig(t, ConfigTypesBlacklist, "sunnyegg21,sunnyeggbot", &userConfig1.UserID)
	createRandomUserConfig(t, ConfigTypesTimerCardDuration, "5", &userConfig1.UserID)
	createRandomUserConfig(t, ConfigTypesTimerCardSoActivation, "false", &userConfig1.UserID)

	arg := ListUserConfigsParams{
		Limit:  5,
		Offset: 0,
	}

	userConfigs, err := testQueries.ListUserConfigs(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, userConfigs, 5)

	for _, userConfig := range userConfigs {
		require.NotEmpty(t, userConfig)
	}
}
