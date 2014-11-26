package poddl

import (
	"encoding/json"
	"fmt"
	"github.com/adrg/xdg"
	"io/ioutil"
)

type Config struct {
	Address  string
	User     string
	Domain   string
	Password string
	Feeds    []string
}

const configFilename = "poddl/config.json"

var requiredKeys = []string{"address", "user", "domain", "password"}

func readConfig() (c *Config, err error) {
	filepath, err := xdg.ConfigFile(configFilename)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
