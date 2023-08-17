package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/Oabraham1/open-blogger/server/auth"
	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/Oabraham1/open-blogger/server/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Body     string `json:"body" binding:"required"`
	UserName string `json:"username" binding:"required,alphanum"`
	Status   string `json:"status" binding:"required,oneof=draft published"`
	Category string `json:"category" binding:"required"`
}

type CreateNewCommentRequest struct {
	PostID   string `json:"post_id" binding:"required"`
	Body     string `json:"body" binding:"required"`
	UserName string `json:"username" binding:"required,alphanum"`
}

type GetPostByIDRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type GetPostsByCategoryRequest struct {
	Category string `uri:"category" binding:"required,alphanum,min=1"`
}

type GetPostsByUsernameRequest struct {
	Username string `uri:"username" binding:"required,alphanum,min=1"`
}

type UpdatePostBodyRequest struct {
	ID       string `json:"post_id" binding:"required"`
	Body     string `json:"body" binding:"required"`
	UserName string `json:"username" binding:"required,alphanum"`
}

type UpdatePostStatusRequest struct {
	ID       string `json:"post_id" binding:"required"`
	UserName string `json:"username" binding:"required,alphanum"`
}

type GetCommentsByPostIDRequest struct {
	PostID string `uri:"id" binding:"required,min=1"`
}

type DeleteCommentByIDRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type DeletePostRequest struct {
	ID string `json:"id" binding:"required,min=1"`
}

type PostResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	UserName     string `json:"username"`
	Status       string `json:"status"`
	Category     string `json:"category"`
	LastModified string `json:"last_modified"`
	PublishedAt  string `json:"published_at"`
}

type CommentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	Body      string `json:"body"`
	UserName  string `json:"username"`
	CreatedAt string `json:"commented_at"`
}

func GetPostResponse(post db.Post) PostResponse {
	return PostResponse{
		ID:           post.ID.String(),
		Title:        post.Title,
		Body:         post.Body,
		UserName:     post.Username,
		Status:       string(post.Status),
		Category:     post.Category,
		LastModified: post.LastModified,
		PublishedAt:  post.PublishedAt,
	}
}

func GetCommentResponse(comment db.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID.String(),
		PostID:    comment.PostID.String(),
		Body:      comment.Body,
		UserName:  comment.Username,
		CreatedAt: comment.CreatedAt,
	}
}

func (server *Server) CreateNewPost(ctx *gin.Context) {
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	user, err := server.DataStore.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	if user.Username != req.UserName {
		server.BadRequestError(ctx)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			server.UnauthorizedError(ctx)
		}
	}()

	authenticationPayload, ok := ctx.MustGet(authorizationPayloadKey).(*auth.AuthPayload)
	if !ok {
		server.UnauthorizedError(ctx)
		return
	}

	if authenticationPayload.Username != user.Username || authenticationPayload.Username != req.UserName {
		server.UnauthorizedError(ctx)
		return
	}

	status := req.Status
	if status != string(db.StatusDraft) && status != string(db.StatusPublished) {
		server.BadRequestError(ctx)
		return
	}
	dbStatus := db.Status(status)

	arg := db.CreateNewPostParams{
		Title:    req.Title,
		Body:     req.Body,
		Username: req.UserName,
		Status:   dbStatus,
		Category: req.Category,
	}

	if dbStatus == db.StatusPublished {
		arg.PublishedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	post, err := server.DataStore.CreateNewPost(ctx, arg)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	server.ReturnOK(ctx, GetPostResponse(post))
}

func (server *Server) CreateNewComment(ctx *gin.Context) {
	var req CreateNewCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	user, err := server.DataStore.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	if user.Username != req.UserName {
		server.BadRequestError(ctx)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			server.UnauthorizedError(ctx)
		}
	}()

	authenticationPayload, ok := ctx.MustGet(authorizationPayloadKey).(*auth.AuthPayload)
	if !ok {
		server.UnauthorizedError(ctx)
		return
	}

	if authenticationPayload.Username != user.Username || authenticationPayload.Username != req.UserName {
		server.UnauthorizedError(ctx)
		return
	}

	// convert postId string int32
	postId, err := uuid.Parse(req.PostID)
	if err != nil {
		server.BadRequestError(ctx)
		return
	}

	// find post
	_, err = server.DataStore.GetPostById(ctx, postId)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	arg := db.CreateNewCommentParams{
		PostID:   postId,
		Body:     req.Body,
		Username: req.UserName,
	}

	comment, err := server.DataStore.CreateNewComment(ctx, arg)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	server.ReturnOK(ctx, GetCommentResponse(comment))
}

func (server *Server) GetPostsByCategory(ctx *gin.Context) {
	var req GetPostsByCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	posts, err := server.DataStore.GetPostsByCategory(ctx, req.Category)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	var rsp []PostResponse
	for _, post := range posts {
		if post.Status == db.StatusPublished {
			rsp = append(rsp, GetPostResponse(post))
		}
	}

	server.ReturnOK(ctx, rsp)
}

func (server *Server) GetPostById(ctx *gin.Context) {
	var req GetPostByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.ID)
	if err != nil {
		server.BadRequestError(ctx)
		return
	}
	post, err := server.DataStore.GetPostById(ctx, postId)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	server.ReturnOK(ctx, GetPostResponse(post))
}

