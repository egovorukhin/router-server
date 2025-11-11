package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
)

type Config struct {
	Root        string       `yaml:"root"`
	Router      Router       `yaml:"router"`
	Port        int          `yaml:"port"`
	Secure      bool         `yaml:"secure"`
	Certificate *Certificate `yaml:"certificate"`
	Logger      *Logger      `yaml:"logger,omitempty"`
	Timeout     Timeout      `yaml:"timeout"`
	Heartbeat   Heartbeat    `yaml:"heartbeat"`
	Client      Client       `yaml:"client"`
}

type Certificate struct {
	Name       string `yaml:"name"`
	Cert       string `yaml:"cert"`
	Key        string `yaml:"key"`
	ClientCert string `yaml:"clientCert"`
}

type Timeout struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

type Heartbeat struct {
	Path       string `yaml:"path"`
	StatusCode int    `yaml:"statusCode"`
}

type Client struct {
	MaxConnsPerHost           int  `yaml:"maxConnsPerHost"`
	MaxIdemponentCallAttempts int  `yaml:"maxIdemponentCallAttempts"`
	MaxResponseBodySize       int  `yaml:"maxResponseBodySize"`
	ReadBufferSize            int  `yaml:"readBufferSize"`
	WriteBufferSize           int  `yaml:"writeBufferSize"`
	ReadTimeout               int  `yaml:"readTimeout"`
	WriteTimeout              int  `yaml:"writeTimeout"`
	MaxConnWaitTimeout        int  `yaml:"maxConnWaitTimeout"`
	StreamResponseBody        bool `yaml:"streamResponseBody"`
	TlsConfig                 *struct {
		InsecureSkipVerify bool `yaml:"insecureSkipVerify"`
	} `yaml:"tls"`
}

var (
	byteHeaderAccept = []byte(fiber.HeaderAccept)
	byteMIMETextHTML = []byte(fiber.MIMETextHTML)
)

func Init(s *Config, appName string) error {

	cli := &fasthttp.Client{
		NoDefaultUserAgentHeader:  true,
		DisablePathNormalizing:    true,
		MaxConnsPerHost:           s.Client.MaxConnsPerHost,
		MaxIdemponentCallAttempts: s.Client.MaxIdemponentCallAttempts,
		ReadBufferSize:            s.Client.ReadBufferSize,
		WriteBufferSize:           s.Client.WriteBufferSize,
		ReadTimeout:               time.Duration(s.Client.ReadTimeout) * time.Second,
		WriteTimeout:              time.Duration(s.Client.WriteTimeout) * time.Second,
		MaxResponseBodySize:       s.Client.MaxResponseBodySize,
		MaxConnWaitTimeout:        time.Duration(s.Client.MaxConnWaitTimeout) * time.Second,
		StreamResponseBody:        s.Client.StreamResponseBody,
	}
	if s.Client.TlsConfig != nil {
		cli.TLSConfig = &tls.Config{
			InsecureSkipVerify: s.Client.TlsConfig.InsecureSkipVerify,
		}
	}

	app := fiber.New(fiber.Config{
		AppName:      appName,
		ReadTimeout:  time.Duration(s.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(s.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(s.Timeout.Idle) * time.Second,
	})
	if len(s.Router.Index) == 0 {
		s.Router.Index = "index.html"
	}
	if len(s.Router.IpClientHeaderName) == 0 {
		s.Router.IpClientHeaderName = "X-Client-IP"
	}
	app.Static("/", s.Root, fiber.Static{
		Index: s.Router.Index,
	})
	if s.Logger != nil {
		app.Use(logger.New(logger.Config{
			Format:       s.Logger.Format,
			TimeFormat:   s.Logger.Time.Format,
			TimeZone:     s.Logger.Time.Zone,
			TimeInterval: time.Duration(s.Logger.Time.Interval) * time.Millisecond,
			Output:       s.Logger,
		}))
	}
	if len(s.Heartbeat.Path) == 0 {
		s.Heartbeat.Path = "/heartbeat"
	}
	if s.Heartbeat.StatusCode == 0 {
		s.Heartbeat.StatusCode = 200
	}
	app.Get(s.Heartbeat.Path, func(ctx *fiber.Ctx) error {
		return ctx.Status(s.Heartbeat.StatusCode).SendString("OK")
	})
	app.Use(func(c *fiber.Ctx) error {
		url := string(c.Request().RequestURI())
		location, err := s.Router.GetLocation(url)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("%s маршрут не найден. проверьте настройки", url))
		}
		req := c.Request()
		var isTextHTML bool
		req.Header.VisitAll(func(key, value []byte) {
			if bytes.Equal(key, byteHeaderAccept) {
				isTextHTML = HasHeader(value, byteMIMETextHTML)
				return
			}
		})
		if isTextHTML {
			return c.Next()
		}
		c.Request().Header.Add(s.Router.IpClientHeaderName, c.IP())
		if err = proxy.Do(c, location, cli); err != nil {
			return err
		}

		resp := c.Response()
		if resp.StatusCode() == fiber.StatusNotFound {
			return c.Next()
		}
		return nil
	})
	app.Use(filesystem.New(filesystem.Config{
		Root:         http.Dir(s.Root),
		PathPrefix:   "/",
		Browse:       false,
		Index:        s.Router.Index,
		NotFoundFile: s.Router.Index,
	}))

	addr := fmt.Sprintf(":%d", s.Port)

	if s.Secure && s.Certificate != nil {
		return app.ListenTLS(addr, s.Certificate.Cert, s.Certificate.Key)
	}

	return app.Listen(addr)
}

func HasHeader(header, value []byte) bool {
	n := bytes.Index(header, value)
	if n < 0 {
		return false
	}
	b := header[n+len(value):]
	if len(b) > 0 && b[0] != ',' {
		return false
	}
	if n == 0 {
		return true
	}
	return header[n-1] == ' '
}
