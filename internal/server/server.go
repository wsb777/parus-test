package server

import (
	"net/http"
	"parus-test/internal/handlers"
	"parus-test/internal/logger"
	"parus-test/internal/middleware"
	"parus-test/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func InitHttpServer(
	authService service.AuthService,
	groupService service.GroupService,
	userService service.UserService,
	fileService service.FileService,

	l *zap.Logger,
) (http.Handler, error) {
	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	// Logging and timeout middleware
	r.Use(logger.HTTPLoggerMiddleware(l))
	r.Use(middleware.TimeoutMiddleware)

	// Initialize handlers
	authHandler := handlers.SignIn(authService)

	groupCreateHandler := handlers.CreateGroup(groupService)
	groupListHandler := handlers.GetGroupList(groupService)

	userCreateHandler := handlers.CreateUser(userService)
	userListHandler := handlers.GetUserList(userService)
	userRevokeTokensHandler := handlers.RevokeUserTokens(authService)

	fileUploadHandler := handlers.UploadFile(fileService)
	fileGetInfoHandler := handlers.GetFileInfo(fileService)
	fileGetLatestVersionInfoHandler := handlers.GetFileLatestVersionInfo(fileService)
	fileGetDataByVersionHandler := handlers.GetFileData(fileService)

	// Swagger route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/", authHandler)
	})

	// User routes
	// r.Route("/users", func(r chi.Router) {
	// 	r.Post("/", userHandler.CreateUser)
	// })

	// Group routes
	// r.Route("/groups", func(r chi.Router) {
	// 	r.Post("/", groupHandler.CreateGroup)
	// })

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))

		// File routes
		r.Route("/files", func(r chi.Router) {
			r.Get("/{file_id}/version/list", fileGetInfoHandler)
			r.Get("/{file_id}/version/latest", fileGetLatestVersionInfoHandler)
			r.Get("/{file_id}/version/{file_version}/download", fileGetDataByVersionHandler)
		})

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminMiddleware)
			r.Get("/users/list", userListHandler)
			r.Post("/users", userCreateHandler)
			r.Post("/users/{user_id}/token/revoke", userRevokeTokensHandler)

			r.Get("/groups/list", groupListHandler)
			r.Post("/groups", groupCreateHandler)

			r.Post("/files", fileUploadHandler)
			r.Post("/files/{file_id}", fileUploadHandler)
		})

	})

	return r, nil
}
