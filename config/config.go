package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"
const configDirName = "sdbi"

type KeyType uint

type sdSetting struct {
	APIURL    string `yaml:"api-url"`
	UserToken string `yaml:"user-token"`
}

// Empty string is available for “CurrentConfName"
type Config struct {
	Config          map[string]sdSetting `yaml:"config"`
	CurrentConfName string               `yaml:"current"`
	configDirPath   string
	configFilePath  string
}

type Configurator interface {
	Load() (*Config, error)
	Use(configName string) error
}

type Blueprint struct {
	APIURL     string
	Token      string
	ConfigName string
}

func New() (*Config, error) {
	conf := new(Config)
	confDirPath, err := os.UserConfigDir()
	if err != nil {
		return conf, fmt.Errorf("failed to read config directory path: %w", err)
	}
	conf.configDirPath = filepath.Join(confDirPath, configDirName)
	conf.configFilePath = filepath.Join(conf.configDirPath, configFileName)

	_, err = os.Stat(conf.configFilePath)
	if err != nil {
		// Create config file and dir if not exist file or dir
		conf, err := conf.initConfig()
		if err != nil {
			return conf, fmt.Errorf("failed to init create config file: %w", err)
		}
	}

	return conf, nil
}

func (conf *Config) initConfig() (*Config, error) {
	err := os.MkdirAll(conf.configDirPath, 0755)
	if err != nil {
		return conf, fmt.Errorf("failed to initialize config: %w", err)
	}

	err = saveConfs(conf)
	if err != nil {
		return conf, fmt.Errorf("failed to initialize config: %w", err)
	}

	return conf, nil
}

func (conf *Config) Load() error {
	buf, err := ioutil.ReadFile(conf.configFilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return nil
}

func (conf *Config) Create(blueprint Blueprint) error {
	if conf.Config != nil {
		_, ok := conf.Config[blueprint.ConfigName]
		if ok {
			return fmt.Errorf("failed to create config: config name is exist")
		}
	} else {
		conf.Config = make(map[string]sdSetting)
	}

	conf.Config[blueprint.ConfigName] = sdSetting{
		APIURL:    blueprint.APIURL,
		UserToken: blueprint.Token,
	}

	err := saveConfs(conf)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	return nil
}

func saveConfs(newConf *Config) error {
	buf, err := yaml.Marshal(newConf)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}
	err = ioutil.WriteFile(newConf.configFilePath, buf, 0744)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
