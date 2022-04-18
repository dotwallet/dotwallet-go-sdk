package main

import (
	"log"
	"os"

	"github.com/dotwallet/dotwallet-go-sdk"
)

func main() {

	// Create the DotWallet client
	c, err := dotwallet.NewClient(
		dotwallet.WithCredentials(
			os.Getenv("DOT_WALLET_CLIENT_ID"),
			os.Getenv("DOT_WALLET_CLIENT_SECRET"),
		),
		dotwallet.WithAutoLoadToken(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	var TransferNftToAddressData *dotwallet.TransferNftToAddressData
	if TransferNftToAddressData, err = c.TransferNftToAddress(
		"b0cc8ab416906ba6565aae3f575858c9bfb1654d9db1725edc786288584bba07",
		"1EszQWy21f4N77A2AVwzx4efGeiCRVnLoP",
		"name",
		"Desc",
		"https://img2.baidu.com/it/u=3895119537,2684520677&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500",
	); err != nil {
		log.Fatalln(err)
	}

	log.Println(
		"fee:", TransferNftToAddressData.Fee,
		"fee_str:", TransferNftToAddressData.FeeStr,
	)
}
