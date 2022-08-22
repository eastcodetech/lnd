package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/lightningnetwork/lnd/nmc_wallet"
)


type DataSchema struct {
	Protocol_Prfix string 
	Channel_Ref string
	// New_Channel_Ref string	
	Tx_Type	int64
	Addr_A_Balance int64
	Addr_B_Balance int64
	
}


func binaryToHexHelper(binary string) string {
        hex, err := strconv.ParseInt(binary, 2, 64)
        if err != nil {
                fmt.Print(err)
                return ""
        }
        return fmt.Sprintf("%x", hex)
}


func closeChannel() (proof, wallet string ){
	
	result_txID := nmc_wallet.NmcVerifyProofTest("04010100a7d0f257574b3777c38926f698ae29b4a03fe8dcef5f10324b8dfd0db7d49fb2a5d9035ced14676dbb91dc4ecf59fc364a2eb83462c3c8f3f14142fe6528f04972feeb62ffff7f200000000001000000010000000000000000000000000000000000000000000000000000000000000000ffffffff292813fae9d6e33e4e3cf660c89ab748396dfb7f932e807a22779794159fc83b05370100000000000000ffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000e4e7bf70dde022151c040d5e70b823ee10ce589136892db09b3620671a718c1a0000000000000000000000000f00000005bf991b3948e0e554a4d66c51a1d6613ce5495f69785371f7e38de6cf34c5c382110c9c02ef3368bb37f0af099c25a51960c0591c4e9dce32bf0ba86c1aa78c64bb88a058e583f54cbd6b90d77bdc7b04b88b69fa4fc5d5c7af0bc4058102b49f67968ce91b786240821acd5df153b1ab9e68b750fa7cbb98e84e59afedb780eaab44738b2d4ef15a6e4edb4caa7df88c41d9bf2bf8d58b8d4138b94e70543c7602ad00")
	// fmt.Println(result_txID)

	parsed_txID := result_txID.Result[0]
	
	rawProof := nmc_wallet.NmcGetRawTransactionTest("liam",parsed_txID)
	// fmt.Println(rawProof)

	decodeRaw := nmc_wallet.NmcDecodeRawTransactionTest(rawProof.Result)
	decoded_txid := decodeRaw.Result.Vout[0].ScriptPubKey.Asm
	// fmt.Println(decoded_txid)

	transformed_tostring_Id := string(decoded_txid)

	//trim off OP_RETURN 
	trimmed_Decoded_Id := strings.ReplaceAll(transformed_tostring_Id, "OP_RETURN ", "")
	// fmt.Println(trimmed_Decoded_Id)

	//decoding the hex, getting the plain "string" data struct
	decode_hex_string, err := hex.DecodeString(trimmed_Decoded_Id)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n", decode_hex_string)


	// converting string to binary 
	var bin_String string 
	for _, c := range decode_hex_string {
		bin_String = fmt.Sprintf("%s%b",bin_String, c)
		// fmt.Print(bin_String)
	}
		// fmt.Println(bin_String)
	data := DataSchema{ "zzz", "FFFFFFFF",65535,1099511627775, 1099511627775}

	/*

	BINARY CONVERSION HERE 
	
	*/
	res := ""
	for _, c := range data.Protocol_Prfix {
        	res = fmt.Sprintf("%s%.8b", res, c)	
	}
	length_Protocol := len(res)
		overP := length_Protocol % 4
		hex_Protocol := ""
		if overP > 0 {
			hex_Protocol = binaryToHexHelper(res[0:overP])
		}
		for i, j := overP, overP+4; j <= length_Protocol; i, j = i+4, j+4 {
		
			hex_Protocol += binaryToHexHelper(res[i:j])
		}


	// resp1 := ""
	// for _, c := range data.New_Channel_Ref {
        // 	resp1 = fmt.Sprintf("%s%.8b", resp1, c)	
	// }
	// length_New_Ch_Ref := len(resp1)
	// 	over1 := length_New_Ch_Ref % 4
	// 	hex_New_Ref := ""
	// 	if over1 > 0 {
	// 		hex_New_Ref = binaryToHexHelper(resp1[0:over1])
	// 	}
	// 	for i, j := over1, over1+4; j <= length_New_Ch_Ref; i, j = i+4, j+4 {
		
	// 		hex_New_Ref += binaryToHexHelper(resp1[i:j])
	// 	}
		
	resp := ""
	for _, c := range data.Channel_Ref {
        	resp = fmt.Sprintf("%s%.8b", resp, c)	
	}
	length_Ch_Ref := len(resp)
		over := length_Ch_Ref % 4
		hex_Ch_Ref := ""
		if over > 0 {
			hex_Ch_Ref = binaryToHexHelper(resp[0:over])
		}
		for i, j := over, over+4; j <= length_Ch_Ref; i, j = i+4, j+4 {
			// chan_ref_bin := resp[i:j]
			hex_Ch_Ref += binaryToHexHelper(resp[i:j])
		
		}
		
	
	tx_to_Binary := strconv.FormatInt(data.Tx_Type, 2)
		
	length_Tx_Type := len(tx_to_Binary)
		over3 := length_Tx_Type % 4
		hex_Tx_Type := ""
		if over3 > 0 {
			hex_Tx_Type = binaryToHexHelper(tx_to_Binary[0:over3])
		}
		for i, j := over3, over3+4; j <= length_Tx_Type; i, j = i+4, j+4 {
			
			hex_Tx_Type += binaryToHexHelper(resp[i:j])
		}


	addr_a_Binary := strconv.FormatInt(data.Addr_A_Balance, 2)
		// fmt.Println("Address A in binary ",addr_a_Binary)

	length_Addr_A := len(addr_a_Binary)
		over4 := length_Addr_A % 4
		hex_Addr_A := ""
		if over4 > 0 {
			hex_Addr_A = binaryToHexHelper(addr_a_Binary[0:over4])
		}
		for i, j := over4, over4+4; j <= length_Addr_A; i, j = i+4, j+4 {
			
			hex_Addr_A += binaryToHexHelper(addr_a_Binary[i:j])
		}


	addr_b_Binary := strconv.FormatInt(data.Addr_B_Balance, 2)
		// fmt.Println("Address B in binary ",addr_b_Binary)	

	length_Addr_B := len(addr_b_Binary)
		over5 := length_Addr_B % 4
		hex_Addr_B := ""
		if over5 > 0 {
			hex_Addr_B = binaryToHexHelper(addr_b_Binary[0:over5])
		}
		for i, j := over5, over5+4; j <= length_Addr_B; i, j = i+4, j+4 {
			
			hex_Addr_B += binaryToHexHelper(addr_b_Binary[i:j])
		}
		
	

	proof = hex_Protocol+ hex_Ch_Ref + hex_Tx_Type + hex_Addr_A + hex_Addr_B
	fmt.Println("Comnined hex values ",proof)
	fmt.Println(len(proof))

	
	return proof, wallet
	
} 


//this function takes in a comnbine hex values 
//breaks up the combined hexs values 
// returns individual strings 
func sliceHex(s string) (string,string,string)   {

	sfx := regexp.MustCompile(`[4-6-f]`).Split(s, -1)
	fmt.Println("From the slice hex function ", sfx)
	// transforming slice to string 
	str1 := strings.Join(sfx, " ")
	//removing white space
	fmt.Println(strings.TrimSpace(str1))


	sfx2 := regexp.MustCompile(`[7-a-f]`).Split(s, -1)
	fmt.Println("From the slice hex function ", sfx2)
	// transforming slice to string 
	str2 := strings.Join(sfx2, " ")
	//removing white space
	fmt.Println(strings.TrimSpace(str2))

	sfx3 := regexp.MustCompile(`[4-6-7-a]`).Split(s, -1)
	fmt.Println("From the slice hex function ", sfx3)
	// transforming slice to string 
	str3 := strings.Join(sfx3, " ")
	//removing white space
	fmt.Println(strings.TrimSpace(str3))

	fmt.Println(reflect.TypeOf(str3))

	return str1,str2,str3
	
}





func main (){
	// closeChannel()
	proof, _ := closeChannel()   
	sliceHex(proof)
	// Chunks(proof)
}


