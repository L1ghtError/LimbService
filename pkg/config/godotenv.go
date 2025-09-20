package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrNotFound = errors.New("could not load config")
)

type goDotEnv struct {
	filenames []string
	params    map[string]string
}

func NewGoDotEnv() Config {
	return &goDotEnv{}
}

func (g *goDotEnv) Init(filenames ...string) error {
	if err := godotenv.Load(filenames...); err != nil {
		return ErrNotFound
	}
	return nil
	// TODO: Remove code above, uncomment code below
	// var err error
	// if g.params, err = godotenv.Read(filenames...); err != nil {
	// 	return ErrNotFound
	// }
	// return nil
}

func (g goDotEnv) Get(key string) string {
	return os.Getenv(key)
	// TODO: Remove code above, uncomment code below
	// return g.params[key]
}
