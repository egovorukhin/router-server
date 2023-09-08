package server

import (
	"fmt"
	"github.com/egovorukhin/egolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"net/http"
	"time"
)

type Config struct {
	Root        string       `yaml:"root"`
	Router      Router       `yaml:"router"`
	Port        int          `yaml:"port"`
	Secure      bool         `yaml:"secure"`
	Certificate *Certificate `yaml:"certificate"`
	Logger      *Logger      `yaml:"logger,omitempty"`
	Timeout     Timeout      `yaml:"timeout"`
}

type Timeout struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

func Init(s *Config, appName string) error {

	app := fiber.New(fiber.Config{
		AppName:      appName,
		ReadTimeout:  time.Duration(s.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(s.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(s.Timeout.Idle) * time.Second,
	})
	if len(s.Router.Index) == 0 {
		s.Router.Index = "index.html"
	}
	app.Static("/", s.Root, fiber.Static{
		Index: s.Router.Index,
	})
	//app.Static("*", s.Root+"/"+s.Router.Index)
	if s.Logger != nil {
		app.Use(logger.New(logger.Config{
			Format:       s.Logger.Format,
			TimeFormat:   s.Logger.Time.Format,
			TimeZone:     s.Logger.Time.Zone,
			TimeInterval: time.Duration(s.Logger.Time.Interval) * time.Millisecond,
			Output:       s.Logger,
		}))
	}
	app.Use(func(c *fiber.Ctx) error {
		url := string(c.Request().RequestURI())
		location, err := s.Router.GetLocation(url)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("%s маршрут не найден. проверьте настройки", url))
		}
		egolog.Info(c.Response().StatusCode())
		if err = proxy.Do(c, location); err != nil {
			return c.Next()
		}
		if c.Response().StatusCode() == fiber.StatusNotFound {
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
