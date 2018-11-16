package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const defaultConfig = `#Config generated on %v
######################################
#            connections             #
######################################
MongoDB:
  host: "localhost"
  port: 27017
  user: "ComToolsUser"
  database: communityTools
  password: "312surlaW"

######################################
#             modules                #
######################################
# enabled   - de/activate the given module
Polls:
  enabled: true

######################################
#          general settings          #
######################################
# Webserver
# this is the default webserver. 
Webserver:
  host: "localhost"
  port: 4337

#Use the locally filesystem?
useFileSystem: true`

type Config struct{}

type MgoConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
}

type ConfigStruct struct {
	MgoConf `yaml:"MongoDB"`

	Webserver struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
	} `yaml:"Webserver"`

	UseFileSystem string `yaml:"useFileSystem"`
}

func (c *Config) LoadConfig() *ConfigStruct {
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		//config file does not exist, creating default config
		file, err := os.Create("config.yml")
		must(err)
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf(defaultConfig, time.Now().Format("Mon Jan 2 15:04:05 2006")))
		must(err)
	}

	bytes, err := ioutil.ReadFile("config.yml")
	must(err)
	conf := &ConfigStruct{}
	err = yaml.Unmarshal(bytes, conf)
	must(err)

	return conf
}

func must(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}
