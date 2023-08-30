package services

import (
	"fmt"
	"sync"

	"scan-eth/pkg/token"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20Token struct {
	Symbol   string
	Decimals uint8
}

type Erc20TokenCache struct {
	sync.RWMutex
	m map[string]*ERC20Token
}

func (c *Erc20TokenCache) get(address string) (*ERC20Token, bool) {
	c.RLock()
	t, ok := c.m[address]
	c.RUnlock()

	return t, ok
}

func (c *Erc20TokenCache) set(address string, t *ERC20Token) {
	c.Lock()
	defer c.Unlock()

	c.m[address] = t
}

var erc20TokenCache = &Erc20TokenCache{
	m: make(map[string]*ERC20Token),
}

func getToken(address string) (*ERC20Token, bool) {
	return erc20TokenCache.get(address)
}

func setToken(address string, t *ERC20Token) {
	erc20TokenCache.set(address, t)
}

func isERC20Token(address common.Address, client *ethclient.Client) (*ERC20Token, error) {
	instance, err := token.NewToken(address, client)
	if err != nil {
		return nil, fmt.Errorf("token NewToken error: %s", err)
	}

	symbol, err := instance.Symbol(nil)
	if err != nil {
		return nil, fmt.Errorf("token symbol error: %s", err)
	}

	decimals, err := instance.Decimals(nil)
	if err != nil {
		return nil, fmt.Errorf("token symbol error: %s", err)
	}

	return &ERC20Token{
		Symbol:   symbol,
		Decimals: decimals,
	}, nil
}
