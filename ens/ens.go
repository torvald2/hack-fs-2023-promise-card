package ens

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

type ENSAdaptor struct {
	OwnerAddress    string
	PrivateKey      string
	MainDomain      string
	RPCUrl          string
	ResolverAddress string
}

func (e *ENSAdaptor) CreateSubdomain(subdomain, receiver string) (string, error) {
	client, err := ethclient.Dial(e.RPCUrl)
	if err != nil {
		return "", err
	}
	registry, err := ens.NewRegistry(client)
	if err != nil {
		return "", err
	}
	opts, err := e.getTxOptions(client)
	if err != nil {
		return "", err
	}
	receiverAddress := common.HexToAddress(receiver)
	resolverAddress := common.HexToAddress(e.ResolverAddress)

	tx, err := registry.SetSubdomainOwner(opts, e.MainDomain, subdomain, receiverAddress)
	if err != nil {
		return "", err
	}

	tx, err = registry.SetResolver(opts, fmt.Sprintf("%s.%s", subdomain, e.MainDomain), resolverAddress)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (e *ENSAdaptor) CreateAvatar(avatarUrl, nick string) (string, error) {
	client, err := ethclient.Dial(e.RPCUrl)
	if err != nil {
		return "", err
	}

	opts, err := e.getTxOptions(client)
	if err != nil {
		return "", err
	}
	resolver, err := ens.NewResolver(client, e.MainDomain)
	if err != nil {
		return "", err
	}
	tx, err := resolver.SetText(opts, "avatar", avatarUrl)

	return tx.Hash().Hex(), err
}
func (e *ENSAdaptor) ResolveAvatar(address string) (avatar string, err error) {
	client, err := ethclient.Dial(e.RPCUrl)
	if err != nil {
		return
	}

	resolver, err := ens.NewResolver(client, e.MainDomain)
	if err != nil {
		return
	}
	avatar, err = resolver.Text("avatar")
	if err != nil {
		return
	}

	return
}

func (e *ENSAdaptor) getTxOptions(client *ethclient.Client) (*bind.TransactOpts, error) {
	from := common.HexToAddress(e.OwnerAddress)
	key, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		return &bind.TransactOpts{}, err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return &bind.TransactOpts{}, err

	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return &bind.TransactOpts{}, err
	}
	gasPrice = big.NewInt(1502287099)

	signer := KeySigner(chainID, key)

	return &bind.TransactOpts{
		From:     from,
		Signer:   signer,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		GasLimit: 390000,
	}, nil
}

func KeySigner(chainID *big.Int, key *ecdsa.PrivateKey) (signerfn bind.SignerFn) {
	signerfn = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		keyAddr := crypto.PubkeyToAddress(key.PublicKey)
		if address != keyAddr {
			return nil, errors.New("not authorized to sign this account")
		}
		return types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	}

	return
}
func NameHash(name string) common.Hash {
	node := common.Hash{}

	if len(name) > 0 {
		labels := strings.Split(name, ".")

		for i := len(labels) - 1; i >= 0; i-- {
			labelSha := crypto.Keccak256Hash([]byte(labels[i]))
			node = crypto.Keccak256Hash(node.Bytes(), labelSha.Bytes())
		}
	}

	return node
}
