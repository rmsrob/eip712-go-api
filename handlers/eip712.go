package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/labstack/echo/v4"
	"github.com/rrobrms/eip712-go-api/pkg/utils"
)

type SignParams struct {
	Address           string `json:"address"`
	ChainId           int64  `json:"chainid"`
	VerifyingContract string `json:"verifyingContract"`
	TokenId           string `json:"tokenId"`
	Amount            string `json:"amount"`
	Deadline          string `json:"deadline"`
}

type ReturnData struct {
	HashToSign common.Hash `json:"hashToSign"`
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
			"Action": {
				{Name: "action", Type: "string"},
				{Name: "tokenId", Type: "string"},
			},
			"NFT": {
				{Name: "action", Type: "Action"},
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
			"action":   map[string]interface{}{"action": "SellNFT", "tokenId": w3Data.TokenId},
			"owner":    w3Data.Address,
			"erc20":    "0xEeEEeeEeeeeEeEEeeEeeeeEeEEeeEeeeeEeEEeeE",
			"amount":   w3Data.Amount,
			"erc721":   "0xEeEEeeEeeeeEeEEeeEeeeeEeEEeeEeeeeEeEEeeE",
			"txHash":   "0xEeEEeeEeeeeEeEEeeEeeeeEeEEeeEeeeeEeEEeeE",
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
	log.Printf("hash to sign %#v\n", hashToSign)

	return c.JSON(http.StatusOK, ReturnData{
		HashToSign: hashToSign,
	})

}
