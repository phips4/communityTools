package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const dEFAULT_CONFIG = `#Config generated on %v
######################################
#            connections             #
######################################
MongoDB:
  host: "localhost"
  port: 27017
  user: ""
  database: communityTools
  password: "312surlaW"

######################################
#             modules                #
######################################
# enabled   - de/activates a module
modules:
  urlshortener:
    shortName: "url"
    enabled: true

  textupload:
    shortName: "txt"
    enabled: true

  fileupload: 
    shortName: "file"
    enabled: true

  imageupload: 
    shortName: "img"
    enabled: true

  polls: 
    shortName: "polls"
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

	Modules struct {
		TextUpload struct {
			ShortName string `yaml:"shortName" json:"short_name"`
			Enabled   bool   `yaml:"enabled" json:"enabled"`
		} `yaml:"textupload" json:"text_upload"`

		FileUpload struct {
			ShortName string `yaml:"shortName" json:"short_name"`
			Enabled   bool   `yaml:"enabled" json:"enabled"`
		} `yaml:"fileupload" json:"file_upload"`

		ImageUpload struct {
			ShortName string `yaml:"shortName" json:"short_name"`
			Enabled   bool   `yaml:"enabled" json:"enabled"`
		} `yaml:"imageupload" json:"image_upload"`

		URLShortener struct {
			ShortName string `yaml:"shortName" json:"short_name"`
			Enabled   bool   `yaml:"enabled" json:"enabled"`
		} `yaml:"urlshortener" json:"url_shortener"`
		
		Polls struct {
			ShortName string `yaml:"shortName" json:"short_name"`
			Enabled   bool   `yaml:"enabled" json:"enabled"`
		} `yaml:"polls" json:"polls"`
	} `yaml:"modules" json:"modules"`

	UseFileSystem string `yaml:"useFileSystem"`
}

func (c *Config) LoadConfig() *ConfigStruct {
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		//config file does not exist, creating default config
		file, err := os.Create("config.yml")
		must(err)
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf(dEFAULT_CONFIG, time.Now().Format("Mon Jan 2 15:04:05 2006")))
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
