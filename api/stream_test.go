package api

// func TestCreateStreamAPI(t *testing.T) {
// 	user := randomUser()
// 	createStreamParams := randomStream(user.ID)
// 	arg := createStreamRequest{
// 		Title:     createStreamParams.Title,
// 		GameName:  createStreamParams.GameName,
// 		StartedAt: createStreamParams.StartedAt.Time.Format(time.RFC3339),
// 	}
// 	rsp := createStreamResponse{
// 		Title:     createStreamParams.Title,
// 		GameName:  createStreamParams.GameName,
// 		StartedAt: createStreamParams.StartedAt.Time.Format(time.RFC3339),
// 	}

// 	testCases := []struct {
// 		name          string
// 		body          createStreamRequest
// 		setupPath     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: arg,
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuth(t, request, tokenMaker, authorizationPrefixKey, user.ID, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().CreateStream(gomock.Any(), gomock.Eq(db.CreateStreamParams{
// 					UserID:    user.ID,
// 					Title:     arg.Title,
// 					GameName:  arg.GameName,
// 					StartedAt: util.StringToTimestamp(arg.StartedAt),
// 					CreatedBy: user.ID,
// 				})).Times(1).Return(db.Stream{
// 					ID:        createStreamParams.ID,
// 					UserID:    createStreamParams.UserID,
// 					Title:     createStreamParams.Title,
// 					GameName:  createStreamParams.GameName,
// 					StartedAt: createStreamParams.StartedAt,
// 					CreatedBy: createStreamParams.CreatedBy,
// 					CreatedAt: createStreamParams.CreatedAt,
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyStream(t, recorder.Body, rsp)
// 			},
// 		},
// 		{
// 			name: "BadRequest",
// 			body: createStreamRequest{},
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuth(t, request, tokenMaker, authorizationPrefixKey, user.ID, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalServerError",
// 			body: arg,
// 			setupPath: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuth(t, request, tokenMaker, authorizationPrefixKey, user.ID, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().CreateStream(gomock.Any(), gomock.Eq(db.CreateStreamParams{
// 					UserID:    user.ID,
// 					Title:     createStreamParams.Title,
// 					GameName:  createStreamParams.GameName,
// 					StartedAt: util.StringToTimestamp(createStreamParams.StartedAt.Time.Format(time.RFC3339)),
// 					CreatedBy: user.ID,
// 				})).Times(1).Return(db.Stream{}, errors.New("error"))
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			url := "/streams"
// 			body, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			req, err := http.NewRequest("POST", url, bytes.NewReader(body))
// 			require.NoError(t, err)

// 			tc.setupPath(t, req, server.tokenMaker)

// 			server.router.ServeHTTP(recorder, req)
// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }

// func TestGetStreamAPI(t *testing.T) {
// }

// func TestListStreamAPI(t *testing.T) {
// }

// func TestGetStreamAttendanceMemberAPI(t *testing.T) {
// }

// func randomStream(userID int64) db.Stream {
// 	return db.Stream{
// 		ID:        util.RandomInt(1, 1000),
// 		UserID:    userID,
// 		Title:     util.RandomString(10),
// 		GameName:  util.RandomString(10),
// 		StartedAt: util.StringToTimestamp(time.Now().Format(time.RFC3339)),
// 		CreatedBy: util.RandomInt(1, 1000),
// 		CreatedAt: util.StringToTimestamp(time.Now().Format(time.RFC3339)),
// 	}
// }

// func requireBodyStream(t *testing.T, body *bytes.Buffer, stream createStreamResponse) {
// 	data, err := io.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotStream createStreamResponse
// 	err = json.Unmarshal(data, &gotStream)
// 	require.NoError(t, err)

// 	require.Equal(t, stream.Title, gotStream.Title)
// 	require.Equal(t, stream.GameName, gotStream.GameName)
// 	require.NotZero(t, gotStream.StartedAt)
// }
