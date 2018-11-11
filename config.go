package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var appConfig *configuration

const (
	DEFAULT_PORT      int    = 8080
	DEFAULT_LOG_LEVEL string = "INFO"
)

type YamlConfiguration struct {
	ServerPort  int    `yaml:"serverPort"`
	LogLevel    string `yaml:"loglevel"`
	ErrorRegex  string `yaml:"errorRegex"`
	PodLabelKey string `yaml:"withLabelKey"`
}

type configuration struct {
	serverPort  int
	logLevel    string
	errorRegex  string
	podLabelKey string
}

func GetLogLevel() log.Level {
	if appConfig == nil {
		return log.WarnLevel
	}
	level, err := log.ParseLevel(appConfig.logLevel)
	if err != nil {
		log.Warnf("Invalid log level(%v), returning WARN: %v\n", appConfig.logLevel, err)
		level = log.WarnLevel
	}
	return level
}

func GetErrorRegex() string {
	if appConfig == nil {
		return ""
	}
	return appConfig.errorRegex
}

func GetPodLabelKey() string {
	if appConfig == nil {
		return ""
	}
	return appConfig.podLabelKey
}

func GetServerPort() int {
	if appConfig == nil {
		return 0
	}
	return appConfig.serverPort
}

func (c *configuration) String() string {
	var msg = "[serverPort: %v, loglevel: %v, errorRegex: %v, podLabelKey: %v]\n"
	return fmt.Sprintf(msg, GetServerPort(), GetLogLevel(), GetErrorRegex(), GetPodLabelKey())
}

func LoadConfiguration() {
	c := &configuration{
		serverPort: DEFAULT_PORT,
		logLevel:   DEFAULT_LOG_LEVEL,
	}
	yc, err := getYamlConfig()
	if err != nil {
		log.Warnf("Error loading the yaml config file: %v\n", err)
	}
	if yc != nil {
		c.withYaml(yc)
	}

	c.loadEnv()
	c.validate()
	appConfig = c
	log.Infof("Configuration Loaded: %s\n", appConfig)
}

func (c *configuration) withYaml(yc *YamlConfiguration) {
	if yc.LogLevel != "" {
		c.logLevel = yc.LogLevel
	}
	if yc.ErrorRegex != "" {
		c.errorRegex = yc.ErrorRegex
	}
	if yc.PodLabelKey != "" {
		c.podLabelKey = yc.PodLabelKey
	}
	if yc.ServerPort != 0 {
		c.serverPort = yc.ServerPort
	}
}

func (c *configuration) validate() {
	if c.serverPort == 0 {
		log.Fatalln("Configuration error: http port must be set!")
	}
	if c.logLevel == "" {
		log.Fatalln("Configuration error: log level must be set!")
	}
}

func (c *configuration) loadEnv() {
	log.Debug("Overriding configuration from environment variables...")
	if os.Getenv("APP_LOG_LEVEL") != "" {
		c.logLevel = os.Getenv("SERVER_LOGLEVEL")
	}

	if os.Getenv("APP_HTTP_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
		if err != nil {
			log.Fatalf("Invalid port specified: %v\n", err)
		}
		c.serverPort = port
	}

	if os.Getenv("APP_ERROR_REGEX") != "" {
		c.errorRegex = os.Getenv("APP_ERROR_REGEX")
	}

	if os.Getenv("APP_WITH_LABEL_KEY") != "" {
		c.podLabelKey = os.Getenv("APP_WITH_LABEL_KEY")
	}
}

func getYamlConfig() (*YamlConfiguration, error) {
	filePath := "./config.yaml"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}
	yamlConfig := &YamlConfiguration{}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read configuration file: %v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal configuration file: %v\n", err)
	}

	return yamlConfig, nil
}
