package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RPC string `yaml:"rpc"`

	UDT struct {
		Deps []struct {
			TxHash  string `yaml:"txHash"`
			Index   uint   `yaml:"index"`
			DepType string `yaml:"depType"`
		} `yaml:"deps"`
		Script struct {
			CodeHash string `yaml:"codeHash"`
			HashType string `yaml:"hashType"`
		} `yaml:"script"`
	} `yaml:"udt"`

	ACP struct {
		Deps []struct {
			TxHash  string `yaml:"txHash"`
			Index   int    `yaml:"index"`
			DepType string `yaml:"depType"`
		} `yaml:"deps"`
		Script struct {
			CodeHash string `yaml:"codeHash"`
			HashType string `yaml:"hashType"`
		} `yaml:"script"`
	} `yaml:"acp"`
}

func Init(path string) (*Config, error) {
	var c Config

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
