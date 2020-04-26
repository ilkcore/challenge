package env

import (
	"github.com/namsral/flag"
)

type Config struct {
	Buffersize int
	Timeout    int
	ImageUrl   string
}

func Get() *Config {
	cfg := Config{}
	buffersize := flag.Int("buffersize", 0, "reading buffer size in bytes")
	timeout := flag.Int("timeout", 0, "http client timeout in seconds")
	imageUrl := flag.String("imageurl", "", "overwrite default image url")
	flag.Parse()

	cfg.Buffersize = *buffersize
	cfg.Timeout = *timeout
	cfg.ImageUrl = *imageUrl

	return &cfg
}
