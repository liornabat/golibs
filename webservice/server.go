package webservice

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"fmt"
	"net/http"

	log "golibs/logging"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/sirupsen/logrus"
)

type RouteType int
type Server struct {
	port         string
	IsReady      bool
	wsServer     *gosocketio.Server
	healthFunc   func() (bool, string)
	readyFunc    func() (bool, string)
	isSocketIO   bool
	isPrometheus bool
	routes       map[string]*route
	isCors       bool
	corsConfig   cors.Config
	jwt          *JwtAuth
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type route struct {
	kind RouteType
	path string
	f    func(c *gin.Context)
}

const (
	Undefined RouteType = 0
	GET       RouteType = 1
	POST      RouteType = 2
	PUT       RouteType = 3
	DELETE    RouteType = 4
)

var logger = log.NewLogger("web/server")

func NewServer(p string) *Server {
	s := &Server{
		port:         p,
		routes:       make(map[string]*route),
		corsConfig:   cors.DefaultConfig(),
		readTimeout:  10 * time.Second,
		writeTimeout: 10 * time.Second,
	}
	gin.SetMode(gin.ReleaseMode)

	return s
}

func (s *Server) SetHealthFunc(f func() (bool, string)) *Server {
	s.healthFunc = f
	return s
}

func (s *Server) SetReadyFunc(f func() (bool, string)) *Server {
	s.readyFunc = f
	return s
}

func (s *Server) SetReadTimeout(value time.Duration) *Server {
	s.readTimeout = value
	return s
}

func (s *Server) SetWriteTimeout(value time.Duration) *Server {
	s.writeTimeout = value
	return s
}

func (s *Server) SetSocketIO(set bool) *Server {
	s.isSocketIO = set
	return s
}

func (s *Server) SetPrometheus(set bool) *Server {
	s.isPrometheus = set
	return s
}

func (s *Server) AddRoute(kind RouteType, path string, f func(c *gin.Context)) *Server {
	s.routes[path] = &route{kind: kind, path: path, f: f}
	return s
}

func (s *Server) SetGinDebug() *Server {
	gin.SetMode(gin.DebugMode)
	return s
}

func (s *Server) SetCorsAllowAll() *Server {
	s.isCors = true
	s.corsConfig = cors.DefaultConfig()
	s.corsConfig.AllowAllOrigins = true
	return s
}

func (s *Server) SetCorsAllowOrigins(origins []string) *Server {
	s.isCors = true
	s.corsConfig.AllowOrigins = origins
	return s
}

func (s *Server) SetCorsAllowMethods(methods []string) *Server {
	s.isCors = true
	s.corsConfig.AllowMethods = methods
	return s
}

func (s *Server) SetCorsAllowHeaders(headers []string) *Server {
	s.isCors = true
	s.corsConfig.AllowHeaders = headers
	return s
}

func (s *Server) SetCorsExposeHeaders(headers []string) *Server {
	s.isCors = true
	s.corsConfig.ExposeHeaders = headers
	return s
}

func (s *Server) SetCorsAllowedCredentials() *Server {
	s.isCors = true
	s.corsConfig.AllowCredentials = true
	return s
}

func (s *Server) SetCorsMaxAge(duration time.Duration) *Server {
	s.isCors = true
	s.corsConfig.MaxAge = duration
	return s
}

func (s *Server) SetJwtAuth(j *JwtAuth) *Server {
	s.jwt = j
	return s
}

func (s *Server) GetIdByJwtClaim(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	return claims["id"].(string)
}

func (s *Server) Run() {
	s.startHttpServer()
}

func (s *Server) startHttpServer() {

	router := gin.Default()

	if s.isCors {
		if err := s.corsConfig.Validate(); err != nil {
			logger.Error(err, "validation error for CORS config. No CORS support will be available")
		} else {
			router.Use(cors.New(s.corsConfig))
			logger.Info("CORS support available")
		}
	} else {
		logger.Info("CORS support is not set.No CORS support will be available")
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hi")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/health",
		func(c *gin.Context) {
			if s.healthFunc != nil {
				if ok, msg := s.healthFunc(); !ok {
					c.String(http.StatusServiceUnavailable, msg)
					return
				}
			}
			c.String(http.StatusOK, "All services are healthy")
		},
	)

	router.GET("/ready",
		func(c *gin.Context) {
			if s.readyFunc != nil {
				if ok, msg := s.readyFunc(); !ok {
					c.String(http.StatusServiceUnavailable, msg)
					return
				}
			}
			c.String(http.StatusOK, "Ready ! :)")
		},
	)

	for _, value := range s.routes {
		switch value.kind {
		case GET:
			router.GET(value.path, value.f)
		case POST:
			router.POST(value.path, value.f)
		case PUT:
			router.PUT(value.path, value.f)
		case DELETE:
			router.DELETE(value.path, value.f)
		}
	}

	if s.isSocketIO {
		s.wsServer = s.getWebSocketServer()
		logrus.AddHook(s)
		router.Any("/socket.io/", gin.WrapH(s.wsServer))
	}

	if s.isPrometheus {
		router.Any("/metrics", gin.WrapH(promhttp.Handler()))
	}

	if s.jwt != nil {
		err := s.jwt.set(router)
		if err != nil {
			logger.Error(err, "error during set of jwt middleware")
		} else {
			logger.Info("jwt middleware was set")
		}

	}
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%s", s.port),
		Handler:        router,
		ReadTimeout:    s.readTimeout,
		WriteTimeout:   s.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go httpServer.ListenAndServe()

}

func (s *Server) sendLogs(msg string) {
	s.wsServer.BroadcastTo("main", "/logs", msg)
}
func (s *Server) getWebSocketServer() *gosocketio.Server {

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		logger.Info(fmt.Sprintf("Client from %s has connected.", c.Ip()))
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		logger.Info(fmt.Sprintf("Client from %s has disconnected", c.Ip()))
	})
	server.On("/join", func(c *gosocketio.Channel, room string) string {
		c.Join(room)
		logger.Info(fmt.Sprintf("Client from %s has joined to %s room", c.Ip(), room))
		return "joined"
	})

	return server
}

func (s *Server) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (s *Server) Fire(entry *logrus.Entry) error {
	var fields string
	for key, value := range entry.Data {
		fields = fmt.Sprintf("%s %s:%s ", fields, key, value)
	}
	msg := fmt.Sprintf("[%s] [%s] (%s) %s", entry.Time.Format("2006-01-02 15:04:05.999"), fields, entry.Level.String(), entry.Message)

	if s.isSocketIO {
		s.wsServer.BroadcastTo("main", "/logs", msg)
	}

	return nil
}
