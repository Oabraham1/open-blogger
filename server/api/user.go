package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Oabraham1/open-blogger/server/auth"
	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/Oabraham1/open-blogger/server/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserAccountRequest struct {
	Username  string `json:"username" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginUserAccountRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type GetUserAccountByUsernameRequest struct {
	Username string `uri:"username" binding:"required,alphanum,min=1"`
}

type UpdateUserInterestsRequest struct {
	Username  string   `json:"username" binding:"required"`
	Interests []string `json:"interests" binding:"required"`
}

type DeleteUserAccountRequest struct {
	Username string `uri:"username" binding:"required"`
}

type UserAccountResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	JoinedOn string `json:"joined_on"`
}

type LoginUserAccountResponse struct {
	SessionID             uuid.UUID           `json:"session_id"`
	AccessToken           string              `json:"access_token"`
	ExpiresAt             time.Time           `json:"expires_at"`
	RefreshToken          string              `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time           `json:"refresh_token_expires_at"`
	UserAccount           UserAccountResponse `json:"user_account"`
}

func GetUserAccountResponse(user db.User) UserAccountResponse {
	return UserAccountResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FirstName + " " + user.LastName,
		JoinedOn: user.CreatedAt.UTC().Format("2006/01/02 15:04:05"),
	}
}

func (server *Server) CreateUserAccount(ctx *gin.Context) {
	var req CreateUserAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateNewUserParams{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := server.DataStore.CreateNewUser(ctx, arg)
	if err != nil {
		if util.ErrorCode(err) == util.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := GetUserAccountResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.DataStore.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.VerifyPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, payload, err := server.Authenticator.CreateToken(user.Username, server.Configurations.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.Authenticator.CreateToken(user.Username, server.Configurations.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userSession, err := server.DataStore.CreateNewUserSession(ctx, db.CreateNewUserSessionParams{
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := LoginUserAccountResponse{
		SessionID:             userSession.ID,
		AccessToken:           token,
		ExpiresAt:             payload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		UserAccount:           GetUserAccountResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) GetUserByUsername(ctx *gin.Context) {
	// Modify previous implementation to check if user is logged in
	var req GetUserAccountByUsernameRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.DataStore.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	authenticationPayload := ctx.MustGet(authorizationPayloadKey).(auth.AuthPayload)
	if authenticationPayload.Username != user.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rsp := GetUserAccountResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) UpdateUserInterests(ctx *gin.Context) {
	var req UpdateUserInterestsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.DataStore.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	authenticationPayload := ctx.MustGet(authorizationPayloadKey).(auth.AuthPayload)
	if authenticationPayload.Username != user.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateUserInterestsByUsernameParams{
		Username:  req.Username,
		Interests: req.Interests,
	}

	err = server.DataStore.UpdateUserInterestsByUsername(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Interests updated successfully"})
}

func (server *Server) DeleteUserAccount(ctx *gin.Context) {
	var req DeleteUserAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println(req)

	user, err := server.DataStore.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	authenticationPayload := ctx.MustGet(authorizationPayloadKey).(auth.AuthPayload)
	if authenticationPayload.Username != user.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Get all Comments by User ID
	comments, err := server.DataStore.GetCommentsByUserName(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	for _, comment := range comments {
		err = server.DataStore.DeleteCommentByID(ctx, comment.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// Get all Posts by User ID
	posts, err := server.DataStore.GetPostsByUserName(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	for _, post := range posts {
		err = server.DataStore.DeletePostByID(ctx, post.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	err = server.DataStore.DeleteUserAccount(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
