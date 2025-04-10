package auth

import (
	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

type AuthModule struct {
	Handler *AuthHandler
}

func NewAuthModule(app *app.Apps) *AuthModule {
	userRepository := users.NewUserRepository(app)
	authService := NewAuthService(userRepository)
	AuthHandler := NewAuthHandler(authService, app)
	return &AuthModule{
		Handler: AuthHandler,
	}
}

func (u *AuthModule) Register(app *app.Apps) error {
	app.Log.Info().Msg("Auth Module Initialized")
	return nil
}

func (a *AuthModule) Route(router *gin.RouterGroup, app *app.Apps) {
	authRoutes := router.Group("/v1/auth")
	{
		authRoutes.POST("/login", a.Handler.Login)
		authRoutes.POST("/register", a.Handler.Register)

		authRoutes.Use(middleware.AuthMiddleware(app))
		authRoutes.POST("/logout", a.Handler.Logout)
		authRoutes.POST("/access-token", a.Handler.GenerateAccessToken)
	}
}
