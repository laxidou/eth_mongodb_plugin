package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	config *viper.Viper
}

//func Init() (config *Config, _ error) {
func Init() (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// viper解析配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	//return &Config{v}, nil
	return v, nil
}




