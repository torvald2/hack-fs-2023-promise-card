package storage

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	PublicKey string
	Roles     []string
	NickName  string
	StartsAt  time.Time
}

func (a *Account) CreateAddress() (key string, err error) {
	getPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}
	thePublicAddress := crypto.PubkeyToAddress(getPrivateKey.PublicKey).Hex()
	privateKeyBytes := crypto.FromECDSA(getPrivateKey)
	a.PublicKey = thePublicAddress
	return hexutil.Encode(privateKeyBytes), nil
}
