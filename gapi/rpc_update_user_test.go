package gapi

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/matodrobec/simplebank/db/mock"
	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/pb"
	"github.com/matodrobec/simplebank/token"
	"github.com/matodrobec/simplebank/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// func randomUser(t *testing.T) (user db.User, password string) {
// 	password = util.RandomString(8)
// 	hashedPass, err := util.HashPassword(password)
// 	require.NoError(t, err)

// 	user = db.User{
// 		Username:       util.RandomOwner(),
// 		HashedPassword: hashedPass,
// 		FullName:       util.RandomOwner(),
// 		Email:          util.RandomEmail(),
// 	}
// 	return
// }

// type eqCreateUserTxParamsMatcher struct {
// 	arg      db.CreateUserTxParams
// 	password string
// 	user     db.User
// }

// func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
// 	// fmt.Println(">> check param matches")
// 	actualArgs, ok := x.(db.CreateUserTxParams)
// 	if !ok {
// 		return false
// 	}

// 	// fmt.Println(">> check password", actualArgs)
// 	err := util.CheckPassword(expected.password, actualArgs.HashedPassword)
// 	if err != nil {
// 		return false
// 	}

// 	// fmt.Println(">> deep equal")
// 	expected.arg.HashedPassword = actualArgs.HashedPassword
// 	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArgs.CreateUserParams) {
// 		return false
// 	}

// 	err = actualArgs.AfterCreate(expected.user)

// 	return err == nil
// }

// // // String describes what the matcher matches.
// func (m eqCreateUserTxParamsMatcher) String() string {
// 	return fmt.Sprintf("matches arg %v and password %v", m.arg, m.password)
// }

// func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
// 	return eqCreateUserTxParamsMatcher{arg, password, user}
// }

func TestUpateeUserApi(t *testing.T) {
	user, _ := randomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()

	testCases := []struct {
		name          string
		request       pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			request: pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					WhereUsername: user.Username,
					FullName: pgtype.Text{
						String: newFullName,
						Valid:  true,
					},
					Email: pgtype.Text{
						String: newEmail,
						Valid:  true,
					},
				}

				updatedUser := db.User{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newFullName,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBeareToken(t, tokenMaker, user.Username, util.DepositorRole, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				udatedUser := res.GetUser()
				require.Equal(t, user.Username, udatedUser.Username)
				require.Equal(t, newEmail, udatedUser.Email)
				require.Equal(t, newFullName, udatedUser.FullName)
			},
		},
		{
			name: "UserNotFound",
			request: pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBeareToken(t, tokenMaker, user.Username, util.DepositorRole, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			request: pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBeareToken(t, tokenMaker, user.Username, util.DepositorRole, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			request: pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUserRequestArgument",
			request: pb.UpdateUserRequest{
				Username: "",
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBeareToken(t, tokenMaker, user.Username, util.DepositorRole, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "PermissionDenied",
			request: pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBeareToken(t, tokenMaker, "not_user", util.DepositorRole, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.PermissionDenied, st.Code())
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			// account := randomAccount()

			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)

			ctx := testCase.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUser(ctx, &testCase.request)

			testCase.checkResponse(t, res, err)

		})
	}
}

// func TestLoginUserApi(t *testing.T) {
// 	user, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		request       gin.H
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			request: gin.H{
// 				"username": user.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(user.Username)).
// 					Times(1).
// 					Return(user, nil)

// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(1)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "UserNotFound",
// 			request: gin.H{
// 				"username": "NotFound",
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.User{}, db.ErrRecordNotFound)

// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "IncorrectPassword",
// 			request: gin.H{
// 				"username": user.Username,
// 				"password": "incorrect",
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(user.Username)).
// 					Times(1).
// 					Return(user, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidUsername",
// 			request: gin.H{
// 				"username": "invalid-user#1",
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		// {
// 		// 	name: "InternalError",
// 		// 	request: gin.H{
// 		// 		"username":  user.Username,
// 		// 		"password":  password,
// 		// 		"full_name": user.FullName,
// 		// 		"email":     user.Email,
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(1).
// 		// 			Return(db.User{}, sql.ErrConnDone)
// 		// 	},
// 		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "DuplicateUsername",
// 		// 	request: gin.H{
// 		// 		"username":  user.Username,
// 		// 		"password":  password,
// 		// 		"full_name": user.FullName,
// 		// 		"email":     user.Email,
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(1).
// 		// 			Return(db.User{}, &pq.Error{Code: db.UniqueViolation})
// 		// 	},
// 		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusForbidden, recorder.Code)
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "InvalidUsername",
// 		// 	request: gin.H{
// 		// 		"username":  "user#",
// 		// 		"password":  password,
// 		// 		"full_name": user.FullName,
// 		// 		"email":     user.Email,
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(0)
// 		// 	},
// 		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
// 		// 	},
// 		// },
// 	}

// 	for i := range testCases {
// 		testCase := testCases[i]

// 		t.Run(testCase.name, func(t *testing.T) {
// 			// account := randomAccount()

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)

// 			testCase.buildStubs(store)

// 			server := newTestServer(t, store)

// 			// recorder
// 			recorder := httptest.NewRecorder()
// 			// request
// 			url := "/users/login"

// 			jsonData, err := json.Marshal(testCase.request)
// 			require.NoError(t, err)

// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)

// 			testCase.checkResponse(t, recorder)

// 		})
// 	}
// }

// func requreBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
// 	data, err := io.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotUser createUserResponse
// 	err = json.Unmarshal(data, &gotUser)
// 	require.NoError(t, err)

// 	require.Equal(t, user.Username, gotUser.Username)
// 	require.Equal(t, user.FullName, gotUser.FullName)
// 	require.Equal(t, user.Email, gotUser.Email)
// 	require.NotContains(t, "hashed_password", string(data))
// }
