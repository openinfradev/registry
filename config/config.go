package config

import (
	"fmt"
	"github.com/openinfradev/registry-builder/util/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Configuration is builder configure yaml
type Configuration struct {
	Default  *Default  `yaml:"default"`
	Database *Database `yaml:"database"`
	Registry *Registry `yaml:"registry"`
	Webhook  *Webhook  `yaml:"webhook"`
	Redis    *Redis    `yaml:"redis"`
	Clair    *Clair    `yaml:"clair"`
	Minio    *Minio    `yaml:"minio"`
}

// Print is printing log Configuration values
func (c *Configuration) Print() {
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Default\n domain[%s]\n port[%s]\n tmp[%s]\n loglevel[%s]", c.Default.Domain, c.Default.Port, c.Default.TmpDir, c.Default.LogLevel))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Database\n type[%s]\n host[%s]\n port[%s]\n user[%s]\n password[%s]\n name[%s]\n xargs[%s]", c.Database.Type, c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.Xargs))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Registry\n name[%s]\n insecure[%v]\n endpoint[%s]\n auth[%s]", c.Registry.Name, c.Registry.Insecure, c.Registry.Endpoint, c.Registry.Auth))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Redis\n endpoint[%s]", c.Redis.Endpoint))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Clair\n endpoint[%s]", c.Clair.Endpoint))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Minio\n data[%s]\n domain[%s]\n ports[from %v to %v]", c.Minio.Data, c.Minio.Domain, c.Minio.StartOfPort, c.Minio.EndOfPort))
}

// Database is db config
type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Xargs    string `yaml:"xargs"`
}

// Registry is registry config
type Registry struct {
	Name     string `yaml:"name"`
	Insecure bool   `yaml:"insecure"`
	Endpoint string `yaml:"endpoint"`
	Auth     string `yaml:"auth"`
}

// Webhook is listeners
type Webhook struct {
	Listener []string `yaml:"listener"`
}

// Default is default config
type Default struct {
	Domain   string `yaml:"domain"`
	Port     string `yaml:"port"`
	TmpDir   string `yaml:"tmp"`
	LogLevel string `yaml:"loglevel"`
}

// Redis is redis config
type Redis struct {
	Endpoint string `yaml:"endpoint"`
}

// Clair is Clair config
type Clair struct {
	Endpoint string `yaml:"endpoint"`
}

// Minio is minio config
type Minio struct {
	Domain      string `yaml:"domain"`
	Data        string `yaml:"data"`
	Ports       string `yaml:"ports"`
	StartOfPort int
	EndOfPort   int
}

var config *Configuration

// LoadConfig is parsing configuration
func LoadConfig() {

	configFile := os.Getenv("BUILDER_CONFIG")
	loglevel := os.Getenv("BUILDER_LOG_LEVEL")
	if configFile == "" {
		configFile = "conf/config.local.yml"
	}

	// configFile := flag.String("config", "conf/config.yml", "config.yml file location (optional)")
	file, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.FATAL("config/config.go", "LoadConfig", err.Error())
	}
	config = &Configuration{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		logger.FATAL("config/config.go", "LoadConfig", err.Error())
	}

	// log level
	if loglevel == "" {
		loglevel = config.Default.LogLevel
	}
	logger.SetLevel(loglevel)

	// minio port
	ports := strings.Split(config.Minio.Ports, "-")
	if len(ports) < 2 {
		logger.FATAL("config/config.go", "LoadConfig", "Wrong ports in minio configuration")
	}
	sp, err := strconv.Atoi(ports[0])
	if err != nil {
		logger.FATAL("config/config.go", "LoadConfig", err.Error())
	}
	ep, err := strconv.Atoi(ports[1])
	if err != nil {
		logger.FATAL("config/config.go", "LoadConfig", err.Error())
	}
	config.Minio.StartOfPort = sp
	config.Minio.EndOfPort = ep

	config.Print()
}

// GetConfig returns configuration
func GetConfig() *Configuration {
	return config
}
