package config

import "github.com/spf13/viper"

type MySqlDbConfig struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	MaxConns      int    `mapstructure:"MAX_CONNS"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func Load(configPath string) (dbConfig MySqlDbConfig, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		return MySqlDbConfig{}, err
	}

	err = viper.Unmarshal(&dbConfig)
	if err != nil {
		return MySqlDbConfig{}, err
	}

	return dbConfig, nil
}
