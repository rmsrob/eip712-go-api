package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/labstack/echo/v4"
	"github.com/rrobrms/eip712-go-api/pkg/utils"
)

type SignParams struct {
	Action            string `json:"action"`
	Address           string `json:"address"`
	ChainId           int64  `json:"chainid"`
	VerifyingContract string `json:"verifyingContract"`
	TokenId           string `json:"tokenId"`
	BuyPrice          string `json:"buyPrice"`
	Amount            string `json:"amount"`
	Deadline          string `json:"deadline"`
	TxHash            string `json:"txHash"`
}

type ReturnData struct {
	HashToSign common.Hash `json:"hashToSign"`
	Signature  string      `json:"signature"`
}

func PostEIP712(c echo.Context) error {

	w3Data := SignParams{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&w3Data)
	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	log.Printf("all params: %#v\n", w3Data)

	testTypedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": utils.EIP712DomainType,
			// "Action": {
			// 	{Name: "action", Type: "string"},
			// 	{Name: "tokenId", Type: "string"},
			// 	{Name: "buyPrice", Type: "string"},
			// },
			"NFT": {
				{Name: "action", Type: "string"},
				{Name: "tokenId", Type: "string"},
				{Name: "buyPrice", Type: "string"},
				{Name: "owner", Type: "string"},
				{Name: "erc20", Type: "string"},
				{Name: "amount", Type: "string"},
				{Name: "erc721", Type: "string"},
				{Name: "txHash", Type: "string"},
				{Name: "deadline", Type: "string"},
			},
		},
		PrimaryType: "NFT",
		Domain: apitypes.TypedDataDomain{
			Name:              "nftsupermarket.eth",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(w3Data.ChainId),
			VerifyingContract: w3Data.VerifyingContract,
		},
		Message: apitypes.TypedDataMessage{
			// "action":   map[string]interface{}{"action": w3Data.Action, "tokenId": w3Data.TokenId, "buyPrice": w3Data.BuyPrice},
			"action": w3Data.Action,
			"tokenId": w3Data.TokenId, 
			"buyPrice": w3Data.BuyPrice,
			"owner":    w3Data.Address,
			"erc20":    "0x4A8b871784A8e6344126F47d48283a87Ea987f27",
			"amount":   w3Data.Amount,
			"erc721":   "0xB9E431Fc34152246BB28453b6ce117829E8A5B0C",
			"txHash":   w3Data.TxHash,
			"deadline": w3Data.Deadline,
		},
	}

	encoded, err := utils.EIP712Encode(&testTypedData)
	if err != nil {
		log.Fatalf("Failed to encode typed data %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	log.Printf("encoded data %#v\n", encoded)

	hashToSign, err := utils.HashForSigning(&testTypedData)
	if err != nil {
		log.Fatalf("Failed to hash typed data %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	log.Printf("hash to sign and stored in db %#v\n", hashToSign)

	// TODO: store hashToSign in db with other data needed.

	// ! Sign the hash to simulate the signature of the user.
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal(err)
	}

	signature, err := crypto.Sign(hashToSign.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	signed := hexutil.Encode(signature)
	// ! end of signing.

	return c.JSON(http.StatusOK, ReturnData{
		HashToSign: hashToSign,
		Signature:  signed,
	})

}
