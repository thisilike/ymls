package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
	"github.com/thisilike/ymls/scraper"
	"github.com/thisilike/ymls/utils"
	"gopkg.in/yaml.v3"
)

var log *logrus.Logger

func main() {
	log.Info("Reading scraper config")
	log.Debugf("test")
	config, err := ReadConfig("ymls.yml")
	if err != nil {
		log.Errorf("failed to read config: '%s'", err.Error())
		return
	}

	cnf := config.(map[string]interface{})
	configList := cnf["scrapers"].([]interface{})
	for _, scraperCnf := range configList {
		log.Error("create scraper")
		scraper, err := scraper.NewSraper(scraperCnf.(map[string]interface{}))
		if err != nil {
			log.Errorf("failed to create scraper: '%s'", err.Error())
			return
		}
		err = scraper.StartScraping()
		if err != nil {
			log.Errorf("failed scraping: '%s'", err.Error())
			return
		}
	}
}

func init() {
	log = config.Logger
	if config.Err != nil {
		fmt.Println(config.Err)
		panic(config.Err.Error())
	}
}

func ReadConfig(path string) (interface{}, error) {
	log.Debug("Reading Config")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("failed to read config file: '%s'", path)
		return nil, errors.New("failed to read config file")
	}

	var config interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Errorf("failed to parse config:'%s'", err.Error())
		return nil, errors.New("failed to parse config")
	}

	schemaText, err := ioutil.ReadFile("ymls.json")
	if err != nil {
		log.Errorf("failed to read schema: 'ymls.json'", err.Error())
		return nil, errors.New("failed to load config")
	}

	config, err = utils.ToStringKeys(config)
	if err != nil {
		log.Errorf("config invalid: '%s'", err.Error())
		return nil, errors.New("config invalid")
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(string(schemaText))); err != nil {
		log.Errorf("failed to add resource to schema compiler: '%s'", err.Error())
		return nil, errors.New("failed to add schema compiler resource")
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		log.Errorf("failed to compile schema: '%s'", err.Error())
		return nil, errors.New("failed to compile schema")
	}

	if err := schema.Validate(config); err != nil {
		log.Errorf("invalid config: '%s'", err.Error())
		return nil, errors.New("invalid config")
	}

	log.Debug("Config Valid!")
	return config, nil
}
