package api

import (
	"fmt"
	"net/http"
	"time"

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

type GetUserByIDRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type UpdateUserInterestsRequest struct {
	ID        string   `json:"id" binding:"required"`
	Interests []string `json:"interests" binding:"required"`
}

type DeleteUserAccountRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UserAccountResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	JoinedOn string `json:"joined_on"`
}

func GetUserAccountResponse(user db.User) UserAccountResponse {
	return UserAccountResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FirstName + " " + user.LastName,
		JoinedOn: user.CreatedAt.UTC().String(),
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
		ID:        uuid.New(),
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Interests: []string{},
		CreatedAt: time.Now(),
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = util.VerifyPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rsp := GetUserAccountResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) GetUserByUsername(ctx *gin.Context) {
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

	rsp := GetUserAccountResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) GetUserByID(ctx *gin.Context) {
	var req GetUserByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Convert string to uuid
	userId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.DataStore.GetUserByID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
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

	// Convert string to uuid
	userId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.DataStore.GetUserByID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.UpdateUserInterestsByIDParams{
		ID:        userId,
		Interests: req.Interests,
	}

	err = server.DataStore.UpdateUserInterestsByID(ctx, arg)
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

	// Convert string to uuid
	userId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.DataStore.GetUserByID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = server.DataStore.DeleteUserByID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
