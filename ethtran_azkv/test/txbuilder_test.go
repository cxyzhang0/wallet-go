package test

import (
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	tran "github.com/cxyzhang0/wallet-go/ethtran_azkv"
	ubi "gitlab.com/Blockdaemon/ubiquity/ubiquity-go-client/v1/pkg/client"
	"math/big"
	"strconv"
	"testing"
)

// TestBuildTx
// It uses Blockdaemon ubiquity API to get the nonce, and gas related parameters.
// address1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
// address2: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
func TestBuildTx(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	to := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	req := tran.TxReq{
		From:   from,
		To:     to,
		Amount: big.NewInt(1e15),
		Gas:    21000, // for standard tx transfering ether from one address to another
		// the following will be set be code
		Nonce:     0,
		GasTipCap: big.NewInt(3200),
		GasFeeCap: big.NewInt(32000),
	}

	_, fromAddr, err := tran.GetAddressPubKey(req.From, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pubkey: %+v", err)
	}

	balances, _, err := ubiAPIClient.AccountsAPI.GetListOfBalancesByAddress(ubiCtx, ubiPlatform, ubiNetwork, fromAddr).Execute()
	if err != nil {
		t.Fatalf("failed to get balance via ubiquity: %+v", err)
	}
	if len(balances) == 0 {
		t.Fatalf("no balance record")
	}
	balance, err := tran.StringToBigInt(balances[0].GetConfirmedBalance())
	if err != nil {
		t.Fatalf("failed to pars big int: %+v", err)
	}

	if balance.Cmp(req.Amount) != 1 {
		t.Fatalf("balance %+v <= req amount %+v", balance, req.Amount)
	}

	req.Nonce = uint64(balances[0].GetConfirmedNonce())

	estimate, _, err := ubiAPIClient.TransactionsAPI.FeeEstimate(ubiCtx, ubiPlatform, ubiNetwork).Execute()
	if err != nil {
		t.Fatalf("failed to get fee estimate via ubiquity: %+v", err)
	}

	mediumMap := estimate.EstimatedFees.Medium.(map[string]interface{})
	gasTipCapFloat, ok := mediumMap["max_priority_fee"]
	if !ok {
		t.Fatalf("max_priority_fee is missing")
	}
	gasTipCap, err := tran.StringToBigInt(strconv.FormatFloat(gasTipCapFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert gasTipCapFloat %s to big int", gasTipCapFloat)
	}
	req.GasTipCap = gasTipCap

	maxTotalFeeFloat, ok := mediumMap["max_total_fee"]
	if !ok {
		t.Fatalf("max_total_fee is missing")
	}
	maxTotalFee, err := tran.StringToBigInt(strconv.FormatFloat(maxTotalFeeFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert maxTotalFeeFloat %s to big int", maxTotalFeeFloat)
	}
	req.GasFeeCap = maxTotalFee

	rawSignedTx, txHash, err := tran.BuildTx(req, _sdk, chainConfig)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", rawSignedTx, txHash)

	receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
	if err != nil {
		t.Fatalf("failed to send tx %s \ntx hash %s\n %+v", rawSignedTx, txHash, err)
	}

	t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
}
