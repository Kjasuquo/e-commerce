package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/decadevs/shoparena/database"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/services"

	"github.com/decadevs/shoparena/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	DB     database.PostgresDb
	Router *router.Router
}

func Start() error {

	values := database.InitDBParams()

	//Setting up the Postgres Database
	var PDB = new(database.PostgresDb)
	var Mail = new(services.Service)
	var Paystack = services.NewPaystack()
	h := &handlers.Handler{DB: PDB, Mail: Mail, Paystack: Paystack}
	err := PDB.Init(values.Host, values.User, values.Password, values.DbName, values.Port)
	if err != nil {
		log.Println("Error trying to Init", err)
		return err
	}

	route, port := router.SetupRouter(h)
	fmt.Println("connected on port ", port)
	err = route.Run(port)
	if err != nil {
		log.Printf("Error from SetupRouter :%v", err)
		return err
	}

	return nil
}

func (s *Server) defineRoutes(router *gin.Engine) {

}

func (s *Server) setupRouter() *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "test" {
		r := gin.New()
		s.defineRoutes(r)
		return r
	}
	r := gin.New()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	// setup cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//https://oja-ecommerce.herokuapp.com
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://oja-ecommerce.herokuapp.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	s.defineRoutes(r)
	return r
}

func (s *Server) Start() {
	r := s.setupRouter()
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if PORT == ":" {
		PORT = ":8080"
	}
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on %s\n", PORT)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
