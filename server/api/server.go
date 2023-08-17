package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Oabraham1/open-blogger/server/auth"
	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/Oabraham1/open-blogger/server/util"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) UnauthorizedError(ctx *gin.Context) {
	err := errors.New("Unauthorized")
	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
}

func (server *Server) InternalServerError(ctx *gin.Context) {
	err := errors.New("internal server error")
	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
}

func (server *Server) BadRequestError(ctx *gin.Context) {
	err := errors.New("bad request")
	ctx.JSON(http.StatusBadRequest, errorResponse(err))
}

func (server *Server) NotFoundError(ctx *gin.Context) {
	err := errors.New("not found")
	ctx.JSON(http.StatusNotFound, errorResponse(err))
}

func (server *Server) ForbiddenError(ctx *gin.Context) {
	err := errors.New("forbidden")
	ctx.JSON(http.StatusForbidden, errorResponse(err))
}

func (server *Server) ConflictError(ctx *gin.Context) {
	err := errors.New("Conflict")
	ctx.JSON(http.StatusConflict, errorResponse(err))
}

func (server *Server) ReturnOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

/* Server serves the HTTP Requests */
type Server struct {
	Router         *gin.Engine
	DataStore      db.Store
	Configurations util.Config
	Authenticator  auth.Authenticator
}

/* NewServer creates a new server */
func NewServer(store db.Store, config util.Config) (*Server, error) {
	authenticator, err := auth.NewPasetoAuthenticator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create authenticator: %w", err)
	}
	server := &Server{
		DataStore:      store,
		Configurations: config,
		Authenticator:  authenticator,
	}
	server.setupRouter()
	return server, nil
}

/* SetupRouter sets up the router */
func (server *Server) setupRouter() {
	router := gin.Default()

	authenticatedRoutes := router.Group("/").Use(AuthenticationMiddleware(server.Authenticator))

	router.POST("/api/user/create", server.CreateUserAccount)
	router.POST("/api/user/login", server.LoginUser)
	authenticatedRoutes.GET("/api/user/getByUsername/:username", server.GetUserByUsername)
	authenticatedRoutes.PUT("/api/user/updateInterests", server.UpdateUserInterests)
	authenticatedRoutes.DELETE("/api/user/delete/:username", server.DeleteUserAccount)

	authenticatedRoutes.POST("/api/post/create", server.CreateNewPost)
	router.GET("/api/post/getByID/:id", server.GetPostById)
	router.GET("/api/post/getByCategory/:category", server.GetPostsByCategory)
	router.GET("/api/post/getByUsername/:username", server.GetPostsByUsername)
	authenticatedRoutes.PUT("/api/post/updateBody", server.UpdatePostBody)
	authenticatedRoutes.PUT("/api/post/publish", server.UpdatePostStatus)
	authenticatedRoutes.DELETE("/api/post/delete/:id", server.DeletePost)

	authenticatedRoutes.POST("/api/comment/create", server.CreateNewComment)
	router.GET("/api/comment/getByPostID/:id", server.GetCommentsByPostID)
	authenticatedRoutes.DELETE("/api/comment/delete/:id", server.DeleteComment)

	server.Router = router
}

/* StartServer starts the server */
func (server *Server) StartServer(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
