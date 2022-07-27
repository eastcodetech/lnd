package main

import (
	"fmt"

	"github.com/eastcodetech/lnd/nmc_client/user"
)

func main() {
	client := user.GetNmcClient()
	fmt.Print(client)
}
