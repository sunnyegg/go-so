package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	mockdb "github.com/sunnyegg/go-so/db/mock"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/util"
)

func TestGetUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		userID        string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.UserID,
			buildStubs: func(store *mockdb.MockStore) {
				userID, err := util.ParseStringToInt64(user.UserID)
				require.NoError(t, err)
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(userID)).Times(1).Return(db.GetUserRow{
					UserLogin:       user.UserLogin,
					UserName:        user.UserName,
					ProfileImageUrl: user.ProfileImageUrl,
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: user.UserID,
			buildStubs: func(store *mockdb.MockStore) {
				userID, err := util.ParseStringToInt64(user.UserID)
				require.NoError(t, err)
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(userID)).Times(1).Return(db.GetUserRow{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			userID: user.UserID,
			buildStubs: func(store *mockdb.MockStore) {
				userID, err := util.ParseStringToInt64(user.UserID)
				require.NoError(t, err)
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(userID)).Times(1).Return(db.GetUserRow{}, errors.New("some error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			userID:     "abc",
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/users/" + tc.userID
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser() loginUserRequest {
	return loginUserRequest{
		UserID:          util.ParseIntToString(int(util.RandomInt(1, 1000))),
		UserLogin:       util.RandomString(10),
		UserName:        util.RandomString(10),
		ProfileImageUrl: util.RandomString(10),
	}
}

func requireBodyUser(t *testing.T, body *bytes.Buffer, user loginUserRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.GetUserRow
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.UserLogin, gotUser.UserLogin)
	require.Equal(t, user.UserName, gotUser.UserName)
	require.Equal(t, user.ProfileImageUrl, gotUser.ProfileImageUrl)
}
