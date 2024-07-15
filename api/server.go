package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sunnyegg/go-so/sqlc"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
}

func NewServer(q *db.Queries) *Server {
	server := &Server{queries: q}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
