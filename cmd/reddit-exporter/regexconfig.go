package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/travisgroth/reddit-exporter/internal/handlers"
	"github.com/travisgroth/reddit-exporter/internal/scanner"
)

var (
	RegexConfigs map[string]map[string]string
)

func addRegexHandlers(s *scanner.Scanner, file string) error {

	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Could not load config file %s: %s", file, err)
	}

	err = v.Unmarshal(&RegexConfigs)
	if err != nil {
		return fmt.Errorf("Could not parse config file %s: %s", file, err)
	}

	log.Infof("Parsing regex config %s", file)

	for matchGroupName, matchGroupConfig := range RegexConfigs {
		regexHandler, _ := handlers.NewRegex(matchGroupName)
		for matchName, matchConfig := range matchGroupConfig {
			err := regexHandler.AddMatch(matchName, matchConfig)

			if err != nil {
				return fmt.Errorf("Could not load regex config %s: %s", file, err)
			}
		}
		s.AddCommentHandler(regexHandler)
		s.AddPostHandler(regexHandler)
	}

	return nil

}
