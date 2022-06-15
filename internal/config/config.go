package config

import "github.com/spf13/viper"

type MySqlDbConfig struct {
	Login    string
	Password string
	Ip       string
	Port     int
	DbName   string
	MaxConns int
}

func Init() (*MySqlDbConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &MySqlDbConfig{
		Login:    viper.GetString("login"),
		Password: viper.GetString("password"),
		Ip:       viper.GetString("ip"),
		Port:     viper.GetInt("port"),
		DbName:   viper.GetString("dbName"),
		MaxConns: viper.GetInt("maxConns"),
	}

	return config, nil
}
