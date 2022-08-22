package nmc_client

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

type Content struct {
	method string
	params []string
}

func GetBtcClient(wallet string) *rpcclient.Client {
	config := rpcclient.ConnConfig{
		Host:         "10.10.10.120:18444/wallet/" + wallet,
		User:         "bitcoinrpc",
		Pass:         "rpc",
		DisableTLS:   true,
		HTTPPostMode: true,
		Params:       chaincfg.RegressionNetParams.Name,
	}

	config_ptr := &config

	client, err := rpcclient.New(config_ptr, nil)
	if err != nil {
		panic(err)
	}

	return client
}

func GetNmcClient(wallet string) *rpcclient.Client {
	config := rpcclient.ConnConfig{
		Host:         "10.10.10.120:8332/wallet/" + wallet,
		User:         "bitcoinrpc",
		Pass:         "rpc",
		DisableTLS:   true,
		HTTPPostMode: true,
		Params:       chaincfg.RegressionNetParams.Name,
	}

	config_ptr := &config

	client, err := rpcclient.New(config_ptr, nil)
	if err != nil {
		panic(err)
	}

	return client
}

func GetBtcClient2() *rpcclient.Client {
	config := rpcclient.ConnConfig{
		Host:         "10.10.10.120:18444",
		User:         "bitcoinrpc",
		Pass:         "rpc",
		DisableTLS:   true,
		HTTPPostMode: true,
		Params:       chaincfg.RegressionNetParams.Name,
	}

	config_ptr := &config

	client, err := rpcclient.New(config_ptr, nil)
	if err != nil {
		panic(err)
	}

	return client
}
