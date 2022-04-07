package main

import (
	"encoding/json"
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
	// 1b7ebb4ff5ab15b13684e54f09f6dd0b10d0aeb7b890c0b31736f554085e5cd6 1NQLjZeZf72d7LFR9aYdZLxAZQVuvSPhy1 get 1b7ebb4ff5ab15b13684e54f09f6dd0b10d0aeb7b890c0b31736f554085e5cd6_1
	// b9f64416f03f0d7e8bb2dc11c2c5d16546e96c877c9751451d2ff1074cdade50 16o8pSZBXJCBPAPLvrT8KB7Q8rdhmwQBz get e5f4592c9bbb3d5eef652c1caf96710ed6b4c5c1e1bdfb188a64d73e83a2868e_67
	// 69f164bcfc834f77a7804b9373235cfe53a18b1897f1ddc83782fce38fdac6f8 no one get anything
	result, err := c.GetNftReceiveAddressesByTxidStr("1b7ebb4ff5ab15b13684e54f09f6dd0b10d0aeb7b890c0b31736f554085e5cd6")
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(b))
}
