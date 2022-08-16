package config

import (
	"fmt"
	"os"

	"eat-and-go/model"
	CONST "eat-and-go/pkg/constvar"

	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var Conf *model.SysConfig

func GetConfig() *model.SysConfig {
	return Conf
}

func SetConfig() error {
	if os.Getenv("APP_ENV") == CONST.ConfigOnline {
		fmt.Println("Fetch config from apollo online.")
		if err := FetchApolloConfig(
			os.Getenv("APOLLO_ADDR"),
			CONST.ApolloAppID,
			os.Getenv("APOLLO_NAMESPACE")+"."+CONST.ApolloConfigType,
			CONST.ApolloConfigType); err != nil {
			return err
		}
	} else {
		fmt.Println("fetch config from local config file.")
		if err := FetchLocalConfig(CONST.ConfigLocalPath); err != nil {
			return err
		}
	}
	return nil
}

func FetchLocalConfig(path string) error {
	var config model.SysConfig
	dir, _ := os.Getwd()
	f, err := os.Open(dir + path)
	if err != nil {
		return err
	}

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return err
	}
	Conf = &config
	fmt.Println(Conf)
	return nil
}

func FetchApolloConfig(addr string, appID string, namespace string, configType string) error {
	var config model.SysConfig
	if err := initApolloConfig(&config, addr, appID, namespace, configType); err != nil {
		return err
	}
	Conf = &config
	fmt.Println(Conf)
	return nil
}

// init apollo config
func initApolloConfig(cfg interface{}, addr string, appID string, namespace string, configType string) error {
	remote.SetAppID(appID)
	v := viper.New()
	v.SetConfigType(configType)

	if err := v.AddRemoteProvider("apollo", addr, namespace); err != nil {
		fmt.Println(err)
		return err
	}
	// error handle...
	if err := v.ReadRemoteConfig(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		// panic(err)
		fmt.Println(err)
		return err
	}

	return nil
}
