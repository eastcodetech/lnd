package nmc_wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
This file is just a bunch of on chain rpc calls we can make to our
test blockchains
*/

// Hex string returns
type Hex struct {
	Result struct {
		hex string
	}
}

type Hex2 struct {
	Result []struct {
		hex string
	}
}

// struct for list unspent
type ListUnspent struct {
	Result []struct {
		Txid          string  `json:"txid"`
		Vout          int     `json:"vout"`
		Generated     bool    `json:"generated"`
		Address       string  `json:"address"`
		ScriptPubKey  string  `json:"scriptPubKey"`
		Amount        float64 `json:"amount"`
		Interest      float64 `json:"interest"`
		Confirmations int     `json:"confirmations"`
		Spendable     bool    `json:"spendable"`
	} `json:"result"`
	Error error  `json:"error"`
	ID    string `json:"id"`
}

// input struct
type Input struct {
	txid string
	vout int
}

//output struct
type Output struct {
	data string
}

// return struct for sign transaction
type SignRawTransaction struct {
	Result struct {
		Hex      string `json:"hex"`
		Complete bool   `json:"complete"`
		Errors   []struct {
			Txid      string `json:"txid"`
			Vout      int    `json:"vout"`
			ScriptSig string `json:"scriptSig"`
			Sequence  int64  `json:"sequence"`
			Error     string `json:"error"`
		} `json:"errors"`
	} `json:"result"`
	Error error  `json:"error"`
	ID    string `json:"id"`
}

// return struct for decoding transactions
type DecodeRawTransaction struct {
	Result struct {
		Txid     string `json:"txid"`
		Size     int    `json:"size"`
		Version  int    `json:"version"`
		Locktime int    `json:"locktime"`
		Vin      []struct {
			Txid      string `json:"txid"`
			Vout      int    `json:"vout"`
			ScriptSig struct {
				Asm string `json:"asm"`
				Hex string `json:"hex"`
			} `json:"scriptSig"`
			Sequence int64 `json:"sequence"`
		} `json:"vin"`
		Vout []struct {
			Value        float64 `json:"value"`
			ValueSat     int     `json:"valueSat"`
			N            int     `json:"n"`
			ScriptPubKey struct {
				Asm string `json:"asm"`
				Hex string `json:"hex"`
				// ReqSigs   int      `json:"reqSigs"`
				Addresses []string `json:"addresses"`
				Type      string   `json:"type"`
			} `json:"scriptPubKey"`
		} `json:"vout"`
		Vjoinsplit []interface{} `json:"vjoinsplit"`
	} `json:"result"`
	Error error  `json:"error"`
	ID    string `json:"id"`
}

type Transaction struct {
	Result struct {
		Amount            int      `json:"amount"`
		Fee               int      `json:"fee"`
		Confirmations     int      `json:"confirmations"`
		Trusted           bool     `json:"trusted"`
		Txid              string   `json:"txid"`
		Walletconflicts   []string `json:"walletconflicts"`
		Time              int      `json:"time"`
		Timereceived      int      `json:"timereceived"`
		Bip125Replaceable string   `json:"bip125-replaceable"`
		Details           []struct {
			InvolvesWatchonly bool   `json:"involvesWatchonly"`
			Address           string `json:"address"`
			Category          string `json:"category"`
			Amount            int    `json:"amount"`
			Vout              int    `json:"vout"`
			Fee               int    `json:"fee"`
			Abandoned         bool   `json:"abandoned"`
		} `json:"details"`
		Hex string `json:"hex"`
	} `json:"result"`
	Err error  `json:"error"`
	Id  string `json:"id"`
}

// Namecoin test proof generator
func NmcGetProofTest(txid string, wallet string) Hex {
	testRequest := fmt.Sprintf(`{"jsonrpc": "2.0", "id":"", "method": "gettxoutproof", "params": [["%s"]]}`, txid)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://10.10.10.120:8332/wallet/%s", wallet), strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Hex
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return Hex{}
}

// Namecoin test unspent transaction
func NmcUnlistTest(wallet string) ListUnspent {
	testRequest := `{"jsonrpc": "2.0", "id":"", "method": "listunspent", "params": []}`
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://10.10.10.120:8332/wallet/%s", wallet), strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j ListUnspent
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return ListUnspent{}
}

// Namecoin test proof verifier
func NmcVerifyProofTest(proof string) Hex2 {
	testRequest := fmt.Sprintf(`{"jsonrpc": "2.0", "id":"", "method": "verifytxoutproof", "params": ["%s"]}`, proof)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332", strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Hex2
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return Hex2{}
}

// Namecoin test transaction methods
func NmcCreateRawTransactionTest(in Input, out Output, wal string) Hex {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "createrawtransaction", "params": [[{"txid":"`, in.txid, `","vout":`, in.vout, `}], {"data": "`, out.data, `"}]}`)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://10.10.10.120:8332/wallet/%s", wal), strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Hex
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return Hex{}
}

func NmcSignRawTransactionTest(hex string, wal string) SignRawTransaction {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "signrawtransactionwithwallet", "params": ["`, hex, `"]}`)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://10.10.10.120:8332/wallet/%s", wal), strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j SignRawTransaction
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		fmt.Println(j.Result)
		return j
	}
	return SignRawTransaction{}
}

