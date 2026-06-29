package app

import (
	"context"
	"log"
	"net/http"

	"github.com/Rowkash/go-gin-auth/internal/auth"
	"github.com/Rowkash/go-gin-auth/internal/common/middleware"
	"github.com/Rowkash/go-gin-auth/internal/config"
	"github.com/Rowkash/go-gin-auth/internal/database"
	"github.com/Rowkash/go-gin-auth/internal/sessions"
	"github.com/Rowkash/go-gin-auth/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	server *http.Server
	gormDB *gorm.DB
	rdb    *redis.Client
}

func NewApp(cfg config.Config) *App {
	gormDB, err := database.NewPostgresConnection(database.Config{URL: cfg.Database.URL})
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	rdb, err := database.NewRedisClient()
	if err != nil {
		log.Fatalf("Redis client initialization failed: %v", err)
	}

	sessionsModule := sessions.NewModule(rdb)
	usersModule := users.NewModule(gormDB)
	authModule := auth.NewModule(auth.ModuleDeps{
		JwtCfg:          cfg.JWT,
		UsersService:    usersModule.Service(),
		SessionsService: sessionsModule.Service(),
	})
	r := gin.Default()
	r.Use(middleware.ErrorHandler())
	api := r.Group("")
	authModule.RegisterRoutes(api)
	usersModule.RegisterRoutes(api, authModule.Middleware())
	return &App{
		server: &http.Server{
			Handler: r,
		},
		gormDB: gormDB,
		rdb:    rdb,
	}
}

func (a *App) Run(addr string) error {
	a.server.Addr = addr
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Closing background resources...")
	if err := a.rdb.Close(); err != nil {
		log.Printf("Warning: Redis client closing failed: %v", err)
	}

	if a.gormDB != nil {
		sqlDB, err := a.gormDB.DB()
		if err != nil {
			log.Printf("Warning: Failed to get sql.DB for closure: %v", err)
		} else {
			log.Println("Closing Postgres connection pool...")
			if err := sqlDB.Close(); err != nil {
				log.Printf("Warning: Postgres connection pool closing failed: %v", err)
			}
		}
	}

	log.Println("Server exited gracefully")
	return nil
}
