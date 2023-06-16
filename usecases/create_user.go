package usecases

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/drand/tlock"
	"github.com/drand/tlock/networks/http"
	"github.com/torvald2/hack-fs-2023-promise-card/ens"
	"github.com/torvald2/hack-fs-2023-promise-card/pinata"
	"github.com/torvald2/hack-fs-2023-promise-card/polybase"
	"github.com/torvald2/hack-fs-2023-promise-card/storage"
)

type CreateUserUseCase struct {
	key               string
	url               string
	namespace         string
	timelockHost      string
	timelockChainHash string
	encryptedKey      []byte
	address           string
	pinataKey         string
	ensRootOwner      string
	ensPrivateKey     string
	mainENS           string
	rpcUrl            string
	resolver          string
}

func NewCreateUserUseCase(key string, url string, namespace string, tlUrl, tlCHash, pinataKey, ensRootOwner, ensPrivateKey, mainENS, rpcURL, resolver string) CreateUserUseCase {
	return CreateUserUseCase{
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
		pinataKey:         pinataKey,
	}

}

func (c *CreateUserUseCase) Execute(nickName string, duration time.Duration, avatar string) error {
	usr := storage.Account{
		NickName: nickName,
	}
	privateKey, err := usr.CreateAddress()
	if err != nil {
		return err
	}
	cl, err := polybase.NewPolybaseClient(c.namespace, c.url)
	if err != nil {
		return err
	}
	args := make([]interface{}, 0)
	args = append(args, usr.PublicKey)
	args = append(args, nickName)
	_, err = cl.CreateRecord("User", args, c.key)
	if err != nil {
		return err
	}
	network, err := http.NewNetwork(c.timelockHost, c.timelockChainHash)
	if err != nil {
		return err
	}
	roundNumber := network.RoundNumber(time.Now().Add(duration))
	var cipherData bytes.Buffer

	if err := tlock.New(network).Encrypt(&cipherData, bytes.NewBuffer([]byte(privateKey)), roundNumber); err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(&cipherData)
	c.encryptedKey = data
	c.address = usr.PublicKey

	pinataService := pinata.New(c.pinataKey)

	image, err := base64.StdEncoding.DecodeString(avatar)
	if err != nil {
		return err
	}
	cid, err := pinataService.PinImage(image, "avatar")
	if err != nil {
		return err
	}

	ensService := ens.ENSAdaptor{
		OwnerAddress:    c.ensRootOwner,
		PrivateKey:      c.ensPrivateKey,
		MainDomain:      c.mainENS,
		RPCUrl:          c.rpcUrl,
		ResolverAddress: c.resolver,
	}

	_, err = ensService.CreateSubdomain(nickName, c.address)
	if err != nil {
		return err
	}
	_, err = ensService.CreateAvatar(fmt.Sprintf("ipfs://%s", cid), nickName)
	if err != nil {
		return err
	}
	return nil
}

func (c *CreateUserUseCase) GetPrivateKey() string {
	return hex.EncodeToString(c.encryptedKey)
}

func (c *CreateUserUseCase) GetAddress() string {
	return c.address
}
