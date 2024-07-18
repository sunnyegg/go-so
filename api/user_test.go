package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	mockdb "github.com/sunnyegg/go-so/db/mock"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/util"
)

func TestGetUserAPI(t *testing.T) {
	user := randomUser()
	getUserRow := db.GetUserRow{
		ID:              user.ID,
		UserLogin:       user.UserLogin,
		UserName:        user.UserName,
		ProfileImageUrl: user.ProfileImageUrl,
	}

	testCases := []struct {
		name          string
		setupPath     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationPrefixKey, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.GetUserRow{
					UserLogin:       user.UserLogin,
					UserName:        user.UserName,
					ProfileImageUrl: user.ProfileImageUrl,
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyUser(t, recorder.Body, getUserRow)
			},
		},
		{
			name: "NotFound",
			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationPrefixKey, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.GetUserRow{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationPrefixKey, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.GetUserRow{}, errors.New("some error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			url := "/users"
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			tc.setupPath(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser() db.User {
	return db.User{
		ID:              util.RandomInt(1, 1000),
		UserID:          util.RandomUserID(),
		UserLogin:       util.RandomString(10),
		UserName:        util.RandomString(10),
		ProfileImageUrl: util.RandomString(10),
		CreatedAt:       util.StringToTimestamp(time.Now().Format(time.RFC3339)),
		UpdatedAt:       util.StringToTimestamp(time.Now().Format(time.RFC3339)),
	}
}

func requireBodyUser(t *testing.T, body *bytes.Buffer, user db.GetUserRow) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.GetUserRow
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.UserLogin, gotUser.UserLogin)
	require.Equal(t, user.UserName, gotUser.UserName)
	require.Equal(t, user.ProfileImageUrl, gotUser.ProfileImageUrl)
}
