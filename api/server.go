package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()
	server.registerRoutes(router)

	return server, nil
}

func (server *Server) registerRoutes(router *gin.Engine) {
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
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
