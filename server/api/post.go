package api

import (
	"net/http"
	"time"

	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Body     string `json:"body" binding:"required"`
	UserName string `json:"username" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Category string `json:"category" binding:"required"`
}

type CreateNewCommentRequest struct {
	PostID   string `json:"post_id" binding:"required"`
	Body     string `json:"body" binding:"required"`
	UserName string `json:"username" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
}

type GetPostByIDRequest struct {
	ID string `json:"id" binding:"required"`
}

type GetPostsByCategoryRequest struct {
	Category string `uri:"category" binding:"required"`
}

type GetPostsByUsernameRequest struct {
	Username string `uri:"username" binding:"required"`
}

type UpdatePostRequest struct {
	ID     string `json:"id" binding:"required"`
	Body   string `json:"body" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

type GetCommentsByPostIDRequest struct {
	PostID string `json:"post_id" binding:"required"`
}

type DeletePostRequest struct {
	ID string `json:"id" binding:"required"`
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
		LastModified: post.LastModified.Local().String(),
		PublishedAt:  post.PublishedAt.UTC().String(),
	}
}

func GetCommentResponse(comment db.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID.String(),
		PostID:    comment.PostID.String(),
		Body:      comment.Body,
		UserName:  comment.Username,
		CreatedAt: comment.CreatedAt.Local().String(),
	}
}

func (server *Server) CreateNewPost(ctx *gin.Context) {
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert userId string to uuid
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	status := req.Status
	if status != string(db.StatusDraft) && status != string(db.StatusPublished) {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dbStatus := db.Status(status)

	arg := db.CreateNewPostParams{
		ID:           uuid.New(),
		Title:        req.Title,
		Body:         req.Body,
		Username:     req.UserName,
		UserID:       userId,
		Status:       dbStatus,
		Category:     req.Category,
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}

	if dbStatus == db.StatusPublished {
		arg.PublishedAt = time.Now()
	}

	post, err := server.DataStore.CreateNewPost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetPostResponse(post))
}

func (server *Server) CreateNewComment(ctx *gin.Context) {
	var req CreateNewCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert userId string to uuid
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert postId string int32
	postId, err := uuid.Parse(req.PostID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateNewCommentParams{
		ID:       uuid.New(),
		PostID:   postId,
		Body:     req.Body,
		Username: req.UserName,
		UserID:   userId,
	}

	comment, err := server.DataStore.CreateNewComment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetCommentResponse(comment))
}

func (server *Server) GetPostsByCategory(ctx *gin.Context) {
	var req GetPostsByCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	posts, err := server.DataStore.GetPostsByCategory(ctx, req.Category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []PostResponse
	for _, post := range posts {
		rsp = append(rsp, GetPostResponse(post))
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) GetPostById(ctx *gin.Context) {
	var req GetPostByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	post, err := server.DataStore.GetPostById(ctx, postId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetPostResponse(post))
}

func (server *Server) GetPostsByUsername(ctx *gin.Context) {
	var req GetPostsByUsernameRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// find user by username
	_, err := server.DataStore.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	posts, err := server.DataStore.GetPostsByUserName(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []PostResponse
	for _, post := range posts {
		rsp = append(rsp, GetPostResponse(post))
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) GetPostsByUserID(ctx *gin.Context) {
	var req GetUserByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert userId string to uuid
	userId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// find user by username
	_, err = server.DataStore.GetUserByID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	posts, err := server.DataStore.GetPostsByUserID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []PostResponse
	for _, post := range posts {
		rsp = append(rsp, GetPostResponse(post))
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) UpdatePostBody(ctx *gin.Context) {
	var req UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert userId string to uuid
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdatePostBodyByPostIDAndUserIDParams{
		ID:     postId,
		Body:   req.Body,
		UserID: userId,
	}

	post, err := server.DataStore.UpdatePostBodyByPostIDAndUserID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetPostResponse(post))
}

func (server *Server) GetCommentsByPostID(ctx *gin.Context) {
	var req GetCommentsByPostIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.PostID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	comments, err := server.DataStore.GetCommentsByPostID(ctx, postId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []CommentResponse
	for _, comment := range comments {
		rsp = append(rsp, GetCommentResponse(comment))
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) DeletePost(ctx *gin.Context) {
	var req DeletePostRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert postId string to uuid
	postId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.DataStore.DeletePostByID(ctx, postId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
