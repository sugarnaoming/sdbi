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

type SdSetting struct {
	APIURL    string `yaml:"api-url"`
	UserToken string `yaml:"user-token"`
	UIURL     string `yaml:"ui-url"`
}

// Empty string is available for “CurrentConfName"
type Config struct {
	Config          map[string]SdSetting `yaml:"config"`
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
	UIURL      string
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
		conf.Config = make(map[string]SdSetting)
	}

	conf.Config[blueprint.ConfigName] = SdSetting{
		APIURL:    blueprint.APIURL,
		UserToken: blueprint.Token,
		UIURL:     blueprint.UIURL,
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

func (conf *Config) Set(blueprint Blueprint) error {
	newConf, ok := conf.Config[blueprint.ConfigName]
	if !ok {
		return fmt.Errorf("%s config dose not exist in %s", blueprint.ConfigName, configFileName)
	}

	if blueprint.APIURL != "" {
		newConf.APIURL = blueprint.APIURL
	}
	if blueprint.Token != "" {
		newConf.UserToken = blueprint.Token
	}
	if blueprint.UIURL != "" {
		newConf.UIURL = blueprint.UIURL
	}

	conf.Config[blueprint.ConfigName] = newConf

	err := saveConfs(conf)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (conf *Config) Delete(configName string) error {
	_, ok := conf.Config[configName]
	if !ok {
		return fmt.Errorf("can not find %s from config.yaml", configName)
	}
	if 1 >= len(conf.Config) {
		return fmt.Errorf("failed to delete: can not delete all configs")
	}

	if conf.CurrentConfName == configName {
		conf.CurrentConfName = ""
	}
	delete(conf.Config, configName)

	err := saveConfs(conf)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (conf *Config) Use(configName string) error {
	if _, ok := conf.Config[configName]; !ok {
		return fmt.Errorf("%s config dose not exist in %s", configName, configFileName)
	}
	conf.CurrentConfName = configName

	err := saveConfs(conf)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	return nil
}

func (conf *Config) CurrentConfig() (SdSetting, error) {
	if conf.CurrentConfName == "" {
		return SdSetting{}, fmt.Errorf("config is not used")
	}
	return conf.Config[conf.CurrentConfName], nil
}
