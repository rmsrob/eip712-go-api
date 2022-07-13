package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
)

/**
 * @StoredHash will be used to retrive the publicKey of the user
 * @Address Fetch the previously stored challenge hash from your database
 * @UserSignature address of the OG signer to be verifyed
 **/
type VerifyParams struct {
	StoredHash    string         `json:"storedHash"`
	Address       common.Address `json:"address"`
	UserSignature string         `json:"userSignature"`
}

type ReturnVerifyData struct {
	Verify bool `json:"verify"`
}

func PostEIP712verify(c echo.Context) error {

	p := VerifyParams{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&p)
	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	log.Printf("all params: %#v\n", p)

	storedHash, err := hex.DecodeString(p.StoredHash)
	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	var userAddress = p.Address
	signature, _ := hex.DecodeString(p.UserSignature)

	if len(signature) != 65 {
		return fmt.Errorf("invalid signature length: %d", len(signature))
	}

	if signature[64] != 27 && signature[64] != 28 {
		return fmt.Errorf("invalid recovery id: %d", signature[64])
	}
	signature[64] -= 27

	pubKeyRaw, err := crypto.Ecrecover(storedHash, signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %s", err.Error())
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	if !bytes.Equal(userAddress.Bytes(), recoveredAddr.Bytes()) {
		return c.JSON(http.StatusOK, ReturnVerifyData{
			Verify: false,
		})
	} else {
		return c.JSON(http.StatusOK, ReturnVerifyData{
			Verify: true,
		})
	}
}
