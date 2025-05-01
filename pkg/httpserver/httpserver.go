package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine          *gin.Engine
	server          *http.Server
	notify          chan error
	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

func New(port string) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())

	s := &Server{
		Engine:          engine,
		address:         ":" + port,
		readTimeout:     5 * time.Second,
		writeTimeout:    5 * time.Second,
		shutdownTimeout: 10 * time.Second,
		notify:          make(chan error, 1),
	}

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.Engine,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
