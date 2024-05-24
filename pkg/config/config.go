package config

type Config struct {
	JWT struct {
		SECRETKEY  string
		EXPIRES_AT int
	}
	MODEM struct {
		PORT string
	}
	GRPC struct {
		PORT string
	}
}
