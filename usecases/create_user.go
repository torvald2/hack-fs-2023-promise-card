package usecases

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"time"

	"github.com/drand/tlock"
	"github.com/drand/tlock/networks/http"
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
}

func NewCreateUserUseCase(key string, url string, namespace string, tlUrl, tlCHash string) CreateUserUseCase {
	return CreateUserUseCase{
		key:               key,
		url:               url,
		namespace:         namespace,
		timelockHost:      tlUrl,
		timelockChainHash: tlCHash,
	}

}

func (c *CreateUserUseCase) Execute(nickName string, duration time.Duration) error {
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

	return nil
}

func (c *CreateUserUseCase) GetPrivateKey() string {
	return hex.EncodeToString(c.encryptedKey)
}

func (c *CreateUserUseCase) GetAddress() string {
	return c.address
}
