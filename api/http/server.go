package http

import (
	"context"
	"net/http"
	"time"

	"github.com/hiroyaonoe/bcop-proxy-controller/api/http/controller"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	echo  *echo.Echo
	proxy *controller.Proxy
	route *controller.Route
}

func NewServer(proxy *controller.Proxy, route *controller.Route) *Server {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("Method", v.Method).
				Str("URI", v.URI).
				Int("Status", v.Status).
				Stringer("Latency", v.Latency).
				Msg("Request")

			return nil
		},
	}))
	return &Server{
		echo:  e,
		proxy: proxy,
		route: route,
	}
}

func (s *Server) SetRoute() {
	proxy := s.echo.Group("/proxy/:proxy-id")

	proxy.PUT("/register", s.proxy.Register)
	proxy.PUT("/activate", s.proxy.Activate)
	proxy.DELETE("", s.proxy.Delete)

	routes := s.echo.Group("/routes")
	routes.PUT("", s.route.Put)
}

func (s *Server) Run(address string) {
	if err := s.echo.Start(address); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("shutting down the server")
	}
}

func (s *Server) Close() {
	log.Info().Msg("server shutdown")
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown")
	}
}
