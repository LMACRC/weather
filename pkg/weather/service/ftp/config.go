package ftp

type Config struct {
	Address  string
	Username string
	Password string
}

func NewConfig() Config {
	return Config{}
}