func NmcSendRawTransactionTest(wal string, hex string) Hex {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "sendrawtransaction", "params": ["`, hex, `", 0]}`)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://10.10.10.120:8332/wallet/%s", wal), strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Hex
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return Hex{}
}

func NmcGetRawTransactionTest(wal string, txid string) Hex {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "getrawtransaction", "params": ["`, txid, `"]}`)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://10.10.10.120:8332/wallet/%s", wal), strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Hex
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return Hex{}
}

func NmcDecodeRawTransactionTest(hex string) DecodeRawTransaction {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "decoderawtransaction", "params": ["`, hex, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332", strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j DecodeRawTransaction
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return DecodeRawTransaction{}
}

type Combine struct {
	Result string
}

type Finalize struct {
	Result struct {
		Psbt     string
		Hex      string
		Complete bool
	}
}

func BtcCombinePSBTTest(psbt1 string, psbt2 string) string {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "combinepsbt", "params": [["`, psbt1, `", "`, psbt2, `"]]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:18444", strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Combine
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j.Result
	}
	return ""
}

func BtcFinalizePsbtTest(psbt string) Finalize {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "finalizepsbt", "params": ["`, psbt, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:18444", strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Finalize
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return Finalize{}
}

func BtcSendRawTxTest(psbt string) string {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "sendrawtransaction", "params": ["`, psbt, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:18444/wallet/kyle", strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Combine
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j.Result
	}
	return ""
}

// Someone involved in the transaction has to do this i think
func BtcGetTransaction(txid string, wallet string) string {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "gettransaction", "params": ["`, txid, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:18444/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j Transaction
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j.Result.Hex
	}
	return ""
}

func BtcDecodeRawTransaction(hex string, wallet string) DecodeRawTransaction {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "decoderawtransaction", "params": ["`, hex, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:18444/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j DecodeRawTransaction
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return DecodeRawTransaction{}
}

type RawTransactionInput struct {
	Txid string
	Vout uint32
}

type RawTransactionResult struct {
	Result string `json:"result"`
	Err    struct {
	} `json:"error"`
	Id string `json:"id"`
}

func NmcCreateRawTransaction(input RawTransactionInput, address string, amount float64, message string, wallet string) (RawTransactionResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "createrawtransaction", "params": [[{"txid": "`, input.Txid, `", "vout": `, input.Vout, `}], 
	[{"`, address, `": `, amount, `}, {"data": "`, message, `"}]]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RawTransactionResult{}, err
	} else {
		defer res.Body.Close()
		var j RawTransactionResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

func NmcCreateEmptyRawTransaction(address string, message string, wallet string) (RawTransactionResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "createrawtransaction", "params": [[], 
	[{"`, address, `": `, 0.01, `}, {"data": "`, message, `"}]]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RawTransactionResult{}, err
	} else {
		defer res.Body.Close()
		var j RawTransactionResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

func NmcSignRawTransactionWithWallet(hex string, wallet string) (SignRawTransaction, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "signrawtransactionwithwallet", "params": ["`, hex, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SignRawTransaction{}, err
	} else {
		defer res.Body.Close()
		var j SignRawTransaction
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

type RawTxResult struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
	Id     string `json:"id"`
}

func NmcSendRawTransaction(hex string, wallet string) (RawTxResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "sendrawtransaction", "params": ["`, hex, `", 0]}`)
	// spew.Dump(testRequest)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RawTxResult{}, err
	} else {
		defer res.Body.Close()
		var j RawTxResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

type AddMultisigAddressResult struct {
	Result struct {
		Address      string
		RedeemScript string
		Descriptor   string
	}
	Error error
	Id    string
}

func NmcAddMultiSig(numSigs uint8, addresses []string, wallet string) (AddMultisigAddressResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "addmultisigaddress", "params": [`, numSigs, `, 
	["`, addresses[0], `", "`, addresses[1], `"]]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AddMultisigAddressResult{}, err
	} else {
		defer res.Body.Close()
		var j AddMultisigAddressResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

func NmcImportAddress(address string, wallet string) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "importaddress", "params": ["`, address, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("Could not import Multisig address for " + wallet)
	}
}

func NmcGenerateToAddress(numBlocks uint8, address string, wallet string) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "generatetoaddress", "params": [
		`, numBlocks, `, "`, address, `"
	]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("Couldn't Generate Blocks to new address for " + wallet)
	}
}

type TxResult struct {
	Result string
	Error  error
	Id     string
}

func NmcSendToAddress(address string, amount uint8, wallet string) (TxResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "sendtoaddress", "params": ["`, address, `", `, amount, `]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TxResult{}, err
	} else {
		defer res.Body.Close()
		var j TxResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

type FundedPsbtResult struct {
	Result struct {
		Psbt      string `json:"psbt"`
		Fee       uint64 `json:"fee"`
		Changepos uint64 `json:"changepos"`
	}
	Error error  `json:"error"`
	Id    string `json:"id"`
}

func NmcCreateFundedPsbt(addresses []string, amount []float32, inputs []string, wallet string) (FundedPsbtResult, error) {
	var addressStr string
	//! May have to look at the vout: it is currently hardcoded to 0
	for _, input := range inputs {
		addressStr = addressStr + fmt.Sprint(`{"txid": "`, input, `", "vout": 0, "sequence": 0},`)
	}
	var receiverStr string
	for index, address := range addresses {
		receiverStr = receiverStr + fmt.Sprint(`{"`, address, `": `, amount[index], `},`)
	}
	// removes trailing comma
	addressStr = addressStr[0 : len(addressStr)-1]
	receiverStr = receiverStr[0 : len(receiverStr)-1]
	var testRequest string
	if len(addresses) > 1 {
		testRequest = fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "walletcreatefundedpsbt", "params": [[`, addressStr, `], [`, receiverStr, `], 0, {"subtractFeeFromOutputs": [0, 1]}]}`)
	} else {
		testRequest = fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "walletcreatefundedpsbt", "params": [[`, addressStr, `], [`, receiverStr, `], 0, {"subtractFeeFromOutputs": [0]}]}`)
	}
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FundedPsbtResult{}, err
	} else {
		defer res.Body.Close()
		var j FundedPsbtResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

type NmcProcessPsbtResult struct {
	Result struct {
		Psbt     string
		Complete bool
	}
	Error error
	Id    string
}

func NmcProcessPsbt(psbt string, wallet string) (NmcProcessPsbtResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "walletprocesspsbt", "params": ["`, psbt, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return NmcProcessPsbtResult{}, err
	} else {
		defer res.Body.Close()
		var j NmcProcessPsbtResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

type NmcCombinePsbtResult struct {
	Result string
	Error  error
	Id     string
}

func NmcCombinePsbt(psbt []string, wallet string) (NmcCombinePsbtResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "combinepsbt", "params": [["`, psbt[0], `", "`, psbt[1], `"]]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return NmcCombinePsbtResult{}, err
	} else {
		defer res.Body.Close()
		var j NmcCombinePsbtResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

type NmcFinalizedPsbtResult struct {
	Result struct {
		Hex      string
		Complete bool
	}
	Error error
	Id    string
}

func NmcFinalizePsbt(psbt string, wallet string) (NmcFinalizedPsbtResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "finalizepsbt", "params": ["`, psbt, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return NmcFinalizedPsbtResult{}, err
	} else {
		defer res.Body.Close()
		var j NmcFinalizedPsbtResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

type getTransactionResult struct {
	Result struct {
		Amount            int      `json:"amount"`
		Fee               int      `json:"fee"`
		Confirmations     int      `json:"confirmations"`
		Blockhash         string   `json:"blockhash"`
		Blockheight       uint64   `json:"blockheight"`
		Blockindex        uint64   `json:"blockindex"`
		Blocktime         uint64   `json:"blocktime"`
		Txid              string   `json:"txid"`
		Walletconflicts   []string `json:"walletconflicts"`
		Time              int      `json:"time"`
		Timereceived      int      `json:"timereceived"`
		Bip125Replaceable string   `json:"bip125-replaceable"`
		Details           []struct {
			InvolvesWatchonly bool   `json:"involvesWatchonly"`
			Address           string `json:"address"`
			Category          string `json:"category"`
			Amount            int    `json:"amount"`
			Label             string `json:"label"`
			Vout              int    `json:"vout"`
			Fee               int    `json:"fee"`
			Abandoned         bool   `json:"abandoned"`
		} `json:"details"`
		Hex string `json:"hex"`
	} `json:"result"`
	Err error  `json:"error"`
	Id  string `json:"id"`
}

// Someone involved in the transaction has to do this i think
func NmcGetTransaction(txid string, wallet string) (getTransactionResult, error) {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "gettransaction", "params": ["`, txid, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return getTransactionResult{}, err
	} else {
		defer res.Body.Close()
		var j getTransactionResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j, nil
	}
}

func NmcDecodeRawTransaction(hex string) DecodeRawTransaction {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "decoderawtransaction", "params": ["`, hex, `"]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332", strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j DecodeRawTransaction
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return DecodeRawTransaction{}
}

type FundRawTxResult struct {
	Result struct {
		Hex       string
		Fee       uint64
		Changepos uint64
	}
	Error error
	Id    string
}

func NmcFundRawTransaction(hex string, wallet string) FundRawTxResult {
	testRequest := fmt.Sprint(`{"jsonrpc": "2.0", "id":"", "method": "fundrawtransaction", "params": ["`, hex, `", {"add_inputs": true, "subtractFeeFromOutputs": [0]}]}`)
	req, _ := http.NewRequest("POST", "http://10.10.10.120:8332/wallet/"+wallet, strings.NewReader(testRequest))
	req.SetBasicAuth("bitcoinrpc", "rpc")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		var j FundRawTxResult
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &j)
		return j
	}
	return FundRawTxResult{}
}