func (server *Server) GetPostsByUsername(ctx *gin.Context) {
	var req GetPostsByUsernameRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	// find user by username
	_, err := server.DataStore.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	posts, err := server.DataStore.GetPostsByUserName(ctx, req.Username)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	var rsp []PostResponse
	for _, post := range posts {
		rsp = append(rsp, GetPostResponse(post))
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) UpdatePostBody(ctx *gin.Context) {
	var req UpdatePostBodyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	user, err := server.DataStore.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			server.UnauthorizedError(ctx)
		}
	}()

	authenticationPayload, ok := ctx.MustGet(authorizationPayloadKey).(*auth.AuthPayload)
	if !ok {
		server.UnauthorizedError(ctx)
		return
	}

	if authenticationPayload.Username != user.Username || authenticationPayload.Username != req.UserName {
		server.UnauthorizedError(ctx)
		return
	}
	// convert postId string to uuid
	postId, err := uuid.Parse(req.ID)
	if err != nil {
		server.BadRequestError(ctx)
		return
	}

	post, err := server.DataStore.GetPostById(ctx, postId)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	if post.Username != req.UserName || authenticationPayload.Username != post.Username {
		server.UnauthorizedError(ctx)
		return
	}

	arg := db.UpdatePostBodyParams{
		ID:           postId,
		Body:         req.Body,
		Username:     req.UserName,
		LastModified: time.Now().Format("2006-01-02 15:04:05"),
	}

	post, err = server.DataStore.UpdatePostBody(ctx, arg)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	server.ReturnOK(ctx, GetPostResponse(post))
}

func (server *Server) UpdatePostStatus(ctx *gin.Context) {
	var req UpdatePostStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	user, err := server.DataStore.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			server.UnauthorizedError(ctx)
		}
	}()

	authenticationPayload, ok := ctx.MustGet(authorizationPayloadKey).(*auth.AuthPayload)
	if !ok {
		server.UnauthorizedError(ctx)
		return
	}

	if authenticationPayload.Username != user.Username || authenticationPayload.Username != req.UserName {
		server.UnauthorizedError(ctx)
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.ID)
	if err != nil {
		server.BadRequestError(ctx)
		return
	}

	post, err := server.DataStore.GetPostById(ctx, postId)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	if post.Username != req.UserName || authenticationPayload.Username != post.Username {
		server.UnauthorizedError(ctx)
		return
	}

	arg := db.UpdatePostStatusParams{
		ID:          postId,
		Status:      db.StatusPublished,
		Username:    req.UserName,
		PublishedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	post, err = server.DataStore.UpdatePostStatus(ctx, arg)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	server.ReturnOK(ctx, GetPostResponse(post))
}

func (server *Server) GetCommentsByPostID(ctx *gin.Context) {
	var req GetCommentsByPostIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.PostID)
	if err != nil {
		server.BadRequestError(ctx)
		return
	}

	comments, err := server.DataStore.GetCommentsByPostID(ctx, postId)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	var rsp []CommentResponse
	for _, comment := range comments {
		rsp = append(rsp, GetCommentResponse(comment))
	}

	server.ReturnOK(ctx, rsp)
}

func (server *Server) DeleteComment(ctx *gin.Context) {
	var req DeleteCommentByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	// convert commentId string to uuid
	commentId, err := uuid.Parse(req.ID)
	if err != nil {
		server.BadRequestError(ctx)
		return
	}

	// find the comment
	comment, err := server.DataStore.GetCommentByID(ctx, commentId)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	// find the user
	user, err := server.DataStore.GetUserByUsername(ctx, comment.Username)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			server.UnauthorizedError(ctx)
		}
	}()

	authenticationPayload, ok := ctx.MustGet(authorizationPayloadKey).(*auth.AuthPayload)
	if !ok {
		server.UnauthorizedError(ctx)
		return
	}

	if authenticationPayload.Username != user.Username || authenticationPayload.Username != comment.Username {
		server.UnauthorizedError(ctx)
		return
	}

	err = server.DataStore.DeleteCommentByID(ctx, commentId)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	server.ReturnOK(ctx, gin.H{"message": "Comment deleted successfully"})
}

func (server *Server) DeletePost(ctx *gin.Context) {
	var req DeletePostRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.ID)
	if err != nil {
		server.BadRequestError(ctx)
		return
	}

	// find the post
	post, err := server.DataStore.GetPostById(ctx, postId)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	// find the user
	user, err := server.DataStore.GetUserByUsername(ctx, post.Username)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			server.UnauthorizedError(ctx)
		}
	}()

	authenticationPayload, ok := ctx.MustGet(authorizationPayloadKey).(*auth.AuthPayload)
	if !ok {
		server.UnauthorizedError(ctx)
		return
	}

	if authenticationPayload.Username != user.Username || authenticationPayload.Username != post.Username {
		server.UnauthorizedError(ctx)
		return
	}

	// Delete all comments by postID
	comments, err := server.DataStore.GetCommentsByPostID(ctx, postId)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}
	for _, comment := range comments {
		err = server.DataStore.DeleteCommentByID(ctx, comment.ID)
		if err != nil {
			server.InternalServerError(ctx)
			return
		}
	}

	err = server.DataStore.DeletePostByID(ctx, postId)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	server.ReturnOK(ctx, gin.H{"message": "Post deleted successfully"})
}
