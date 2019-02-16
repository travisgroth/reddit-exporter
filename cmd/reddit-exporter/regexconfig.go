package main

import (
	"fmt"

	"github.com/turnage/graw/botfaces"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/travisgroth/reddit-exporter/internal/handlers"
)

var (
	RegexConfigs map[string]map[string]string
)

type contentScanner interface {
	AddCommentHandler(botfaces.CommentHandler)
	AddPostHandler(botfaces.PostHandler)
}

func addRegexHandlers(s contentScanner, file string) error {

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
		log.Debugf("Adding matchgroup: %s", matchGroupName)
		regexHandler, _ := handlers.NewRegex(matchGroupName)
		for matchName, matchConfig := range matchGroupConfig {
			log.Debugf("Adding match: %s: '%s'", matchName, matchConfig)
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
