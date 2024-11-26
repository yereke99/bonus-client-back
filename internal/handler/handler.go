package handler

import (
	"bonus-client-back/config"
	"bonus-client-back/internal/service"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service   *service.Services
	zapLogger *zap.Logger
	appConfig *config.Config
}

func NewHandler(service *service.Services, zapLogger *zap.Logger, appConfig *config.Config) *Handler {
	return &Handler{
		service:   service,
		zapLogger: zapLogger,
		appConfig: appConfig,
	}
}

func (h *Handler) InitHandler() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Length", "Authorization", "X-CSRF-Token", "Content-Type", "Accept", "X-Requested-With", "Bearer", "Authority"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type", "application/json", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return origin == "https://api.worldbonussystem.com" },
	}))

	r.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	// Customer Email code send
	r.POST("/request-otp")
	r.POST("/login")

	// Refresh token
	r.POST("/refresh-token")

	// С токеном клиента
	r.POST("/get-profile")

	r.PATCH("/customer/profile")
	r.GET("/customer/partners")
	r.GET("/customer/transactions")
	r.DELETE("/customer/profile")

	// С токеном торговый точки
	r.GET("/company-asset/profile")

	return r
}
