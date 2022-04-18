package dotwallet

import (
	"bytes"
	"container/list"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
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
			Param:    Param,
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

// MintNft can obtain information authorized by DotWallet users via their user access_token
func (c *Client) TransferNftToAddress(Txid string, Address string, Name string, Desc string, PicUrl string) (*TransferNftToAddressData, error) {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		TransferNftToAddress,
		&transferNftToAddressParam{
			Txid:    Txid,
			Address: Address,
			Name:    Name,
			Desc:    Desc,
			PicUrl:  PicUrl,
		},
		http.StatusOK,
		c.Token(),
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	resp := new(transferNftToAddressResponse)
	if err = json.Unmarshal(
		response.Body, &resp,
	); err != nil {
		return nil, err
	}

	// Error?
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Message)
	}

	return &resp.Data.TransferNftToAddressData, nil
}

func ParseNftVoutScript(pkScript []byte) (btcutil.Address, error) {
	if len(pkScript) != NFT_VOUT_LEN {
		return nil, errors.New("not nft vout")
	}
	if !bytes.HasPrefix(pkScript, NFT_VOUT_SCRIPT_PREFIX) {
		return nil, errors.New("not nft vout")
	}
	addr, err := btcutil.NewAddressPubKeyHash(pkScript[NFT_VOUT_LEN-20:], &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (c *Client) VerifyCastingNftTransactionByRawTx(rawTx string) (bool, error) {
	msgTx, err := c.DeserializeRawTx(rawTx)
	if err != nil {
		return false, err
	}
	return c.VerifyCastingNftTransaction(msgTx)
}

func (c *Client) VerifyCastingNftTransactionByTxid(txid string) (bool, error) {
	msgTx, err := c.GetMsgTxByStr(txid)
	if err != nil {
		return false, err
	}
	return c.VerifyCastingNftTransaction(msgTx)
}

func (c *Client) VerifyCastingNftTransaction(msgTx *wire.MsgTx) (bool, error) {
	nftVinCount := 0
	for _, vin := range msgTx.TxIn {
		if !txscript.IsPushOnlyScript(vin.SignatureScript) {
			continue
		}

		datas, err := txscript.PushedData(vin.SignatureScript)
		if err != nil {
			return false, err
		}

		if len(datas) != NFT_VIN_PUSHED_DATA_COUNT {
			continue
		}

		preMsgTx, err := c.GetMsgTx(&vin.PreviousOutPoint.Hash)
		if err != nil {
			return false, err
		}
		preVout := preMsgTx.TxOut[vin.PreviousOutPoint.Index]
		_, err = ParseNftVoutScript(preVout.PkScript)
		if err != nil {
			continue
		}
		nftVinCount++
	}

	nftVoutCount := 0
	continuity := true
	var opReturnScript []byte = nil
	for index, vout := range msgTx.TxOut {
		if bytes.HasPrefix(vout.PkScript, []byte{0x00, 0x6a}) {
			opReturnScript = vout.PkScript
			continue
		}
		_, err := ParseNftVoutScript(vout.PkScript)
		if err != nil {
			continue
		}
		if nftVinCount != index {
			continuity = false
		}
		nftVoutCount++
	}
	if nftVinCount > 0 || nftVoutCount == 0 || !continuity || opReturnScript == nil {
		return false, nil
	}

	datas, err := txscript.PushedData(opReturnScript)
	if err != nil {
		return false, nil
	}
	if len(datas) != 2 {
		return false, nil
	}

	nftAuthInfo := &NftAuthInfo{}
	err = json.Unmarshal(datas[1], nftAuthInfo)
	if err != nil {
		return false, nil
	}

	sig, err := btcec.ParseSignature(nftAuthInfo.Sig, btcec.S256())
	if err != nil {
		panic(err)
	}
	pub, err := btcec.ParsePubKey(nftAuthInfo.Pub, btcec.S256())
	if err != nil {
		panic(err)
	}
	hasher := sha256.New()
	key := fmt.Sprintf("%s:%d", msgTx.TxIn[0].PreviousOutPoint.Hash.String(), msgTx.TxIn[0].PreviousOutPoint.Index)
	hasher.Write([]byte(key))
	keyHash := hasher.Sum(nil)
	ok := sig.Verify(
		keyHash,
		pub,
	)
	if !ok {
		return false, nil
	}
	return true, nil
}

func (c *Client) VerifyNftCastingOpReturn(msgTx *wire.MsgTx) bool {
	var opReturnScript []byte = nil
	for _, vout := range msgTx.TxOut {
		if bytes.HasPrefix(vout.PkScript, []byte{0x00, 0x6a}) {
			opReturnScript = vout.PkScript
			break
		}
	}
	datas, err := txscript.PushedData(opReturnScript)
	if err != nil {
		return false
	}
	if len(datas) != 2 {
		return false
	}

	nftAuthInfo := &NftAuthInfo{}
	err = json.Unmarshal(datas[1], nftAuthInfo)
	if err != nil {
		return false
	}

	sig, err := btcec.ParseSignature(nftAuthInfo.Sig, btcec.S256())
	if err != nil {
		return false
	}
	pub, err := btcec.ParsePubKey(nftAuthInfo.Pub, btcec.S256())
	if err != nil {
		return false
	}
	hasher := sha256.New()
	key := fmt.Sprintf("%s:%d", msgTx.TxIn[0].PreviousOutPoint.Hash.String(), msgTx.TxIn[0].PreviousOutPoint.Index)
	hasher.Write([]byte(key))
	keyHash := hasher.Sum(nil)
	ok := sig.Verify(
		keyHash,
		pub,
	)
	if !ok {
		return false
	}
	return true
}

func (c *Client) GetNftReceiveAddressesByTxidStr(txid string) ([]*AddressBadgeCodePair, error) {
	msgTx, err := c.GetMsgTxByStr(txid)
	if err != nil {
		return nil, err
	}
	return c.GetNftReceiveAddresses(msgTx)
}

func (c *Client) GetNftReceiveAddresses(msgTx *wire.MsgTx) ([]*AddressBadgeCodePair, error) {
	l := list.New()
	l.PushBack(msgTx)
	nftTxInfos := make([]*NftTxInfo, 0, 8)
	for l.Len() > 0 {
		elem := l.Front()
		l.Remove(elem)
		currentMsgTx := elem.Value.(*wire.MsgTx)
		nftTxInfo := &NftTxInfo{
			Txid:            currentMsgTx.TxHash().String(),
			NftPreOutPoints: make([]*TxidIndexPair, 0, 1),
			NftOutPoints:    make([]*AddressIndexPair, 0, 8),
			Type:            -1,
		}
		nftTxInfos = append(nftTxInfos, nftTxInfo)
		for _, vin := range currentMsgTx.TxIn {
			if !txscript.IsPushOnlyScript(vin.SignatureScript) {
				continue
			}

			datas, err := txscript.PushedData(vin.SignatureScript)
			if err != nil {
				return nil, err
			}

			if len(datas) != NFT_VIN_PUSHED_DATA_COUNT {
				continue
			}

			preMsgTx, err := c.GetMsgTx(&vin.PreviousOutPoint.Hash)
			if err != nil {
				return nil, err
			}
			preVout := preMsgTx.TxOut[vin.PreviousOutPoint.Index]
			_, err = ParseNftVoutScript(preVout.PkScript)
			if err != nil {
				continue
			}
			txidIndexPair := &TxidIndexPair{
				Txid:  vin.PreviousOutPoint.Hash.String(),
				Index: int(vin.PreviousOutPoint.Index),
			}
			nftTxInfo.NftPreOutPoints = append(nftTxInfo.NftPreOutPoints, txidIndexPair)
			l.PushBack(preMsgTx)
		}

		continuity := true
		nftVoutCount := 0
		for index, vout := range currentMsgTx.TxOut {
			addr, err := ParseNftVoutScript(vout.PkScript)
			if err != nil {
				continue
			}
			addressIndexPair := &AddressIndexPair{
				Address: addr.EncodeAddress(),
				Index:   index,
			}
			nftTxInfo.NftOutPoints = append(nftTxInfo.NftOutPoints, addressIndexPair)
			if nftVoutCount != index {
				continuity = false
			}
			nftVoutCount++
		}

		if len(nftTxInfo.NftPreOutPoints) == 0 && len(nftTxInfo.NftOutPoints) == 0 {
			nftTxInfo.Type = NFT_TX_TYPE_IRRELEVANT
			break
		}

		if len(nftTxInfo.NftPreOutPoints) == 0 && len(nftTxInfo.NftOutPoints) > 0 {
			//casting
			if !continuity || !c.VerifyNftCastingOpReturn(currentMsgTx) {
				nftTxInfo.Type = NFT_TX_TYPE_ERROR
				break
			}
			nftTxInfo.Type = NFT_TX_TYPE_CASTING
			break
		}

		if len(nftTxInfo.NftPreOutPoints) == 1 && len(nftTxInfo.NftOutPoints) == 0 {
			//destory
			nftTxInfo.Type = NFT_TX_TYPE_DESTORY
			break
		}

		if len(nftTxInfo.NftPreOutPoints) == 1 && len(nftTxInfo.NftOutPoints) == 1 {
			//transfer
			nftTxInfo.Type = NFT_TX_TYPE_TRANSFER
			continue
		}

		if len(nftTxInfo.NftPreOutPoints) > 1 {
			//destory
			nftTxInfo.Type = NFT_TX_TYPE_DESTORY
			break
		}
	}
	nftTxInfosCount := len(nftTxInfos)

	if nftTxInfosCount == 0 {
		return nil, errors.New("nftTxInfosCount should not be zero")
	}

	firstNftTxInfo := nftTxInfos[nftTxInfosCount-1]
	//追回去的第一笔
	if firstNftTxInfo.Type != NFT_TX_TYPE_CASTING {
		return []*AddressBadgeCodePair{}, nil
	}

	if nftTxInfosCount == 1 {
		result := make([]*AddressBadgeCodePair, 0, len(firstNftTxInfo.NftOutPoints))
		for _, nftOutPoint := range firstNftTxInfo.NftOutPoints {
			addressBadgeCodePair := &AddressBadgeCodePair{
				Address:   nftOutPoint.Address,
				BadgeCode: fmt.Sprintf("%s_%d", firstNftTxInfo.Txid, nftOutPoint.Index+1),
			}
			result = append(result, addressBadgeCodePair)
		}
		return result, nil
	}

	secondNftTxInfo := nftTxInfos[nftTxInfosCount-2]
	if len(secondNftTxInfo.NftPreOutPoints) != 1 {
		return nil, errors.New("secondNftTxInfo.NftPreOutPoints should be 1")
	}

	badgeCode := fmt.Sprintf("%s_%d", firstNftTxInfo.Txid, secondNftTxInfo.NftPreOutPoints[0].Index+1)

	lastNftTxInfo := nftTxInfos[0]

	if len(lastNftTxInfo.NftOutPoints) == 0 {
		return nil, errors.New("lastNftTxInfo.NftOutPoints count should be zero")
	}

	return []*AddressBadgeCodePair{
		{
			Address:   lastNftTxInfo.NftOutPoints[0].Address,
			BadgeCode: badgeCode,
		},
	}, nil

}
