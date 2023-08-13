package api

import (
	db "github.com/Oabraham1/open-blogger/server/db/sqlc"
	"github.com/gin-gonic/gin"
)

/* Server serves the HTTP Requests */
type Server struct {
	Router    *gin.Engine
	DataStore *db.Store
}
