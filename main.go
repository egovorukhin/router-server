package main

import (
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egoconf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Root        string       `yaml:"root"`
	Location    string       `yaml:"location"`
	Port        int          `yaml:"port"`
	Secure      bool         `yaml:"secure"`
	Certificate *Certificate `yaml:"certificate"`
	Timeout     Timeout      `yaml:"timeout"`
}

type Certificate struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type Timeout struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

func (s *Server) Init() error {

	app := fiber.New(fiber.Config{
		AppName:      fmt.Sprintf("%s v%s", info.GetApplicationName(), info.GetVersion()),
		ReadTimeout:  time.Duration(s.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(s.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(s.Timeout.Idle) * time.Second,
	})

	app.Static("/", s.Root)
	app.Use(func(c *fiber.Ctx) error {
		return proxy.Do(c, s.Location+string(c.Request().RequestURI()))
	})

	addr := fmt.Sprintf(":%d", s.Port)

	if s.Secure && s.Certificate != nil {
		return app.ListenTLS(addr, s.Certificate.Cert, s.Certificate.Key)
	}

	return app.Listen(addr)
}

func start(errChan chan error) {

	// Загружаем конфигурацию приложения
	server := &Server{}
	err := egoconf.Load("config.yml", server)
	if err != nil {
		errChan <- err
		return
	}

	// Запускаем сервер
	errChan <- server.Init()
}

func waitSignal(errChan chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL) //nolint:govet
	errChan <- fmt.Errorf("%s", <-c)
}

func init() {
	info.SetName("ALAcall Call Center Server")
	info.SetVersion(0, 0, 1)
}

func main() {

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan error, 2)
	// Ждем сигнал от ОС
	go waitSignal(errChan)
	// Стартовая горутина
	go start(errChan)

	if err := <-errChan; err != nil {
		log.Println(err)
	}
}
