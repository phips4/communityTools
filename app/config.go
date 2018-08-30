package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
# ownServer - should this module run on its own webserver?
#             If not it will run on the default webserver.
#             This option is useful for splitting up every module to its own sub-domain e.g.
# host      - ip or hostname for the module own webserver, ignore if ownServer is false
# port      - port for the module own webserver, ignore if ownServer is false

Polls:
  enabled: true
  ownServer: false
  host: "localhost"
  port: 54321

######################################
#          general settings          #
######################################
# Webserver
# this is the default webserver. It's used if a module
# doesn't run on its own webserver. 
Webserver:
  enabled: true
  host: "localhost"
  port: 4337

#Use the locally filesystem?
useFileSystem: true`

type Config struct{}

type ConfigStruct struct {
	MongoDB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Database string `yaml:"database"`
		Password string `yaml:"password"`
	} `yaml:"MongoDB"`

	//TODO: implement filesystem first
	/*
		OpenstackSwift struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"OpenstackSwift"`
	*/

	// modules
	ModulePoll `yaml:"Polls"`

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
		panic(err)
	}
}
