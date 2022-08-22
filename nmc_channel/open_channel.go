package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/lightningnetwork/lnd/nmc_wallet"
)

// This is a sample script of all the calls that would have to be made to
// mimic what the lightning OpenChannel function does.
// This will print the txid for the broadcasted nmc transaction.
//! Sometimes the opening channel errors out but the money is still being sent to the multisig
func main() {
	// amount to be entered into the multisig wallet
	amount, _ := btcutil.NewAmount(1.0)

	btcAddress1 := "bcrt1qekn2wfen3qe90vxe9semv2mqm0wt3880ha4pne"
	btcAddress2 := "bcrt1qavyf856ehh9gdehn2xne030q38qvzdsugy24aq"
	pubKey1 := "03f5bdb7460ff0103acb865495730312b68b0866a96cbcf0d48f14f1c98706a71d"
	pubKey2 := "03e4c8af01a9ddfbbb8de335283a6282b6b96831b75c53f4164698f610622f976d"

	txid, err := OpenChannel(btcAddress1, "brandon", pubKey1, btcAddress2, "liam", pubKey2, amount)
	if len(txid) < 2 || err != nil {
		panic(fmt.Sprint("Did not receive the txid", err))
	}

	nmcAddress1 := "n3yYYSfH932jt8kWGVj8hrtWVkaXrujvY1"
	nmcWalletName := "brandon"

	// These are just random prefix
	binaryPrefix := "011001010111010001100011"
	binaryTxType := "0000"

	// Will have to adjust these numbers to be the correct number of bits
	binaryBalance1 := intToBinary(int64(amount))

	// These function just return random binary values
	channelRef := getChannelRef()
	newChannelRef := getNewChannelRef()

	// Will have to decide what scheme we should use
	binary := fmt.Sprintf("%s%s%s%s%s", binaryPrefix, binaryTxType, binaryBalance1,
		channelRef, newChannelRef)

	hex := binaryToHex(binary)

	txid, err = nmcBroadcastMessage(hex, nmcAddress1, nmcWalletName)
	if err != nil {
		fmt.Println("Error broadcasting message")
		panic(err)
	}

	fmt.Print(txid)
}

