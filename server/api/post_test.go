package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Oabraham1/open-blogger/server/auth"
	mockdb "github.com/Oabraham1/open-blogger/server/db/mock"
	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/Oabraham1/open-blogger/server/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateDummyPost(t *testing.T, user db.User) db.Post {
	return db.Post{
		Title:    "Test Post",
		Body:     "This is a test post",
		Username: user.Username,
		Category: "testCategory",
	}
}

func generateDummyComment(t *testing.T, user db.User, post db.Post) db.Comment {
	return db.Comment{
		Body:     "This is a test comment",
		Username: user.Username,
		PostID:   post.ID,
	}
}

func TestCreatePost(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	post.Status = "draft"

	testCases := []struct {
		name               string
		body               gin.H
		setUpAuthenticator func(t *testing.T, request *http.Request, authenticator auth.Authenticator)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":    post.Title,
				"body":     post.Body,
				"category": post.Category,
				"username": post.Username,
				"status":   post.Status,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNewPostParams{
					Title:    post.Title,
					Body:     post.Body,
					Username: post.Username,
					Category: post.Category,
					Status:   post.Status,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateNewPost(gomock.Any(), arg).
					Times(1).
					Return(post, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"title":    post.Title,
				"body":     post.Body,
				"category": post.Category,
				"username": post.Username,
				"status":   post.Status,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNewPostParams{
					Title:    post.Title,
					Body:     post.Body,
					Username: post.Username,
					Category: post.Category,
					Status:   post.Status,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(0)
				store.EXPECT().
					CreateNewPost(gomock.Any(), arg).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Incorrect Username",
			body: gin.H{
				"title":    post.Title,
				"body":     post.Body,
				"category": post.Category,
				"username": user.Username,
				"status":   post.Status,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, "incorrectUserName", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateNewPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Empty Title",
			body: gin.H{
				"title":    "",
				"body":     post.Body,
				"category": post.Category,
				"username": post.Username,
				"status":   post.Status,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(0)
				store.EXPECT().
					CreateNewPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Status",
			body: gin.H{
				"title":    post.Title,
				"body":     post.Body,
				"category": post.Category,
				"username": post.Username,
				"status":   "invalidStatus",
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1)
				store.EXPECT().
					CreateNewPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "User Not Found",
			body: gin.H{
				"title":    post.Title,
				"body":     post.Body,
				"category": post.Category,
				"username": post.Username,
				"status":   post.Status,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, util.ErrRecordNotFound)
				store.EXPECT().
					CreateNewPost(gomock.Any(), gomock.Any()).
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

			url := "/api/post/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateNewComment(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	comment := generateDummyComment(t, user, post)

	invalidPostId := uuid.New()

	testCases := []struct {
		name               string
		body               gin.H
		setUpAuthenticator func(t *testing.T, request *http.Request, authenticator auth.Authenticator)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"body":     comment.Body,
				"username": comment.Username,
				"post_id":  comment.PostID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNewCommentParams{
					Body:     comment.Body,
					Username: comment.Username,
					PostID:   comment.PostID,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					CreateNewComment(gomock.Any(), arg).
					Times(1).
					Return(comment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"body":     comment.Body,
				"username": comment.Username,
				"post_id":  comment.PostID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNewCommentParams{
					Body:     comment.Body,
					Username: comment.Username,
					PostID:   comment.PostID,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(0)
				store.EXPECT().
					CreateNewComment(gomock.Any(), arg).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Incorrect Username",
			body: gin.H{
				"body":     comment.Body,
				"username": user.Username,
				"post_id":  comment.PostID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, "incorrectUserName", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateNewComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Nonexistent Post",
			body: gin.H{
				"body":     comment.Body,
				"username": comment.Username,
				"post_id":  invalidPostId,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNewCommentParams{
					Body:     comment.Body,
					Username: comment.Username,
					PostID:   comment.PostID,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetPostById(gomock.Any(), invalidPostId).
					Times(1).
					Return(db.Post{}, util.ErrRecordNotFound)
				store.EXPECT().
					CreateNewComment(gomock.Any(), arg).
					Times(0).
					Return(db.Comment{}, util.ErrRecordNotFound)
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

			url := "/api/comment/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
