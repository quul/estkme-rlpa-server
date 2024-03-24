package config

import (
	"errors"
)

type Config struct {
	LpacVersion         string
	DataDir             string
	DontDownload        bool
	RLPAListenIP        string
	RLPAListenPortRange string

	APIServerMode          string // singleUser or multiUser
	APIServerListenAddress string
	DBFilePath             string
	Secret                 string

	DefaultUserName string
	DefaultPassword string
	EnableRegister  bool
}

var C = &Config{}

var (
	ErrLpacVersionEmpty = errors.New("lpac version is empty")
)

func (c *Config) IsValid() error {
	// TODO
	//if _, err := net.ResolveTCPAddr("tcp", c.ListenAddress); err != nil {
	//	return err
	//}
	//if c.LpacVersion == "" {
	//	return ErrLpacVersionEmpty
	//}
	//return nil
}
