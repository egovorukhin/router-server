package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"path/filepath"
	"time"

	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egolog"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/gofiber/fiber/v3/middleware/static"
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
	Buffer      Buffer       `yaml:"buffer"`
	BodyLimit   int          `yaml:"bodyLimit"`
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

type Buffer struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
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
	AllowPrivateIPs           bool `yaml:"allowPrivateIPs"`
	TlsConfig                 *struct {
		InsecureSkipVerify bool `yaml:"insecureSkipVerify"`
	} `yaml:"tls"`
}

var (
	byteHeaderAccept   = []byte(fiber.HeaderAccept)
	byteMIMETextHTML   = []byte(fiber.MIMETextHTML)
	ipClientHeaderName = "x-client-ip"
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

	// В v3 прокси по умолчанию блокирует loopback/приватные адреса (защита от SSRF).
	// Бэкенды роутера задаются администратором в конфиге, поэтому это можно разрешить.
	policy := proxy.DefaultSecurityPolicy()
	policy.AllowPrivateIPs = s.Client.AllowPrivateIPs
	proxy.WithSecurityPolicy(policy)

	app := fiber.New(fiber.Config{
		AppName:         appName,
		ReadTimeout:     time.Duration(s.Timeout.Read) * time.Second,
		WriteTimeout:    time.Duration(s.Timeout.Write) * time.Second,
		IdleTimeout:     time.Duration(s.Timeout.Idle) * time.Second,
		ReadBufferSize:  s.Buffer.Read,
		WriteBufferSize: s.Buffer.Write,
		BodyLimit:       s.BodyLimit,
	})
	if len(s.Router.Index) == 0 {
		s.Router.Index = "index.html"
	}
	if len(s.Router.IpClientHeaderName) == 0 {
		s.Router.IpClientHeaderName = ipClientHeaderName
	}
	app.Get("/*", static.New(s.Root, static.Config{IndexNames: []string{s.Router.Index}}))
	if s.Logger != nil {
		app.Use(logger.New(logger.Config{
			Format:       s.Logger.Format,
			TimeFormat:   s.Logger.Time.Format,
			TimeZone:     s.Logger.Time.Zone,
			TimeInterval: time.Duration(s.Logger.Time.Interval) * time.Millisecond,
			Stream:       s.Logger,
		}))
	}
	if len(s.Heartbeat.Path) == 0 {
		s.Heartbeat.Path = "/heartbeat"
	}
	if s.Heartbeat.StatusCode == 0 {
		s.Heartbeat.StatusCode = 200
	}
	app.Get(s.Heartbeat.Path, func(ctx fiber.Ctx) error {
		return ctx.Status(s.Heartbeat.StatusCode).SendString("OK")
	})
	app.Use(func(c fiber.Ctx) error {
		url := string(c.Request().RequestURI())
		location, err := s.Router.GetLocation(url)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("%s маршрут не найден. проверьте настройки", url))
		}
		req := c.Request()
		var isTextHTML, isXClientIP bool
		for key, value := range req.Header.All() {
			if bytes.Equal(key, byteHeaderAccept) {
				isTextHTML = HasHeader(value, byteMIMETextHTML)
				continue
			}
			isXClientIP = bytes.Equal(bytes.ToLower(key), []byte(ipClientHeaderName))
		}
		if isTextHTML {
			egolog.Info("isTextHTML:", isTextHTML)
			egolog.Info("%v", c.String())
			return c.Next()
		}
		if !isXClientIP {
			c.Request().Header.Add(s.Router.IpClientHeaderName, c.IP())
		}
		c.Request().Header.Add("Router-Server", "v"+info.GetVersion().String())
		if err = proxy.Do(c, location, cli); err != nil {
			return err
		}
		resp := c.Response()
		fmt.Println(resp.StatusCode())
		return nil
	})
	// SPA fallback: отдаём index для всех неразрешённых маршрутов
	app.Use(func(c fiber.Ctx) error {
		return c.SendFile(filepath.Join(s.Root, s.Router.Index))
	})

	addr := fmt.Sprintf(":%d", s.Port)

	if s.Secure && s.Certificate != nil {
		return app.Listen(addr, fiber.ListenConfig{CertFile: s.Certificate.Cert, CertKeyFile: s.Certificate.Key})
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
