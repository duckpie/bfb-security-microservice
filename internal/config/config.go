package config

type Config struct {
	Services      ServicesConfigs     `toml:"services"`
	Microservices MicroservicesConfig `toml:"microservices"`
}

type ServicesConfigs struct {
	Server ServerConfig `toml:"server"`
	Redis  RedisConfig  `toml:"redis"`
}

type MicroservicesConfig struct {
	UserMs UserMsConfig `toml:"user"`
}

type UserMsConfig struct {
	Host string `toml:"HOST"`
	Port int64  `toml:"PORT"`
}

type ServerConfig struct {
	Port          int64  `toml:"PORT"`
	AccessSecret  string `toml:"ACCESS_SECRET"`
	RefreshSecret string `toml:"REFRESH_SECRET"`
}

type RedisConfig struct {
	Host string `toml:"HOST"`
	Port int64  `toml:"PORT"`
	DB   int8   `toml:"DB"`
}

func NewConfig() *Config {
	return &Config{}
}
