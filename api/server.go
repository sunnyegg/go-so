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

	// users
	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUser)

	// streams
	router.POST("/streams", server.createStream)
	router.GET("/streams/:id", server.getStream)
	router.GET("/streams", server.listStream)
	router.GET("/streams/attendance_members", server.getStreamAttendanceMember)

	// attendance members
	router.POST("/attendance_members", server.createAttendanceMember)

	// user_configs
	router.POST("/user_configs", server.createUserConfig)
	router.GET("/user_configs", server.getUserConfig)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
