package task

import (
	"log"
)

var Config Configuration

type Configuration struct {
	ListDefaultLimit int `toml:"list_default_limit"`
	ListMaxLimit     int `toml:"list_max_limit"`
}

func LoadConfig(cfg *Configuration) error {
	Config = *cfg

	if Config.ListDefaultLimit == 0 {
		log.Fatal("List default limit is required")
	}

	if Config.ListMaxLimit == 0 {
		log.Fatal("List max limit is required")
	}

	return nil
}
