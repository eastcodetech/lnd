package nmc_network

import (
	"github.com/btcsuite/btcd/rpcclient"
)

type Content struct {
	method string
	params []string
}

func GetNmcClient() *rpcclient.Client {
	config := rpcclient.ConnConfig{
		Host:         "10.10.10.120:8332/wallet/erik",
		User:         "bitcoinrpc",
		Pass:         "rpc",
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	config_ptr := &config

	client, err := rpcclient.New(config_ptr, nil)
	if err != nil {
		panic(err)
	}

	return client
}
