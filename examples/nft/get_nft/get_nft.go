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
	)
	if err != nil {
		log.Fatalln(err)
	}

	var NftData *dotwallet.NftData
	if NftData, err = c.GetNft(
		"36413e6c5d5955bd7090642b3f0c4b14b606489c318b943cc6510055fc421088",
	); err != nil {
		log.Fatalln(err)
	}

	log.Println(
		"code_hash:", NftData.CodeHash,
		"param:", NftData.Param,
	)
}
