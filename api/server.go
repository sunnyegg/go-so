package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sunnyegg/go-so/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// auth
	router.POST("/auth/login", server.loginUser)

	// users
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