// This will often fail for two reasons
// 1. Insufficient funds
// 2. A Psbt that shouldn't be complete is complete
func OpenChannel(address1 string, name1 string, pubkey1 string, address2 string, name2 string, pubkey2 string, amount btcutil.Amount) (string, error) {

	config1 := rpcclient.ConnConfig{
		Host:         "10.10.10.120:18444/wallet/" + name1,
		Endpoint:     "ws",
		User:         "bitcoinrpc",
		Pass:         "rpc",
		Params:       chaincfg.RegressionNetParams.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	config2 := rpcclient.ConnConfig{
		Host:         "10.10.10.120:18444/wallet/" + name2,
		Endpoint:     "ws",
		User:         "bitcoinrpc",
		Pass:         "rpc",
		Params:       chaincfg.RegressionNetParams.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	user1, err := rpcclient.New(&config1, nil)
	if err != nil {
		panic(fmt.Sprint("Error creating user1 client: ", err))
	}

	user2, err := rpcclient.New(&config2, nil)
	if err != nil {
		panic(fmt.Sprint("Error creating user2 client: ", err))
	}

	user1Address, err := btcutil.DecodeAddress(pubkey1, &chaincfg.RegressionNetParams)
	user2Address, err := btcutil.DecodeAddress(pubkey2, &chaincfg.RegressionNetParams)
	addresses := []btcutil.Address{user2Address, user1Address}

	// I'm not sure what the last param is for this function
	// leaving it blank seems to work
	multisigAddress, err := user1.AddMultisigAddress(2, addresses, "")

	err = user1.ImportAddress(multisigAddress.String())
	if err != nil {
		panic(fmt.Sprint("Error importing multisig to user1: ", err))
	}
	err = user2.ImportAddress(multisigAddress.String())
	if err != nil {
		panic(fmt.Sprint("Error importing multisig to user2: ", err))
	}

	maxTries := int64(2)
	_, err = user1.GenerateToAddress(1, multisigAddress, &maxTries)
	if err != nil {
		panic(fmt.Sprint("Error generating to multisig address: ", err))
	}

	txid, err := user1.SendToAddress(multisigAddress, amount)

	_, err = user1.GenerateToAddress(2, multisigAddress, &maxTries)
	if err != nil {
		panic(fmt.Sprint("Error generating to multisig address: ", err))
	}

	psbtInput := []btcjson.PsbtInput{btcjson.PsbtInput{
		Txid:     txid.String(),
		Vout:     0,
		Sequence: 0,
	}}

	psbtOutput := []btcjson.PsbtOutput{btcjson.NewPsbtOutput(address1, amount)}
	locktime := uint32(0)
	options := btcjson.WalletCreateFundedPsbtOpts{}
	bip32Derivs := true

	fundedPsbtResult, err := user1.WalletCreateFundedPsbt(
		psbtInput, psbtOutput, &locktime, &options, &bip32Derivs)
	if err != nil {
		panic(fmt.Sprint("Error funding Psbt: ", err))
	}
	sign := true
	user1ProcessPsbtResult, err := user1.WalletProcessPsbt(
		fundedPsbtResult.Psbt, &sign, rpcclient.SigHashAll, &bip32Derivs)
	if err != nil {
		panic(fmt.Sprint("Error processing Psbt for user1: ", err))
	}
	// ! for some reason this keeps coming back as complete
	if user1ProcessPsbtResult.Complete {
		panic(fmt.Sprint("Half signed psbt was complete for user1"))
	}

	user2ProcessPsbtResult, err := user2.WalletProcessPsbt(
		fundedPsbtResult.Psbt, &sign, rpcclient.SigHashAll, &bip32Derivs)
	if err != nil {
		panic(fmt.Sprint("Error processing Psbt for user2: ", err))
	}
	if user2ProcessPsbtResult.Complete {
		panic(fmt.Sprint("Half signed psbt was complete for user2"))
	}

	combinedPsbtStr := nmc_wallet.BtcCombinePSBTTest(
		user1ProcessPsbtResult.Psbt, user2ProcessPsbtResult.Psbt)

	finalizedPsbt := nmc_wallet.BtcFinalizePsbtTest(combinedPsbtStr)

	txid1 := nmc_wallet.BtcSendRawTxTest(finalizedPsbt.Result.Hex)

	return txid1, nil

}

// Creates the NMC transaction with the message
func nmcBroadcastMessage(hex string, nmcAddress1 string, nmcWalletName1 string) (string, error) {

	txid, err := nmc_wallet.NmcCreateEmptyRawTransaction(nmcAddress1, hex, nmcWalletName1)
	if err != nil {
		return "", err
	}

	txid1 := nmc_wallet.NmcFundRawTransaction(txid.Result, nmcWalletName1)

	txid2, err1 := nmc_wallet.NmcSignRawTransactionWithWallet(txid1.Result.Hex, nmcWalletName1)
	if err1 != nil {
		return "", err
	}

	txid3, err2 := nmc_wallet.NmcSendRawTransaction(txid2.Result.Hex, nmcWalletName1)

	return txid3.Result, err2
}

func getAddress(client1 rpcclient.Client, address1 string) (btcutil.Address, error) {

	client1AddInfo, err := client1.GetAddressInfo(address1)
	if err != nil {
		fmt.Print("Error gettting client1's info")
		return nil, err
	}

	client1PubKey := client1AddInfo.PubKey

	client1Address, err := btcutil.DecodeAddress(*client1PubKey, &chaincfg.RegressionNetParams)
	if err != nil {
		fmt.Println("Couldn't decode address")
		return nil, err
	}

	return client1Address, nil
}

func getUtxo(client *rpcclient.Client) (btcjson.ListUnspentResult, error) {

	unspentResponse, err := client.ListUnspent()
	if err != nil {
		return btcjson.ListUnspentResult{}, err
	}

	return unspentResponse[0], nil
}

// Converts a long binary string to a hex string
func binaryToHex(binary string) string {
	length := len(binary)
	over := length % 4
	hex := ""
	if over > 0 {
		hex = binaryToHexHelper(binary[0:over])
	}

	for i, j := over, over+4; j <= length; i, j = i+4, j+4 {
		hex += binaryToHexHelper(binary[i:j])
	}
	if len(hex)%2 != 0 {
		hex = "0" + hex
	}
	return hex
}

// Can only convert short binary strings to hex
func binaryToHexHelper(binary string) string {
	hex, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	return fmt.Sprintf("%x", hex)
}

func intToBinary(satoshies int64) string {
	int := big.NewInt(satoshies)
	binary := fmt.Sprintf("%b", int) // 11111011110
	return binary
}

// Returns array of balances, user A's balance is at index 0
// and user B's balance is at index 1
// could use mapping and iterator for order
func getAddresses(decodedTx nmc_wallet.DecodeRawTransaction, address1 string, address2 string) []uint64 {
	var addressesTemp []string
	var balancesTemp []uint64
	for _, element := range decodedTx.Result.Vout {
		// Needs to match a users address because it might be the multisig wallet address
		if element.ScriptPubKey.Addresses[0] == address1 || element.ScriptPubKey.Addresses[0] == address2 {
			addressesTemp = append(addressesTemp, element.ScriptPubKey.Addresses[0])
			balancesTemp = append(balancesTemp, uint64(element.Value*100000000))
		}
	}
	var balances [2]uint64
	if addressesTemp[0] > addressesTemp[1] {
		balances[0] = balancesTemp[1]
		balances[1] = balancesTemp[0]

		return balances[:]
	}
	return balancesTemp[:]
}

func getChannelRef() string {
	return "1010101111111110101"
}

func getNewChannelRef() string {
	return "0000000000000000000"
}
