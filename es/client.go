package es

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit2/config"
)

var Client *elastic.Client

const (
	retryCount   = 15
	retrySeconds = 2
)

func init() {
	var err error
	logger := jsonlog.DefaultLogger
	options := getElasticSearchConfig()
	for r := retryCount; r > 0; r-- {
		Client, err = elastic.NewClient(options...)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to connect to ElasticSearch database. Retrying in %s seconds.", retrySeconds), err.Error())
			time.Sleep(retrySeconds * time.Second)
		} else {
			logger.Info("Successfully connected to ElasticSearch database.", nil)
			return
		}
	}
	logger.Error("Failed to connect to ElasticSearch database. Not retrying.", nil)
}

func getElasticSearchConfig() []elastic.ClientOptionFunc {
	return []elastic.ClientOptionFunc{
		getElasticSearchUrlConfig(),
		getElasticSearchAuthConfig(),
	}
}

func getElasticSearchUrlConfig() elastic.ClientOptionFunc {
	return elastic.SetURL(config.EsAddress)
}

func getElasticSearchAuthConfig() elastic.ClientOptionFunc {
	authType, authValue := getElasticSearchAuthTypeAndValue()
	logger := jsonlog.DefaultLogger
	switch authType {
	case "basic":
		logger.Debug("Configuring ElasticSearch client with basic auth.", nil)
		return getElasticSearchBasicAuth(authValue)
	case "none":
		logger.Debug("Configuring ElasticSearch client with null auth.", nil)
	default:
		logger.Error("Could not configure ElasticSearch client auth: bad auth format.", nil)
		os.Exit(1)
	}
	return configNoop
}

func getElasticSearchBasicAuth(auth string) elastic.ClientOptionFunc {
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) == 2 {
		return elastic.SetBasicAuth(parts[0], parts[1])
	} else {
		logger := jsonlog.DefaultLogger
		logger.Error("Could not configure ElasticSearch client basic auth: missing username or password.", nil)
		os.Exit(1)
		return nil
	}
}

func getElasticSearchAuthTypeAndValue() (string, string) {
	parts := strings.SplitN(config.EsAuthentication, ":", 2)
	if len(parts) == 0 {
		return "none", ""
	} else if len(parts) == 1 {
		return parts[0], ""
	} else {
		return parts[0], parts[1]
	}
}

func configNoop(_ *elastic.Client) error { return nil }
