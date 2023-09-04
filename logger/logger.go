package logger

import "github.com/egovorukhin/egolog"

type Config struct {
	DirPath  string    `yaml:"dir_path"`
	Info     string    `yaml:"info"`
	Error    string    `yaml:"error"`
	Debug    string    `yaml:"debug"`
	Rotation *Rotation `yaml:"rotation,omitempty"`
}

type Rotation struct {
	Size   int    `yaml:"size"`
	Format string `yaml:"format"`
	Path   string `yaml:"path"`
}

func Init() error {
	cfg := egolog.Config{
		FileName: "app",
		Info:     "3",
		Error:    "3|16",
		Rotation: &egolog.Rotation{
			Size:   10240,
			Format: "%name_old",
		},
	}
	return egolog.InitLogger(cfg, nil)
}
