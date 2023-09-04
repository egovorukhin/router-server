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
	"regexp"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	Root        string       `yaml:"root"`
	Router      Router       `yaml:"router"`
	Port        int          `yaml:"port"`
	Secure      bool         `yaml:"secure"`
	Certificate *Certificate `yaml:"certificate"`
	Timeout     Timeout      `yaml:"timeout"`
}

type Router struct {
	Pattern  string   `yaml:"pattern"`
	Index    string   `yaml:"index"`
	Location Location `yaml:"location"`
}

type Location map[string]struct {
	Url    string `yaml:"url"`
	Remove bool   `yaml:"remove"`
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

func (r Router) GetLocation(url string) (string, error) {
	rgx, err := regexp.Compile(r.Pattern)
	if err != nil {
		return "", err
	}
	match := rgx.FindAllString(url, -1)
	if len(match) > 0 {
		if location, ok := r.Location[match[0]]; ok {
			if location.Remove {
				url = strings.Replace(url, match[0], location.Url, 1)
			}
			return location.Url + url, nil
		}
	}
	return "", err
}

func (s *Server) Init() error {

	app := fiber.New(fiber.Config{
		AppName:      fmt.Sprintf("%s v%s", info.GetApplicationName(), info.GetVersion()),
		ReadTimeout:  time.Duration(s.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(s.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(s.Timeout.Idle) * time.Second,
	})
	if len(s.Router.Index) == 0 {
		s.Router.Index = "index.html"
	}

	app.Static("/", "./dist", fiber.Static{
		Index: s.Router.Index,
	})
	app.Static("*", s.Root+"/"+s.Router.Index)
	app.Use(func(c *fiber.Ctx) error {
		url := string(c.Request().RequestURI())
		location, err := s.Router.GetLocation(url)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("%s маршрут не найден. проверьте настройки", url))
		}
		return proxy.Do(c, location)
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
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	errChan <- fmt.Errorf("%s", <-c)
}

func init() {
	info.SetName("ALAcall Router Server")
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
