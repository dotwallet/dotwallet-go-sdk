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

	var NftMintData *dotwallet.NftMintData
	if NftMintData, err = c.MintNft(
		"094dc7fa82e3ccb0e112c5f58f220cddbcd97811ff1bfc4ff349a3ec7862379e",
		"c34e415745e1c2a7f3ff91ea9a0f24e3ba829e718f9c821dffce9bb60aeb3691",
	); err != nil {
		log.Fatalln(err)
	}

	log.Println(
		"fee:", NftMintData.Fee,
		"fee_str:", NftMintData.FeeStr,
		"txid:", NftMintData.Txid,
	)
}
