package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sunnyegg/go-so/cron"
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
	server.registerRoutes()
	server.registerCron()

	return server, nil
}

func (server *Server) registerRoutes() {
	router := gin.Default()

	// auth
	router.GET("/auth/login", server.loginUser)
	router.POST("/auth/refresh", server.refreshUser)
	router.GET("/auth/state", server.createState)

	authRoutes := router.Group("/").Use(authMiddleware(server))

	// users
	authRoutes.GET("/users", server.getUser)

	// streams
	authRoutes.POST("/streams", server.createStream)
	authRoutes.GET("/streams/:id", server.getStream)
	authRoutes.GET("/streams", server.listStream)
	authRoutes.GET("/streams/attendance_members", server.getStreamAttendanceMember)

	// attendance members
	authRoutes.POST("/attendance_members", server.createAttendanceMember)

	// user_configs
	authRoutes.POST("/user_configs", server.createUserConfig)
	authRoutes.GET("/user_configs", server.getUserConfig)

	// twitch
	authRoutes.GET("/twitch/user", server.getTwitchUser)
	authRoutes.POST("/twitch/chat/connect", server.connectChat)
	router.POST("/twitch/eventsub", server.handleEventsub)

	// ws
	router.GET("/ws", server.ws)

	server.router = router
}

func (server *Server) registerCron() {
	cronClient := cron.NewCron()
	cronClient.AddFunc("@every 1h", cron.ValidateToken(context.Background(), server.store, server.config))

	cronClient.Start()
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
