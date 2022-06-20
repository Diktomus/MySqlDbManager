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
	viper.AddConfigPath("config/app")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &MySqlDbConfig{
		Login:    viper.GetString("LOGIN"),
		Password: viper.GetString("PASSWORD"),
		Ip:       viper.GetString("IP"),
		Port:     viper.GetInt("PORT"),
		DbName:   viper.GetString("DB_NAME"),
		MaxConns: viper.GetInt("MAX_CONNS"),
	}

	return config, nil
}
