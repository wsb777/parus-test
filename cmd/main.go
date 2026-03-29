package main

import (
	"fmt"
	"net/http"

	"parus-test/internal/config"
	"parus-test/internal/database"
	"parus-test/internal/database/repo"
	"parus-test/internal/logger"
	"parus-test/internal/server"
	"parus-test/internal/service"
	"parus-test/internal/storage"

	_ "parus-test/docs"

	"go.uber.org/zap"
)

// @title           Parus Test
// @version         1.0
// @description     Сервис для публикации и проверки обновлений, с учетом проверки целостности файлов через хеш-суммы

func main() {

	// Initialize config
	cfg := config.Load()

	// Initialize logger

	loggers, err := logger.New(cfg.AppEnv)

	if err != nil {
		panic(fmt.Errorf("logger init failed: %w", err))
	}

	// Initialize database
	db := database.NewDatabasePG(cfg, loggers.DB)
	dbConn, err := db.Connect()
	if err != nil {
		loggers.App.Fatal("Failed to connect database:", zap.Error(err))
	}
	err = db.Migrate(dbConn)
	if err != nil {
		loggers.App.Fatal("Failed to migrate database:", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repo.NewUserRepo(dbConn)
	groupRepo := repo.NewGroupRepo(dbConn)
	fileRepo := repo.NewFileRepo(dbConn)
	tokenRepo := repo.NewTokenRepo(dbConn)

	// Initialize storage
	fileStorage := &storage.LocalDiskStorage{
		BasePath: "/uploads",
	}

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenRepo, loggers.Service)
	userService := service.NewUserService(userRepo, loggers.Service)
	groupService := service.NewGroupService(groupRepo, loggers.Service)
	fileService := service.NewFileService(fileRepo, fileStorage, loggers.Service)

	// Create admin group and admin user

	if groupID, err := groupService.CreateAdminGroup(); err == nil {
		userService.CreateAdmin(cfg.AdminName, cfg.AdminPass, groupID)
	}

	handler, err := server.InitHttpServer(authService, groupService, userService, fileService, loggers.HTTP)
	if err != nil {
		loggers.App.Fatal("init failed", zap.Error(err))
	}

	loggers.App.Info("server started", zap.String("port", cfg.Port))

	// TODO Из за TLS ломаются логи, в идеале написать адептер, чтобы логировался в json для ELK
	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		if err := http.ListenAndServeTLS(
			fmt.Sprintf(":%s", cfg.Port),
			cfg.TLSCert,
			cfg.TLSKey,
			handler,
		); err != nil {
			loggers.App.Fatal("server crashed", zap.Error(err))
		}
	} else {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), handler); err != nil {
			loggers.App.Fatal("server crashed", zap.Error(err))
		}
	}
}
