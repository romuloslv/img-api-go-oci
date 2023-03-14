package config

import (
	"bytes"
	_ "embed"
	"log"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var defaultConfiguration []byte

type Export struct {
	ImageId       string
	CompartmentId string
	BucketName    string
	Namespace     string
	ObjectName    string
}

type Config struct {
	Export *Export
}

func Read() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("OCI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetConfigType("yml")

	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfiguration)); err != nil {
		log.Fatalln(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}
	return &config, nil
}
