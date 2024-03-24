package apiServer

import (
	"github.com/damonto/estkme-rlpa-server/internal/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewAPIServer() *gin.Engine {
	r := gin.Default()

	setupServer()

	store := cookie.NewStore([]byte(config.C.Secret))
	r.Use(sessions.Sessions("rlpa", store))

	api := r.Group("/api")
	api.POST("/login", login)
	api.GET("/logout", logout)
	api.GET("/registerEnableCheck", registerEnabledCheck)
	api.POST("/register", register)

	loginRequired := r.Group("/api")
	loginRequired.Use(AuthRequired)
	{
		loginRequired.GET("/socket", socketHandler)
	}

	return r
}
