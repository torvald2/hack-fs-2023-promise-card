package usecases

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/drand/tlock"
	"github.com/drand/tlock/networks/http"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/torvald2/hack-fs-2023-promise-card/ens"
	"github.com/torvald2/hack-fs-2023-promise-card/polybase"
)

type GetUserUseCase struct {
	timelockHost      string
	timelockChainHash string
	key               string
	url               string
	namespace         string
	ensRootOwner      string
	ensPrivateKey     string
	mainENS           string
	rpcUrl            string
	resolver          string
}

func NewGetUserUseCase(key string, url string, namespace string, tlUrl, tlCHash, pinataKey, ensRootOwner, ensPrivateKey, mainENS, rpcURL, resolver string) GetUserUseCase {
	return GetUserUseCase{
		key:               key,
		url:               url,
		namespace:         namespace,
		timelockHost:      tlUrl,
		timelockChainHash: tlCHash,
		ensRootOwner:      ensRootOwner,
		ensPrivateKey:     ensPrivateKey,
		mainENS:           mainENS,
		rpcUrl:            rpcURL,
		resolver:          resolver,
	}

}

func (c *GetUserUseCase) Execute(key, address string) (token string, err error) {
	network, err := http.NewNetwork(c.timelockHost, c.timelockChainHash)
	if err != nil {
		return
	}
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return
	}
	var data bytes.Buffer

	if err := tlock.New(network).Decrypt(&data, bytes.NewBuffer(keyBytes)); err != nil {
		return "", err
	}
	keyEncrypted, _ := ioutil.ReadAll(&data)
	usersKey, err := crypto.HexToECDSA(string(keyEncrypted[2:]))
	if err != nil {
		return
	}
	userAddress := crypto.PubkeyToAddress(usersKey.PublicKey).Hex()

	if userAddress != address {
		err = fmt.Errorf("Bad user address: %v", userAddress)
		return
	}

	cl, err := polybase.NewPolybaseClient(c.namespace, c.url)
	if err != nil {
		return
	}
	userData, err := cl.GetRecord("User", c.key, address)
	if err != nil {
		return
	}
	ensService := ens.ENSAdaptor{
		OwnerAddress:    c.ensRootOwner,
		PrivateKey:      c.ensPrivateKey,
		MainDomain:      c.mainENS,
		RPCUrl:          c.rpcUrl,
		ResolverAddress: c.resolver,
	}
	avatar, err := ensService.ResolveAvatar(address)
	if err != nil {
		return
	}
	userData["avatar"] = avatar
	privateKeyBytes := crypto.FromECDSA(usersKey)

	accessToken, err := CreateAccessToken(10*time.Minute, userData, privateKeyBytes, "promisecards")
	if err != nil {
		return
	}

	return accessToken, nil
}

func CreateAccessToken(ttl time.Duration, content interface{}, k []byte, iis string) (token string, err error) {

	key, err := crypto.ToECDSA(k)

	//key, err := jwt.ParseECPrivateKeyFromPEM(k)
	if err != nil {
		return
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content             // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.
	claims["iis"] = iis

	token, err = jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(key)
	return
}
