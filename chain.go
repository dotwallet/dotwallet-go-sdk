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

// GetMerkleProof get merkle proof
func (c *Client) GetMerkleProof(txHash *chainhash.Hash) (*MerkelProof, error) {
	return c.GetMerkleProofByHashStr(txHash.String())
}

// GetMerkleProofByHashStr get merkle proof by hash string
func (c *Client) GetMerkleProofByHashStr(txID string) (*MerkelProof, error) {
	getMerkleProofRequest := &GetMerkleProofRequest{
		Txid: txID,
	}
	response, err := c.Request(http.MethodPost, getMerkleProof, getMerkleProofRequest, http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	getMerkleProofResponse := &GetMerkleProofResponse{}
	if err = json.Unmarshal(
		response.Body, getMerkleProofResponse,
	); err != nil {
		return nil, err
	}
	return &getMerkleProofResponse.Data.MerkelProof, nil
}

// VerifyMerkleProof verify merkle proof
func (c *Client) VerifyMerkleProof(txID string, nIndex int, nodes []string, exceptMerkleRoot string) bool {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return false
	}
	var peerHash *chainhash.Hash
	for i := 0; i < len(nodes); i, nIndex = i+1, nIndex>>1 {
		node := nodes[i]
		if node == "*" {
			hash = blockchain.HashMerkleBranches(hash, hash)
			continue
		}
		if peerHash, err = chainhash.NewHashFromStr(node); err != nil {
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

// DeserializeRawTx deserialize the raw tx
func (c *Client) DeserializeRawTx(rawTx string) (*wire.MsgTx, error) {
	serializedTx, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, err
	}
	msgTx := wire.NewMsgTx(2)
	if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); err != nil {
		return nil, err
	}
	return msgTx, nil
}

// SerializeRawTx serialize the raw tx
func (c *Client) SerializeRawTx(tx *wire.MsgTx) string {
	buf := make([]byte, 0, tx.SerializeSize())
	buff := bytes.NewBuffer(buf)
	_ = tx.Serialize(buff) // TODO: capture this error?
	return hex.EncodeToString(buff.Bytes())
}

// GetMsgTxByStr get message tx by string
func (c *Client) GetMsgTxByStr(txID string) (*wire.MsgTx, error) {
	getRawTransactionRequest := &GetRawTransactionRequest{
		TxID: txID,
	}
	response, err := c.Request(
		http.MethodPost, getRawTransaction,
		getRawTransactionRequest, http.StatusOK, nil,
	)
	if err != nil {
		return nil, err
	}
	getRawTransactionResponse := &GetRawTransactionResponse{}
	if err = json.Unmarshal(
		response.Body, getRawTransactionResponse,
	); err != nil {
		return nil, err
	}
	return c.DeserializeRawTx(getRawTransactionResponse.Data.RawTx)
}

// GetMsgTx get message tx
func (c *Client) GetMsgTx(hash *chainhash.Hash) (*wire.MsgTx, error) {
	return c.GetMsgTxByStr(hash.String())
}
