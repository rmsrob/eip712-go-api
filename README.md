# EIP 712 - Go api - example

> We are creating the hash to be signed on the user wallet and can verify it.

## Usage

```sh
go run main.go
```

> POST request

```sh
curl --request POST \
  --url http://localhost:1337/ \
  --header 'Content-Type: application/json' \
  --data '{
	"address": "0xEeEEeeEeeeeEeEEeeEeeeeEeEEeeEeeeeEeEEeeE",
	"chainid": 1,
	"verifyingContract": "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
	"tokenId": "138",
	"amount": "222222",
	"deadline": "19780"
}
'

### resp
{
	"hashToSign": "0xc1d79dbec465caddf6f578c56a93b67999de931428d1a96eb8ce3e26795b7e95"
}
```
