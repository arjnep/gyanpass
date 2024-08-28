package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arjnep/gyanpass/config"
	"github.com/arjnep/gyanpass/internal/db"
	httpUser "github.com/arjnep/gyanpass/internal/delivery/http/user"
	"github.com/arjnep/gyanpass/internal/delivery/middleware"
	"github.com/arjnep/gyanpass/internal/repository"
	"github.com/arjnep/gyanpass/internal/usecase"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadConfig()
	db.SetupPostgres()
}

func main() {
	router := gin.New()
	err := os.MkdirAll("log", 0755)
	if err != nil {
		log.Fatalf("Error Creating log directory: %v", err)
		return
	}

	logFile, err := os.Create("log/api.log")
	if err != nil {
		log.Fatalf("Error Creating log file: %v", err)
	}

	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(logFile)

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery())
	// router.Use(middleware.CORS())
	router.NoRoute(middleware.NoRouteHandler())

	database := db.GetDB()
	cfg := config.GetConfig()

	userRepo := repository.NewUserRepository(database)
	// bookRepo := repository.NewBookRepository(database)
	// exchangeRepo := repository.NewExchangeRepository(database)

	jwtService := jwt.NewJWTService(cfg)
	// notificationService := notification.NewNotificationService()

	userUsecase := usecase.NewUserUsecase(userRepo, jwtService)
	// bookUsecase := usecase.NewBookUsecase(bookRepo)
	// exchangeUsecase := usecase.NewExchangeUsecase(exchangeRepo, userRepo, bookRepo, notificationService)

	httpUser.NewUserHandler(&httpUser.Config{
		R:           router,
		UserUsecase: userUsecase,
		JwtService:  jwtService,
	})
	// http.NewBookHandler(router, bookUsecase, jwtService)
	// http.NewExchangeHandler(router, exchangeUsecase, bookUsecase, jwtService)

	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Error Running Server: %v", err)
		}
	}()

	fmt.Println("*************************************************************************")
	fmt.Println("GyanPass", cfg.Server.Version, "is live on port ", ":"+cfg.Server.Port)
	fmt.Println("*************************************************************************")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server Shutted Down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.Timeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Forced to Shutdown: %v\n", err)
	}

	<-ctx.Done()
	log.Println("Timeout of 5 seconds...")

	log.Println("Server Exiting...")

}
