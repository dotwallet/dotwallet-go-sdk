package main

import (
	"log"
	"os"

	"github.com/dotwallet/dotwallet-go-sdk"
)

func main() {
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

	msgTx, err := c.GetMsgTxByStr("dbe0bd86245b42983058615d0249a4f9cd2dda49c5e369866d4c0f2e300993ad")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(c.SerializeRawTx(msgTx))
}
