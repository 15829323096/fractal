package sdk

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/fractalplatform/fractal/accountmanager"
	"github.com/fractalplatform/fractal/common"
	"github.com/fractalplatform/fractal/consensus/dpos"
	"github.com/fractalplatform/fractal/crypto"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	// 	tAssetName   = "testasset"
	// 	tAssetSymbol = "tat"
	// 	tAmount      = new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1e8))
	// 	tDecimals    = uint64(8)
	// 	tAssetID     uint64
	rpchost         = "http://127.0.0.1:8545"
	systemaccount   = "ftsystemio"
	systemprivkey   = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	systemassetname = "ftoken"
	systemassetid   = uint64(1)
	chainid         = big.NewInt(1)
	tValue          = new(big.Int).Mul(big.NewInt(300000), big.NewInt(1e18))
	tGas            = uint64(90000)
)

func TestAccount(t *testing.T) {
	Convey("Account", t, func() {
		api := NewAPI(rpchost)
		var systempriv, _ = crypto.HexToECDSA(systemprivkey)
		sysAcct := NewAccount(api, common.StrToName(systemaccount), systempriv, systemassetid, math.MaxUint64, true, chainid)
		// CreateAccount
		priv, pub := GenerateKey()
		accountName := common.StrToName(GenerateAccountName("test", 8))
		hash, err := sysAcct.CreateAccount(common.StrToName(systemaccount), tValue, systemassetid, tGas, &accountmanager.AccountAction{
			AccountName: accountName,
			PublicKey:   pub,
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// Transfer
		hash, err = sysAcct.Transfer(accountName, tValue, systemassetid, tGas)
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// UpdateAccount
		acct := NewAccount(api, accountName, priv, systemassetid, math.MaxUint64, true, chainid)
		_, npub := GenerateKey()
		hash, err = acct.UpdateAccount(common.StrToName(systemaccount), new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &accountmanager.AccountAction{
			AccountName: accountName,
			PublicKey:   npub,
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// DestroyAccount
	})
}

func TestAsset(t *testing.T) {
	// Convey("types.IssueAsset", t, func() {
	// 	api := NewAPI(rpchost)
	// 	var systempriv, _ = crypto.HexToECDSA(systemprivkey)
	// 	sysAcct := NewAccount(api, common.StrToName(systemaccount), systempriv, systemassetid, math.MaxUint64, true, chainid)
	// 	priv, pub := GenerateKey()
	// 	accountName := common.StrToName(GenerateAccountName("test", 8))
	// 	hash, err := sysAcct.CreateAccount(common.StrToName(systemaccount), tValue, systemassetid, tGas, &accountmanager.AccountAction{
	// 		AccountName: accountName,
	// 		PublicKey:   pub,
	// 	})
	// 	So(err, ShouldBeNil)
	// 	So(hash, ShouldNotBeNil)

	// 	acct := NewAccount(api, accountName, priv, systemassetid, math.MaxUint64, true, chainid)
	// 	assetname := common.StrToName(GenerateAccountName("asset", 8)).String()
	// 	// IssueAsset
	// 	hash, err = acct.IssueAsset(accountName, new(big.Int).Div(tValue, big.NewInt(10)), systemassetid, tGas, &asset.AssetObject{
	// 		AssetName:  assetname,
	// 		Symbol:     assetname[len(assetname)-4:],
	// 		Amount:     new(big.Int).Mul(big.NewInt(10000000), big.NewInt(1e18)),
	// 		Decimals:   18,
	// 		Owner:      accountName,
	// 		Founder:    accountName,
	// 		AddIssue:   big.NewInt(0),
	// 		UpperLimit: big.NewInt(0),
	// 	})
	// 	So(err, ShouldBeNil)
	// 	So(hash, ShouldNotBeNil)

	// 	// acct.UpdateAsset()
	// 	// acct.IncreaseAsset()
	// 	// acct.SetAssetOwner()
	// 	// acct.DestroyAsset()
	// })
}

func TestDPOS(t *testing.T) {
	Convey("DPOS", t, func() {
		api := NewAPI(rpchost)
		var systempriv, _ = crypto.HexToECDSA(systemprivkey)
		sysAcct := NewAccount(api, common.StrToName(systemaccount), systempriv, systemassetid, math.MaxUint64, true, chainid)
		priv, pub := GenerateKey()
		accountName := common.StrToName(GenerateAccountName("prod", 8))
		hash, err := sysAcct.CreateAccount(common.StrToName(systemaccount), tValue, systemassetid, tGas, &accountmanager.AccountAction{
			AccountName: accountName,
			PublicKey:   pub,
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		priv2, pub2 := GenerateKey()
		accountName2 := common.StrToName(GenerateAccountName("voter", 8))
		hash, err = sysAcct.CreateAccount(common.StrToName(systemaccount), tValue, systemassetid, tGas, &accountmanager.AccountAction{
			AccountName: accountName2,
			PublicKey:   pub2,
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// RegCadidate
		acct := NewAccount(api, accountName, priv, systemassetid, math.MaxUint64, true, chainid)
		acct2 := NewAccount(api, accountName2, priv2, systemassetid, math.MaxUint64, true, chainid)
		hash, err = acct.RegCadidate(accountName, new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &dpos.RegisterCadidate{
			Url:   fmt.Sprintf("www.%s.com", accountName.String()),
			Stake: new(big.Int).Div(tValue, big.NewInt(3)),
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// VoteCadidate
		hash, err = acct2.VoteCadidate(accountName, new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &dpos.VoteCadidate{
			Cadidate: accountName.String(),
			Stake:    new(big.Int).Div(tValue, big.NewInt(3)),
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// UnvoteCadidate
		hash, err = acct2.UnvoteCadidate(accountName, new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas)
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// VoteCadidate
		hash, err = acct2.VoteCadidate(accountName, new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &dpos.VoteCadidate{
			Cadidate: systemaccount,
			Stake:    new(big.Int).Div(tValue, big.NewInt(3)),
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// ChangeCadidate
		hash, err = acct2.ChangeCadidate(accountName, new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &dpos.ChangeCadidate{
			Cadidate: accountName.String(),
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// UnvoteVoter
		hash, err = acct.UnvoteVoter(accountName, new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &dpos.RemoveVoter{
			Voter: accountName2.String(),
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		hash, err = sysAcct.KickedCadidate(common.StrToName(systemaccount), new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &dpos.KickedCadidate{
			Cadidates: []string{accountName.String()},
		})
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)

		// UnRegCadidate
		hash, err = acct.UnRegCadidate(accountName, new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas)
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)
	})
}
func TestManual(t *testing.T) {
	Convey("Manual", t, func() {
		api := NewAPI(rpchost)
		var systempriv, _ = crypto.HexToECDSA(systemprivkey)
		sysAcct := NewAccount(api, common.StrToName(systemaccount), systempriv, systemassetid, math.MaxUint64, true, chainid)

		hash, err := sysAcct.KickedCadidate(common.StrToName(systemaccount), new(big.Int).Mul(tValue, big.NewInt(0)), systemassetid, tGas, &dpos.KickedCadidate{
			Cadidates: []string{"ftcadidate1", "ftcadidate2", "ftcadidate3"},
			Invalid:   true,
		})
		fmt.Println("====", err, hash.String())
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeNil)
	})
}
