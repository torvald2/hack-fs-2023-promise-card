package ens

import "testing"

func TestCreateSubdomain(t *testing.T) {
	ensService := ENSAdaptor{
		OwnerAddress:    "0x55D3a2918CBc15467Cba3c965d0Bb89352B047c6",
		PrivateKey:      "eb9eb40f944c4b2d3d8aa16a9e89e3f933ac6fc0e20513901cacb4bbc73548ea",
		MainDomain:      "promisecard.eth",
		RPCUrl:          "https://eth-goerli.g.alchemy.com/v2/wnn2ogm_fc2xS605Ja6LaINbBuYWovzb",
		ResolverAddress: "0xd7a4F6473f32aC2Af804B3686AE8F1932bC35750",
	}
	hash, err := ensService.CreateSubdomain("first", "0xa67f7826C808d836ca7aE99d3aa183b7E6DCC3B2")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "" {
		t.Errorf("Hash should be not empty")
	}
}
