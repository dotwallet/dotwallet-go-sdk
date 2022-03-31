package dotwallet

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetNft can obtain information authorized by DotWallet users via their user access_token
func (c *Client) GetNft(Txid string) (*NftData, error) {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		getNft,
		&getNftParam{
			Txid: Txid,
		},
		http.StatusOK,
		c.Token(),
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	resp := new(nftResponse)
	if err = json.Unmarshal(
		response.Body, &resp,
	); err != nil {
		return nil, err
	}

	// Error?
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Message)
	}

	return &resp.Data.NftData, nil
}

// MintNft can obtain information authorized by DotWallet users via their user access_token
func (c *Client) MintNft(CodeHash string, Param string) (*NftMintData, error) {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		mintNft,
		&mintNftParam{
			CodeHash: CodeHash,
			Param: Param,
		},
		http.StatusOK,
		c.Token(),
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	resp := new(nftMintResponse)
	if err = json.Unmarshal(
		response.Body, &resp,
	); err != nil {
		return nil, err
	}

	// Error?
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Message)
	}

	return &resp.Data.NftMintData, nil
}
