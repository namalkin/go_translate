package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/namalkin/go_translate/pkg/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
		}

		// для проверки авторизации
		translations := api.Group("/translations", h.userIdentity)
		{
			translations.POST("/", h.createTranslation)
			translations.GET("/", h.getAllTranslations)
			translations.GET("/:id", h.getTranslationById)
			translations.PUT("/:id", h.updateTranslation)
			translations.DELETE("/:id", h.deleteTranslation)
			translations.POST("/post100", h.createTestTranslations)
			translations.POST("/post100k", h.createTestTranslations100k)
			translations.DELETE("/delete_all", h.deleteAllTranslations)
		}
	}

	return router
}
