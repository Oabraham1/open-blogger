package api

import (
	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/Oabraham1/open-blogger/server/util"
	"github.com/gin-gonic/gin"
)

/* Server serves the HTTP Requests */
type Server struct {
	Router         *gin.Engine
	DataStore      db.Store
	Configurations util.Config
}

/* NewServer creates a new server */
func NewServer(store db.Store, config util.Config) *Server {
	server := &Server{
		DataStore:      store,
		Configurations: config,
	}
	server.setupRouter()
	return server
}

/* SetupRouter sets up the router */
func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/api/user/create", server.CreateUserAccount)
	router.POST("/api/user/login", server.LoginUser)
	router.GET("/api/user/getByUsername/:userName", server.GetUserByUsername)
	router.GET("/api/user/getById/:id", server.GetUserByID)
	router.POST("/api/user/updateInterests", server.UpdateUserInterests)
	router.DELETE("/api/user/delete/:id", server.DeleteUserAccount)

	router.POST("/api/post/create", server.CreateNewPost)
	router.GET("/api/post/getByID/:id", server.GetPostById)
	router.GET("/api/post/getByCategory", server.GetPostsByCategory)
	router.GET("/api/post/getByUsername/:username", server.GetPostsByUsername)
	router.POST("/api/post/updateBody", server.UpdatePostBody)
	router.POST("/api/post/publish", server.UpdatePostStatus)
	router.DELETE("/api/post/delete/:id", server.DeletePost)

	router.POST("/api/comment/create", server.CreateNewComment)
	router.GET("/api/comment/getByPostID/:id", server.GetCommentsByPostID)
	router.DELETE("/api/comment/delete/:id", server.DeleteComment)

	server.Router = router
}

/* StartServer starts the server */
func (server *Server) StartServer(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
