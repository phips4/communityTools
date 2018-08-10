package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

const defaultConfig = `#Config generated on %v
#MongoDB settings
MongoDB:
  host: "localhost"
  port: 27017
  user: "ComToolsUser"
  database: communityTools
  password: "312surlaW"
  setupCollections: true

#Swift Openstack settings (ignore it if you are not going to use it)
OpenstackSwift:
  host: "localhost"
  port: 12345
  user: "ComToolsUser"
  password: "312surlaWkcatsnepo"

#Generall settings
#Use the locally filesystem instead of OpenStack Swift?
useFileSystem: true
`

type Config struct{}

type ConfigStruct struct {
	MongoDB struct {
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		User       string `yaml:"user"`
		Database       string `yaml:"database"`
		Password   string `yaml:"password"`
		SetupCollections bool `yaml:"setupCollections"`
	} `yaml:"MongoDB"`
	OpenstackSwift struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"OpenstackSwift"`

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
