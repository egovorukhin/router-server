package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ucc-front-server/logger"
	"ucc-front-server/server"

	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egoconf"
	"github.com/egovorukhin/egolog"
)

type Config struct {
	server.Config `yaml:",inline"`
}

func init() {
	info.SetName("Router Server")
	info.SetVersion(0, 0, 7)
}

func main() {

	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}

	app := fmt.Sprintf("%s v%s", info.GetApplicationName(), info.GetVersion())

	egolog.Info("start %s", app)
	defer egolog.Info("stop")

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan error, 2)
	// Ждем сигнал от ОС
	go waitSignal(errChan)
	// Стартовая горутина
	go start(errChan, app)

	if err := <-errChan; err != nil {
		egolog.Error(err)
	}
}

func start(errChan chan error, appName string) {

	// Загружаем конфигурацию приложения
	cfg := Config{}
	err := egoconf.Load("config.yml", &cfg)
	if err != nil {
		errChan <- err
		return
	}

	// Запускаем сервер
	errChan <- server.Init(&cfg.Config, appName)
}

func waitSignal(errChan chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	errChan <- fmt.Errorf("%s", <-c)
}
