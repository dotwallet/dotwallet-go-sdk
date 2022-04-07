package dotwallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func (c *Client) GetMerkleProof(txhash *chainhash.Hash) (*MerkelProof, error) {
	return c.GetMerkleProofByHashStr(txhash.String())
}

func (c *Client) GetMerkleProofByHashStr(txid string) (*MerkelProof, error) {
	getMerkleProofRequest := &GetMerkleProofRequest{
		Txid: txid,
	}
	response, err := c.Request(http.MethodPost, getMerkleProof, getMerkleProofRequest, http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	getMerkleProofResponse := &GetMerkleProofResponse{}
	err = json.Unmarshal(response.Body, getMerkleProofResponse)
	if err != nil {
		return nil, err
	}
	return &getMerkleProofResponse.Data.MerkelProof, nil
}

func (c *Client) VerifyMerkleProof(txid string, nIndex int, nodes []string, exceptMerkleRoot string) bool {
	hash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return false
	}
	for i := 0; i < len(nodes); i, nIndex = i+1, nIndex>>1 {
		node := nodes[i]
		if node == "*" {
			hash = blockchain.HashMerkleBranches(hash, hash)
			continue
		}
		peerHash, err := chainhash.NewHashFromStr(node)
		if err != nil {
			return false
		}
		if nIndex&1 != 0 {
			hash = blockchain.HashMerkleBranches(peerHash, hash)
			continue
		}
		hash = blockchain.HashMerkleBranches(hash, peerHash)
	}
	return hash.String() == exceptMerkleRoot
}

func (c *Client) DeserializeRawTx(rawTx string) (*wire.MsgTx, error) {
	serializedTx, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, err
	}
	msgtx := wire.NewMsgTx(2)
	err = msgtx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		return nil, err
	}
	return msgtx, nil
}

func (c *Client) SerializeRawTx(tx *wire.MsgTx) string {
	buf := make([]byte, 0, tx.SerializeSize())
	buff := bytes.NewBuffer(buf)
	tx.Serialize(buff)
	rawtxByte := buff.Bytes()
	rawtx := hex.EncodeToString(rawtxByte)
	return rawtx
}

func (c *Client) GetMsgTxByStr(txid string) (*wire.MsgTx, error) {
	getRawtransactionRequest := &GetRawTransactionRequest{
		Txid: txid,
	}
	response, err := c.Request(http.MethodPost, getRawtransaction, getRawtransactionRequest, http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	getRawTransactionResponse := &GetRawTransactionResponse{}
	err = json.Unmarshal(response.Body, getRawTransactionResponse)
	if err != nil {
		return nil, err
	}
	return c.DeserializeRawTx(getRawTransactionResponse.Data.RawTx)
}

func (c *Client) GetMsgTx(hash *chainhash.Hash) (*wire.MsgTx, error) {
	return c.GetMsgTxByStr(hash.String())
}
