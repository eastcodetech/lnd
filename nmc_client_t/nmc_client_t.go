package main

import (
	"github.com/eastcodetech/lnd/nmc_client"
)

func main() {
	nmc_client = nmc_client.user.GetNmcClient()
	fmt.Print(nmc_client)
}
