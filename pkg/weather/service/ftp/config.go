package ftp

type Config struct {
	Address  string
	Username string
	Password string
	Retries  int
}

func NewConfig() Config {
	return Config{
		Retries: 5,
	}
}
