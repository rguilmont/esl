package configuration

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/viper"
)

// Context are all settings needed to access an elasticsearch cluster
type Context struct {
	URL                  *string
	Filter               *string
	Refresh              *int64
	VerifySSLCertificate *bool

	Username *string
	Password *string
	Index    *string
}

// Configuration of the application, basicaly a list of context.
type Configuration map[string]Context

const configName = "config"

// default path to configuration
var defaultPath = []string{"/etc/esl/", "$HOME/.esl/", "."}

var (
	defaultFilter    = "*"
	defaultRefresh   = 200 * time.Millisecond
	defaultSSLVerify = true
	defaultIndex     = "logstash*"
)

// LoadConf loads configuration, containing all contexts
func LoadConf(overridedConf io.Reader) (*Configuration, error) {
	if overridedConf != nil {
		err := viper.ReadConfig(overridedConf)
		if err != nil {
			return nil, err
		}

	} else {
		viper.SetConfigName(configName) // name of config file (without extension)
		for _, p := range defaultPath {
			viper.AddConfigPath(p)
		}
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}
	var runtimeConf Configuration

	err := viper.Unmarshal(&runtimeConf)
	if err != nil {
		return nil, err
	}
	return &runtimeConf, nil
}

// Context return context, if exists. If not, return nil and an error.
func (conf Configuration) Context(c string) (*Context, error) {
	curentContext, ok := conf[c]

	if !ok {
		return nil, fmt.Errorf("Unknown context %v", c)
	}

	// Set default values if nil.
	if curentContext.Refresh == nil {
		curentContext.Refresh = new(int64)
		*curentContext.Refresh = int64(defaultRefresh)
	} else {
		*curentContext.Refresh = *curentContext.Refresh * int64(time.Millisecond)
	}
	if curentContext.VerifySSLCertificate == nil {
		curentContext.VerifySSLCertificate = &defaultSSLVerify
	}
	if curentContext.Filter == nil {
		curentContext.Filter = &defaultFilter
	}
	return &curentContext, nil
}
