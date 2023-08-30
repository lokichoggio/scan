package config

import (
	"scan-eth/pkg/log"
	"scan-eth/pkg/mysql"
)

type Config struct {
	EthAddr   string `yaml:"ethAddr"`
	EthWsAddr string `yaml:"ethWsAddr"`

	ERC20ABIFilePath string `yaml:"erc20ABIFilePath"`

	StartBlockNumber        int64 `yaml:"startBlockNumber"`
	ScanIntervalBlock       int64 `yaml:"scanIntervalBlock"`
	ScanIntervalTimeSeconds int64 `yaml:"scanIntervalTimeSeconds"`

	SpecificAddressList []string `yaml:"specificAddressList"`

	Db  mysql.Config `yaml:"db"`
	Log log.Config   `yaml:"log"`
}

func (c *Config) ContainAddress(addressList []string) bool {
	for _, specificAddress := range c.SpecificAddressList {
		for _, customAddress := range addressList {
			if specificAddress == customAddress {
				return true
			}
		}
	}

	return false
}
