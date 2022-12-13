package algotran_pkcs11

import (
	"crypto/ed25519"
	"github.com/algorand/go-algorand-sdk/crypto"
	transaction "github.com/algorand/go-algorand-sdk/future"
	//"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
)

type TxReq struct {
	From        kmssdk.KeyLabel
	FromPubKey  *ed25519.PublicKey
	FromAddress *types.Address
	FromAddr    string
	To          kmssdk.KeyLabel
	ToAddr      string
	Amount      uint64
	TxParams    *types.SuggestedParams
	Note        []byte
	AssetID     uint64
}

func BuildPaymentTx(req TxReq, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	tx, err := transaction.MakePaymentTxn(
		req.FromAddr,
		req.ToAddr,
		req.Amount,
		req.Note,
		"",
		*req.TxParams,
	)
	/*
			note := req.Note //[]byte("Hello Payment")
			var minFee uint64 = transaction.MinTxnFee
			genID := req.TxParams.GenesisID
			genHash := req.TxParams.GenesisHash
			firstValidRound := uint64(req.TxParams.FirstRoundValid)
			lastValidRound := uint64(req.TxParams.LastRoundValid)
		tx, err := transaction.MakePaymentTxnWithFlatFee(
			req.FromAddr,
			req.ToAddr,
			minFee,
			req.Amount,
			firstValidRound,
			lastValidRound,
			note,
			"",
			genID,
			genHash,
			)
	*/
	if err != nil {
		return
	}

	txid, stxBytes, err = SignTransaction(req, tx, sdk)

	//crypto.SignMultisigTransaction()
	//crypto.SignTransaction()
	//types.SignedTxn
	//types.Signature
	return
}

func BuildAssetAcceptanceTx(req TxReq, assetID uint64, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	tx, err := transaction.MakeAssetAcceptanceTxn(
		req.FromAddr,
		req.Note,
		*req.TxParams,
		assetID,
	)
	if err != nil {
		return
	}

	txid, stxBytes, err = SignTransaction(req, tx, sdk)

	//crypto.SignMultisigTransaction()
	//crypto.SignTransaction()
	//types.SignedTxn
	//types.Signature
	return
}

type MultisigTxReq struct {
	From        []kmssdk.KeyLabel
	PubKeys     []ed25519.PublicKey
	FromAddress *types.Address
	FromAddr    string
	M           int
	To          kmssdk.KeyLabel
	ToAddr      string
	Amount      uint64
	TxParams    *types.SuggestedParams
	Note        []byte
	AssetID     uint64
}

func BuildMultisigPaymentTx(req MultisigTxReq, ma crypto.MultisigAccount, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	tx, err := transaction.MakePaymentTxn(
		req.FromAddr,
		req.ToAddr,
		req.Amount,
		req.Note,
		"",
		*req.TxParams,
	)
	/*
		var minFee uint64 = transaction.MinTxnFee
		note := req.Note //[]byte("Hello Payment")
		genID := req.TxParams.GenesisID
		genHash := req.TxParams.GenesisHash
		firstValidRound := uint64(req.TxParams.FirstRoundValid)
		lastValidRound := uint64(req.TxParams.LastRoundValid)
		tx, err := transaction.MakePaymentTxn(
			req.FromAddr,
			req.ToAddr,
			minFee,
			req.Amount,
			firstValidRound,
			lastValidRound,
			note,
			"",
			genID,
			genHash,
		)
	*/
	if err != nil {
		return
	}

	return makeMultisig(req, ma, tx, sdk)
	/*
		txid, stxBytes, err = SignMultisigTransaction(req.From[0], req.PubKeys[0], ma, tx, sdk)
		if err != nil {
			return "", nil, err
		}
		fmt.Printf("partially signed multsig transaction txid %d: %s", 0, txid)

		for i := 1; i < req.M; i++ {
			txid, stxBytes, err = AppendMultisigTransaction(req.From[i], req.PubKeys[i], ma, stxBytes, sdk)
			if err != nil {
				return "", nil, err
			}
			fmt.Printf("partially signed multsig transaction txid %d: %s", i, txid)
		}

		return
	*/
}

type AssetParam struct {
	CreatorAddress  string
	AssetName       string
	UnitName        string
	Total           uint64
	Decimals        uint32
	DefaultFrozen   bool
	URL             string
	MetaDataHash    string
	ManagerAddress  string
	ReserveAddress  string
	FreezeAddress   string
	ClawbackAddress string
}

func BuildMultisigCreateAssetTx(req MultisigTxReq, ma crypto.MultisigAccount, param AssetParam, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	tx, err := transaction.MakeAssetCreateTxn(
		param.CreatorAddress,
		req.Note,
		*req.TxParams,
		param.Total,
		param.Decimals,
		param.DefaultFrozen,
		param.ManagerAddress,
		param.ReserveAddress,
		param.FreezeAddress,
		param.ClawbackAddress,
		param.UnitName,
		param.AssetName,
		param.URL,
		param.MetaDataHash,
	)
	if err != nil {
		return
	}

	return makeMultisig(req, ma, tx, sdk)
}

// transfer from Reserve address to recipient
func BuildMultisigMintAssetTx(req MultisigTxReq, ma crypto.MultisigAccount, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	tx, err := transaction.MakeAssetTransferTxn(
		req.FromAddr,
		req.ToAddr,
		req.Amount,
		req.Note,
		*req.TxParams,
		"",
		req.AssetID,
	)
	if err != nil {
		return
	}

	return makeMultisig(req, ma, tx, sdk)
}

// transfer from a normal address
func BuildTransferAssetTx(req TxReq, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	tx, err := transaction.MakeAssetTransferTxn(
		req.FromAddr,
		req.ToAddr,
		req.Amount,
		req.Note,
		*req.TxParams,
		"",
		req.AssetID,
	)
	if err != nil {
		return
	}

	txid, stxBytes, err = SignTransaction(req, tx, sdk)
	return
}
