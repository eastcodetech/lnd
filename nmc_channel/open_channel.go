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

// Opens a channel and logs the open channel into Namecoin
func main() {
	fundingAmount := 1.0

	btcAddress1 := "bcrt1qekn2wfen3qe90vxe9semv2mqm0wt3880ha4pne"
	btcAddress2 := "bcrt1qavyf856ehh9gdehn2xne030q38qvzdsugy24aq"
	pubKey1 := "03f5bdb7460ff0103acb865495730312b68b0866a96cbcf0d48f14f1c98706a71d"
	pubKey2 := "03e4c8af01a9ddfbbb8de335283a6282b6b96831b75c53f4164698f610622f976d"
	btcWallet1 := "brandon"
	btcWallet2 := "liam"

	txid, err := OpenChannel(btcAddress1, btcWallet1, pubKey1, btcAddress2, btcWallet2, pubKey2, fundingAmount)
	if len(txid) < 2 || err != nil {
		panic(fmt.Sprint("Did not receive the txid", err))
	}

	nmcAddress1 := "n3yYYSfH932jt8kWGVj8hrtWVkaXrujvY1"
	nmcWalletName := "brandon"

	// These are just random prefixes for now
	binaryPrefix := "011001010111010001100011"
	binaryTxType := "0000"

	// Will have to adjust these numbers to be the correct number of bits
	binaryBalance1 := intToBinary(int64(fundingAmount))

	// These function just return random binary values for now
	//* We will have to figure how to get the values from Lnd
	channelRef := getChannelRef()
	newChannelRef := getNewChannelRef()

	// We can skip this if we are only using hex
	binary := fmt.Sprintf("%s%s%s%s%s", binaryPrefix, binaryTxType, binaryBalance1,
		channelRef, newChannelRef)

	hex := binaryToHex(binary)
	fmt.Println("The message is " + hex)

	txid, err = nmcBroadcastMessage(hex, nmcAddress1, nmcWalletName)
	if err != nil {
		panic(fmt.Sprint("Failed to broadcast message ", err))
	}

	fmt.Print("The transaction id in NMC is " + txid)
}

func OpenChannel(address1 string, btcWallet1 string, pubkey1 string, address2 string, btcWallet2 string, pubkey2 string, amount float64) (string, error) {

	addresses := []string{pubkey1, pubkey2}

	multisigAddress, err := nmc_wallet.BtcAddMultiSig(2, addresses, btcWallet1)
	if err != nil {
		panic("Could not make multisig")
	}

	nmc_wallet.BtcImportAddress(multisigAddress.Result.Address, btcWallet1)
	nmc_wallet.BtcImportAddress(multisigAddress.Result.Address, btcWallet2)

	nmc_wallet.BtcGenerateToAddress(1, multisigAddress.Result.Address, btcWallet1)

	_, err = nmc_wallet.BtcSendToAddress(multisigAddress.Result.Address, amount, btcWallet1)
	if err != nil {
		panic(fmt.Sprint("Could not send to Btc ", err))
	}

	nmc_wallet.BtcGenerateToAddress(1, multisigAddress.Result.Address, btcWallet1)

	psbtResult, err := nmc_wallet.BtcCreateFundedPsbt(address1, amount, btcWallet1)

	user1Psbt, err := nmc_wallet.BtcProcessPsbt(psbtResult.Result.Psbt, btcWallet1)
	if user1Psbt.Result.Complete {
		panic("User1's Psbt was complete")
	}
	user2Psbt, err := nmc_wallet.BtcProcessPsbt(psbtResult.Result.Psbt, btcWallet2)
	if user2Psbt.Result.Complete {
		panic("User1's Psbt was complete")
	}

	completePsbt := nmc_wallet.BtcCombinePSBTTest(user1Psbt.Result.Psbt, user2Psbt.Result.Psbt)

	finalizedPsbt := nmc_wallet.BtcFinalizePsbtTest(completePsbt, btcWallet1)

	txid := nmc_wallet.BtcSendRawTxTest(finalizedPsbt.Result.Hex, btcWallet1)
	return txid, nil

}

// Creates the NMC transaction with the message
func nmcBroadcastMessage(hex string, nmcAddress1 string, nmcWalletName1 string) (string, error) {

	txid, err := nmc_wallet.NmcCreateEmptyRawTransaction(nmcAddress1, hex, nmcWalletName1)
	if err != nil {
		panic(fmt.Sprint("Failed to create empty raw transaction on NMC ", err))
	}
	txid1, err := nmc_wallet.NmcFundRawTransaction(txid.Result, nmcWalletName1)
	if err != nil {
		panic(fmt.Sprint("Failed to fund raw transaction on NMC ", err))
	}

	txid2, err := nmc_wallet.NmcSignRawTransactionWithWallet(txid1.Result.Hex, nmcWalletName1)
	if err != nil {
		panic(fmt.Sprint("Failed to sign raw transaction on NMC ", err))
	}

	txid3, err := nmc_wallet.NmcSendRawTransaction(txid2.Result.Hex, nmcWalletName1)
	if err != nil {
		panic(fmt.Sprint("Failed to send raw transaction on NMC ", err))
	}
	return txid3.Result, err
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
