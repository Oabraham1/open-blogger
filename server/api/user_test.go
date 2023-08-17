package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Oabraham1/open-blogger/server/auth"
	mockdb "github.com/Oabraham1/open-blogger/server/db/mock"
	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/Oabraham1/open-blogger/server/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func addAuth(t *testing.T, request *http.Request, authenticator auth.Authenticator, authType string, username string, duration time.Duration) {
	token, payload, err := authenticator.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	authorizationHeader := fmt.Sprintf("%s %s", authType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateNewUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateNewUserParams)
	if !ok {
		return false
	}

	err := util.VerifyPassword(arg.Password, e.password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateNewUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func generateDummyUser(t *testing.T) (user db.User, password string) {
	password = "password123"
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username:  "testuser",
		Password:  hashedPassword,
		Email:     "test@email.com",
		FirstName: "test",
		LastName:  "user",
	}
	return
}

func TestCreateNewUserAccount(t *testing.T) {
	user, password := generateDummyUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNewUserParams{
					Username:  user.Username,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
					Password:  user.Password,
				}
				store.EXPECT().
					CreateNewUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotUser db.User
				err = json.Unmarshal(data, &gotUser)

				require.NoError(t, err)
				require.Equal(t, user.Username, gotUser.Username)
				require.Equal(t, user.Email, gotUser.Email)
				require.Empty(t, gotUser.Password)
			},
		},

		{
			name: "InternalError",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, util.ErrUniqueViolation)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":   "invalid-user#1",
				"password":   password,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      "invalid-email",
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":   user.Username,
				"password":   "123",
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/user/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUser(t *testing.T) {
	user, password := generateDummyUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateNewUserSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, util.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Username,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-user#1",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/user/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	user, _ := generateDummyUser(t)
	testCases := []struct {
		name               string
		username           string
		setUpAuthenticator func(t *testing.T, request *http.Request, authenticator auth.Authenticator)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotUser db.User
				err = json.Unmarshal(data, &gotUser)
				require.NoError(t, err)
				require.Equal(t, user.Username, gotUser.Username)
			},
		},
		{
			name:               "Unauthorized",
			username:           user.Username,
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/user/getByUsername/%s", tc.username)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateUserInterests(t *testing.T) {
	user, _ := generateDummyUser(t)
	testCases := []struct {
		name               string
		username           string
		body               gin.H
		setUpAuthenticator func(t *testing.T, request *http.Request, authenticator auth.Authenticator)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			body: gin.H{
				"username":  user.Username,
				"interests": []string{"test", "interests"},
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserInterestsByUsernameParams{
					Username:  user.Username,
					Interests: []string{"test", "interests"},
				}
				store.EXPECT().
					UpdateUserInterestsByUsername(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response gin.H
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)
				require.Equal(t, "Interests updated successfully", response["message"])
			},
		},
		{
			name:     "Unauthorized",
			username: user.Username,
			body: gin.H{
				"username":  user.Username,
				"interests": []string{"test", "interests"},
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserInterestsByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "InvalidUsername",
			username: user.Username,
			body: gin.H{
				"username":  "invalid-user#1",
				"interests": []string{"test", "interests"},
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserInterestsByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "UserNotFound",
			username: user.Username,
			body: gin.H{
				"username":  user.Username,
				"interests": []string{"test", "interests"},
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, util.ErrRecordNotFound)
				store.EXPECT().
					UpdateUserInterestsByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/user/updateInterests"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteUserAccount(t *testing.T) {
	user, _ := generateDummyUser(t)
	testCases := []struct {
		name               string
		username           string
		setUpAuthenticator func(t *testing.T, request *http.Request, authenticator auth.Authenticator)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetCommentsByUserName(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return([]db.Comment{}, nil)
				store.EXPECT().
					GetPostsByUserName(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return([]db.Post{}, nil)
				store.EXPECT().
					GetUserSessionsByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return([]db.Session{}, nil)
				store.EXPECT().
					DeleteUserAccount(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response gin.H
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)
				require.Equal(t, "User account deleted successfully", response["message"])
			},
		},
		{
			name:               "Unauthorized",
			username:           user.Username,
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUserAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "InvalidUsername",
			username: "invalid-user#1",
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUserAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "UserNotFound",
			username: user.Username,
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, util.ErrRecordNotFound)
				store.EXPECT().
					DeleteUserAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/user/delete/%s", tc.username)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
