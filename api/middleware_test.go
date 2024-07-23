package api

// func addAuth(
// 	t *testing.T,
// 	request *http.Request,
// 	tokenMaker token.Maker,
// 	authorizationType string,
// 	userID int64,
// 	duration time.Duration,
// ) {
// 	token, _, err := tokenMaker.MakeToken(userID, duration)
// 	require.NoError(t, err)

// 	authHeader := fmt.Sprintf("%s %s", authorizationType, token)
// 	request.Header.Set(authorizationHeaderKey, authHeader)
// }

// func TestAuthMiddleware(t *testing.T) {
// 	testCases := []struct {
// 		name          string
// 		setupPath     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore, tokenMaker token.Maker, request *http.Request)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuth(t, request, tokenMaker, authorizationPrefixKey, 1, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker, request *http.Request) {
// 				token := request.Header.Get(authorizationHeaderKey)
// 				require.NotEmpty(t, token)

// 				accessToken := strings.Fields(token)[1]

// 				// get payload from token
// 				payload, err := tokenMaker.VerifyToken(accessToken)
// 				require.NoError(t, err)

// 				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(db.GetSessionParams{
// 					ID:     util.UUIDToUUID(payload.ID),
// 					UserID: payload.UserID,
// 				})).Times(1).Return(db.Session{
// 					ID:                   util.UUIDToUUID(payload.ID),
// 					UserID:               payload.UserID,
// 					RefreshToken:         util.RandomString(32),
// 					UserAgent:            util.RandomString(10),
// 					ClientIp:             util.RandomString(10),
// 					IsBlocked:            false,
// 					ExpiresAt:            util.StringToTimestamp(time.Now().Add(time.Minute).Format(time.RFC3339)),
// 					CreatedAt:            util.StringToTimestamp(time.Now().Format(time.RFC3339)),
// 					EncryptedTwitchToken: util.RandomString(10),
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name:      "Unauthorized",
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(db.GetSessionParams{
// 					ID:     session.ID,
// 					UserID: session.UserID,
// 				})).Times(1).Return(db.Session{
// 					ID:        session.ID,
// 					UserID:    session.UserID,
// 					CreatedAt: session.CreatedAt,
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidAuthorizationType",
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuth(t, request, tokenMaker, "invalid-type", 1, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(db.GetSessionParams{
// 					ID:     session.ID,
// 					UserID: session.UserID,
// 				})).Times(1).Return(db.Session{
// 					ID:        session.ID,
// 					UserID:    session.UserID,
// 					CreatedAt: session.CreatedAt,
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidAuthorizationHeader",
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuth(t, request, tokenMaker, "", 1, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(db.GetSessionParams{
// 					ID:     session.ID,
// 					UserID: session.UserID,
// 				})).Times(1).Return(db.Session{
// 					ID:        session.ID,
// 					UserID:    session.UserID,
// 					CreatedAt: session.CreatedAt,
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "ExpiredToken",
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuth(t, request, tokenMaker, authorizationPrefixKey, 1, -time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(db.GetSessionParams{
// 					ID:     session.ID,
// 					UserID: session.UserID,
// 				})).Times(1).Return(db.Session{
// 					ID:        session.ID,
// 					UserID:    session.UserID,
// 					CreatedAt: session.CreatedAt,
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			server := newTestServer(t, nil)
// 			authPath := "/auth"

// 			server.router.GET(authPath, authMiddleware(server), func(ctx *gin.Context) {
// 				ctx.JSON(http.StatusOK, gin.H{})
// 			})

// 			recorder := httptest.NewRecorder()
// 			request, err := http.NewRequest(http.MethodGet, authPath, nil)
// 			require.NoError(t, err)

// 			tc.setupPath(t, request, server.tokenMaker)
// 			server.router.ServeHTTP(recorder, request)

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store, server.tokenMaker, request)

// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }
