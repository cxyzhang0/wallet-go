package test

import (
	tran "github.com/cxyzhang0/wallet-go/ethtran_pkcs11"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"math/big"
	"strconv"
	"testing"
)

// TestBuildTx
// It uses Blockdaemon ubiquity API to get the nonce, and gas related parameters.
// address1: 0x389ac41522E3019886ACB003843E62d84FfA70bB
// address3: 0xf9aD19e7a38FaDB98C8A5cC7a14aBcfE80AC657b
func TestBuildTx(t *testing.T) {
	pf := "Slot Token 0"
	from := kmssdk.KeyLabel{
		Prefix:    pf,
		KeyRing:   "WIM-test",
		Key:       "secp256k1-hsm-1",
		Version:   1,
		Algorithm: kmssdk.Secp256k1,
	}
	to := kmssdk.KeyLabel{
		Prefix:    pf,
		KeyRing:   "WIM-test",
		Key:       "secp256k1-hsm-1",
		Version:   3,
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

	balances, _, err := ubiAPIClient.AccountsAPI.GetListOfBalancesByAddress(ubiCtx, ubiPlatform, ubiNetwork, fromAddr).Execute()
	if err != nil {
		t.Errorf("failed to get balance via ubiquity: %+v", err)
	}
	if len(balances) == 0 {
		t.Errorf("no balance record")
	}
	balance, err := tran.StringToBigInt(balances[0].GetConfirmedBalance())
	if err != nil {
		t.Errorf("failed to pars big int: %+v", err)
	}
	//balanceFloat, _, err := big.ParseFloat(balances[0].GetConfirmedBalance(), 10, 0, big.ToZero)
	//if err != nil {
	//	t.Errorf("failed to pars big int: %+v", err)
	//}
	//balance, _ := balanceFloat.Int(nil)
	if balance.Cmp(req.Amount) != 1 {
		t.Errorf("balance %+v <= req amount %+v", balance, req.Amount)
	}
	req.Nonce = uint64(balances[0].GetConfirmedNonce())

	estimate, _, err := ubiAPIClient.TransactionsAPI.FeeEstimate(ubiCtx, ubiPlatform, ubiNetwork).Execute()
	if err != nil {
		t.Errorf("failed to get fee estimate via ubiquity: %+v", err)
	}

	mediumMap := estimate.EstimatedFees.Medium.(map[string]interface{})
	gasTipCapFloat, ok := mediumMap["max_priority_fee"]
	if !ok {
		t.Errorf("max_priority_fee is missing")
	}
	gasTipCap, err := tran.StringToBigInt(strconv.FormatFloat(gasTipCapFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Errorf("failed to convert gasTipCapFloat %s to big int", gasTipCapFloat)
	}
	req.GasTipCap = gasTipCap

	maxTotalFeeFloat, ok := mediumMap["max_total_fee"]
	if !ok {
		t.Errorf("max_total_fee is missing")
	}
	maxTotalFee, err := tran.StringToBigInt(strconv.FormatFloat(maxTotalFeeFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Errorf("failed to convert maxTotalFeeFloat %s to big int", maxTotalFeeFloat)
	}
	req.GasFeeCap = maxTotalFee

	rawSignedTx, txHash, err := tran.BuildTx(req, _sdk, chainConfig)
	if err != nil {
		t.Errorf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", rawSignedTx, txHash)

	//receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
	//t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
}
