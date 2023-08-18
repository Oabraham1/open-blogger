package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		{
			name: "User Not Found",
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
					Return(db.User{}, util.ErrRecordNotFound)
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

func TestGetPostByCategory(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	post.Status = "published"
	post.LastModified = time.Now().Format("2006-01-02 15:04:05")

	testCases := []struct {
		name          string
		category      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			category: post.Category,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostsByCategory(gomock.Any(), gomock.Eq(post.Category)).
					Times(1).
					Return([]db.Post{post}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var posts []PostResponse

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &posts)
				require.NoError(t, err)
				require.Equal(t, post.Body, posts[0].Body)
			},
		},
		{
			name:     "No Posts Found",
			category: post.Category,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostsByCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Post{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var posts []db.Post

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &posts)
				require.NoError(t, err)
				require.Equal(t, len([]db.Post{}), len(posts))
			},
		},
		{
			name:     "Invalid Category",
			category: "$",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostsByCategory(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/post/getByCategory/%s", tc.category)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetPostById(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	post.Status = "published"
	post.LastModified = time.Now().Format("2006-01-02 15:04:05")

	testCases := []struct {
		name          string
		id            string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   post.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var postResponse PostResponse

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &postResponse)
				require.NoError(t, err)
				require.Equal(t, post.Body, postResponse.Body)
			},
		},
		{
			name: "No Posts Found",
			id:   post.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Post{}, util.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Invalid Post ID",
			id:   "  ",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/post/getByID/%s", tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetPostByUserName(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	post.Status = "published"
	post.LastModified = time.Now().Format("2006-01-02 15:04:05")

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: post.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetPostsByUserName(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return([]db.Post{post}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var posts []PostResponse

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &posts)
				require.NoError(t, err)
				require.Equal(t, post.Body, posts[0].Body)
			},
		},
		{
			name:     "No Posts Found",
			username: post.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetPostsByUserName(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Post{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var posts []db.Post

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &posts)
				require.NoError(t, err)
				require.Equal(t, len([]db.Post{}), len(posts))
			},
		},
		{
			name:     "Invalid Username",
			username: "$",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostsByUserName(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "User Not Found",
			username: "invalidUsername",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq("invalidUsername")).
					Times(1).
					Return(db.User{}, util.ErrRecordNotFound)
				store.EXPECT().
					GetPostsByUserName(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/post/getByUsername/%s", tc.username)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdatePostBody(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	post.Status = "published"
	post.LastModified = time.Now().Format("2006-01-02 15:04:05")

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
				"body":     post.Body,
				"username": post.Username,
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostBodyParams{
					Body:         post.Body,
					Username:     post.Username,
					ID:           post.ID,
					LastModified: post.LastModified,
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
					UpdatePostBody(gomock.Any(), arg).
					Times(1).
					Return(post, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var postResponse PostResponse

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &postResponse)
				require.NoError(t, err)
				require.Equal(t, post.Body, postResponse.Body)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"body":     post.Body,
				"username": post.Username,
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostBodyParams{
					Body:     post.Body,
					Username: post.Username,
					ID:       post.ID,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(0)
				store.EXPECT().
					UpdatePostBody(gomock.Any(), arg).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Updating Another User's Post",
			body: gin.H{
				"body":     post.Body,
				"username": post.Username,
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, "anotherUserName", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostBodyParams{
					Body:     post.Body,
					Username: post.Username,
					ID:       post.ID,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					UpdatePostBody(gomock.Any(), arg).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"body":     post.Body,
				"username": "_FJ",
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, post.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostBodyParams{
					Body:     post.Body,
					Username: post.Username,
					ID:       post.ID,
				}
				store.EXPECT().
					UpdatePostBody(gomock.Any(), arg).
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

			url := "/api/post/updateBody"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdatePostStatus(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	post.Status = "draft"
	post.PublishedAt = time.Now().Format("2006-01-02 15:04:05")
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
				"username": post.Username,
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, post.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostStatusParams{
					Username:    post.Username,
					ID:          post.ID,
					PublishedAt: post.PublishedAt,
					Status:      db.StatusPublished,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					UpdatePostStatus(gomock.Any(), arg).
					Times(1).
					Return(db.Post{
						ID:          post.ID,
						Username:    post.Username,
						Body:        post.Body,
						Status:      "published",
						PublishedAt: post.PublishedAt,
						Category:    post.Category,
						CreatedAt:   post.CreatedAt,
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var postResponse PostResponse

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &postResponse)
				require.NoError(t, err)
				require.NotEqual(t, string(post.Status), postResponse.Status)
				require.Equal(t, string(db.StatusPublished), postResponse.Status)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"username": post.Username,
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostStatusParams{
					Username:    post.Username,
					ID:          post.ID,
					PublishedAt: post.PublishedAt,
					Status:      db.StatusPublished,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(0)
				store.EXPECT().
					UpdatePostStatus(gomock.Any(), arg).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Publishing Another User's Post",
			body: gin.H{
				"username": post.Username,
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, "anotherUser", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostStatusParams{
					Username:    post.Username,
					ID:          post.ID,
					PublishedAt: post.PublishedAt,
					Status:      db.StatusPublished,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					UpdatePostStatus(gomock.Any(), arg).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"username": "_FJ",
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, post.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostStatusParams{
					Username:    post.Username,
					ID:          post.ID,
					PublishedAt: post.PublishedAt,
					Status:      db.StatusPublished,
				}
				store.EXPECT().
					UpdatePostStatus(gomock.Any(), arg).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Post Not Found",
			body: gin.H{
				"username": post.Username,
				"post_id":  invalidPostId,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, post.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostStatusParams{
					Username:    post.Username,
					ID:          invalidPostId,
					PublishedAt: post.PublishedAt,
					Status:      db.StatusPublished,
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(invalidPostId)).
					Times(1).
					Return(db.Post{}, util.ErrRecordNotFound)
				store.EXPECT().
					UpdatePostStatus(gomock.Any(), arg).
					Times(0).
					Return(db.Post{}, util.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "User Not Found",
			body: gin.H{
				"username": "invalidUsername",
				"post_id":  post.ID,
			},
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, post.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq("invalidUsername")).
					Times(1).
					Return(db.User{}, util.ErrRecordNotFound)
				store.EXPECT().
					UpdatePostStatus(gomock.Any(), gomock.Any()).
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

			url := "/api/post/publish"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	comment := generateDummyComment(t, user, post)

	randomPostId := uuid.New()

	testCases := []struct {
		name          string
		postID        string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			postID: post.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					GetCommentsByPostID(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return([]db.Comment{comment}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var comments []CommentResponse

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &comments)
				require.NoError(t, err)
				require.Equal(t, comment.Body, comments[0].Body)
			},
		},
		{
			name:   "No Comments Found",
			postID: post.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					GetCommentsByPostID(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Comment{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var comments []db.Comment

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &comments)
				require.NoError(t, err)
				require.Equal(t, len([]db.Comment{}), len(comments))
			},
		},
		{
			name:   "Invalid Post ID",
			postID: "  ",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCommentsByPostID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Post Not Found",
			postID: randomPostId.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(randomPostId)).
					Times(1).
					Return(db.Post{}, util.ErrRecordNotFound)
				store.EXPECT().
					GetCommentsByPostID(gomock.Any(), randomPostId).
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

			url := fmt.Sprintf("/api/comment/getByPostID/%s", tc.postID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)
	comment := generateDummyComment(t, user, post)

	randomCommentId := uuid.New()

	testCases := []struct {
		name               string
		commentID          string
		setUpAuthenticator func(t *testing.T, request *http.Request, authenticator auth.Authenticator)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			commentID: comment.ID.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCommentByID(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(comment, nil)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					DeleteCommentByID(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:               "Unauthorized",
			commentID:          comment.ID.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCommentByID(gomock.Any(), gomock.Eq(comment.ID)).
					Times(0)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(0)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(0)
				store.EXPECT().
					DeleteCommentByID(gomock.Any(), gomock.Eq(comment.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Deleting Another User's Comment",
			commentID: comment.ID.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, "anotherUser", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCommentByID(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(comment, nil)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(comment.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					DeleteCommentByID(gomock.Any(), gomock.Eq(comment.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Invalid Comment ID",
			commentID: "  ",
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCommentByID(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteCommentByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Comment Not Found",
			commentID: randomCommentId.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCommentByID(gomock.Any(), gomock.Eq(randomCommentId)).
					Times(1).
					Return(db.Comment{}, util.ErrRecordNotFound)
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteCommentByID(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/comment/delete/%s", tc.commentID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeletePost(t *testing.T) {
	user, _ := generateDummyUser(t)
	post := generateDummyPost(t, user)

	randomPostId := uuid.New()

	testCases := []struct {
		name               string
		postID             string
		setUpAuthenticator func(t *testing.T, request *http.Request, authenticator auth.Authenticator)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			postID: post.ID.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetCommentsByPostID(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return([]db.Comment{}, nil)
				store.EXPECT().
					DeletePostByID(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:               "Unauthorized",
			postID:             post.ID.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(0)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(0)
				store.EXPECT().
					DeletePostByID(gomock.Any(), gomock.Eq(post.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "Deleting Another User's Post",
			postID: post.ID.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, "anotherUser", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(post.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					DeletePostByID(gomock.Any(), gomock.Eq(post.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "Invalid Post ID",
			postID: "  ",
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeletePostByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Post Not Found",
			postID: randomPostId.String(),
			setUpAuthenticator: func(t *testing.T, request *http.Request, authenticator auth.Authenticator) {
				addAuth(t, request, authenticator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(randomPostId)).
					Times(1).
					Return(db.Post{}, util.ErrRecordNotFound)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeletePostByID(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/post/delete/%s", tc.postID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setUpAuthenticator(t, request, server.Authenticator)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
