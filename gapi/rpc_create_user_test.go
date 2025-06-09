package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/matodrobec/simplebank/db/mock"
	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/pb"
	"github.com/matodrobec/simplebank/util"
	"github.com/matodrobec/simplebank/worker"
	mockwk "github.com/matodrobec/simplebank/worker/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPass, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	// fmt.Println(">> check param matches")
	actualArgs, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	// fmt.Println(">> check password", actualArgs)
	err := util.CheckPassword(expected.password, actualArgs.HashedPassword)
	if err != nil {
		return false
	}

	// fmt.Println(">> deep equal")
	expected.arg.HashedPassword = actualArgs.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArgs.CreateUserParams) {
		return false
	}

	err = actualArgs.AfterCreate(expected.user)

	return err == nil
}

// // String describes what the matcher matches.
func (m eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", m.arg, m.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		request       pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			request: pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username:       user.Username,
						HashedPassword: user.HashedPassword,
						FullName:       user.FullName,
						Email:          user.Email,
					},
					AfterCreate: func(user db.User) error {
						return nil
					},
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{
						User: user,
					}, nil)

				taskArg := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().
					DistributedTaskSendEmail(
						gomock.Any(),
						taskArg,
						gomock.Any(),
					).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.Email, createdUser.Email)
				require.Equal(t, user.FullName, createdUser.FullName)
			},
		},
		{
			name: "InternalError",
			request: pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{
						User: user,
					}, sql.ErrConnDone)

				// taskArg := &worker.PayloadSendVerifyEmail{
				// 	Username: user.Username,
				// }
				taskDistributor.EXPECT().
					DistributedTaskSendEmail(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		// {
		// 	name: "InternalError",
		// 	request: gin.H{
		// 		"username":  user.Username,
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.User{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "DuplicateUsername",
		// 	request: gin.H{
		// 		"username":  user.Username,
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.User{}, &pq.Error{Code: db.UniqueViolation})
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusForbidden, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidUsername",
		// 	request: gin.H{
		// 		"username":  "user#",
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			// account := randomAccount()

			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)

			testCase.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)

			res, err := server.CreateUser(context.Background(), &testCase.request)

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
