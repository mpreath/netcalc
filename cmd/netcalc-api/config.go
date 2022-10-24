package main

type Config struct {
	HttpPort int
}

func LoadConfig() *Config {
	return &Config{
		HttpPort: 3000,
	}
}
