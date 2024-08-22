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
	server.registerRoutes()
	server.registerCron()

	return server, nil
}

func (server *Server) registerRoutes() {
	router := gin.Default()

	router.Use(corsMiddleware())

	// auth
	router.GET("/auth/login", server.loginUser)
	router.POST("/auth/refresh", server.refreshUser)
	router.GET("/auth/state", server.createState)
	router.POST("/auth/logout", server.logoutUser)

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
	router.POST("/twitch/eventsub", server.handleEventsub)

	authRoutes.GET("/twitch/channel", server.getChannelInfo)
	authRoutes.GET("/twitch/stream", server.getStreamInfo)
	authRoutes.GET("/twitch/user", server.getTwitchUser)
	authRoutes.GET("/twitch/chat/connect", server.connectChat)
	authRoutes.POST("/twitch/chat/connect", server.connectChat)
	authRoutes.POST("/twitch/chat/message", server.sendChatMessage)
	authRoutes.POST("/twitch/chat/shoutout", server.sendShoutout)

	// ws
	router.GET("/ws/:id", server.ws)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
