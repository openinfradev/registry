package config

import (
	"builder/repository"
	"builder/service"
	"builder/util/logger"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Configuration is builder configure yaml
type Configuration struct {
	Default  *Default  `yaml:"default"`
	Database *Database `yaml:"database"`
	Registry *Registry `yaml:"registry"`
	Redis    *Redis    `yaml:"redis"`
	Clair    *Clair    `yaml:"clair"`
}

// Print is printing log Configuration values
func (c *Configuration) Print() {
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Default\n domain[%s]\n port[%s]\n tmp[%s]\n loglevel[%s]", c.Default.Domain, c.Default.Port, c.Default.TmpDir, c.Default.LogLevel))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Database\n type[%s]\n host[%s]\n port[%s]\n user[%s]\n password[%s]\n name[%s]\n xargs[%s]", c.Database.Type, c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.Xargs))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Registry\n name[%s]\n insecure[%v]\n endpoint[%s]\n auth[%s]", c.Registry.Name, c.Registry.Insecure, c.Registry.Endpoint, c.Registry.Auth))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Redis\n endpoint[%s]", c.Redis.Endpoint))
	logger.INFO("config/config.go", "Configuration", fmt.Sprintf("Clair\n endpoint[%s]", c.Clair.Endpoint))
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

// LoadConfig returns basicinfo & dbinfo
func LoadConfig() (*service.BasicInfo, *repository.DBInfo) {

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
	conf := &Configuration{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		logger.FATAL("config/config.go", "LoadConfig", err.Error())
	}

	conf.Print()

	// log level
	if loglevel == "" {
		loglevel = conf.Default.LogLevel
	}
	logger.SetLevel(loglevel)

	dbinfo := repository.DBInfo{
		DBtype: conf.Database.Type,
		DBhost: conf.Database.Host,
		DBport: conf.Database.Port,
		DBuser: conf.Database.User,
		DBpass: conf.Database.Password,
		DBname: conf.Database.Name,
		DBxarg: conf.Database.Xargs,
	}

	basicinfo := service.BasicInfo{
		RegistryName:     conf.Registry.Name,
		RegistryInsecure: conf.Registry.Insecure,
		RegistryEndpoint: conf.Registry.Endpoint,
		TemporaryPath:    conf.Default.TmpDir,
		RedisEndpoint:    conf.Redis.Endpoint,
		ClairEndpoint:    conf.Clair.Endpoint,
		AuthURL:          conf.Registry.Auth,
		ServiceDomain:    conf.Default.Domain,
		ServicePort:      conf.Default.Port,
	}
	return &basicinfo, &dbinfo
}
