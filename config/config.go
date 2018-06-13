package config

import (
	"github.com/BurntSushi/toml"
)

type sqs struct {
	QueueUrl string `toml:"queue_url"`
}

type worker struct {
	Url     string `toml:"url"`
	Workers int    `toml:"workers"`
}

type Config struct {
	Sqs    sqs    `toml:"sqs"`
	Worker worker `toml:"worker"`
}

func LoadConfig(path string) (Config, error) {
	var c Config
	if _, err := toml.DecodeFile(path, &c); err != nil {
		return c, err
	}

	return c, nil
}
