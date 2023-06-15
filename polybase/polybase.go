package polybase

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type PolybaseClient struct {
	namespace string
	url       string
}

func (c *PolybaseClient) createAuthHeader(body []byte, key string) (header string, err error) {
	//Format message to sign
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return
	}
	now := time.Now()
	timestamp := now.UnixMilli()
	messageBodyString := string(body)
	message := fmt.Sprintf("%d.%s", timestamp, messageBodyString)
	hash := crypto.Keccak256Hash([]byte(message))
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return
	}
	sigBase64 := hex.EncodeToString(sig)
	header = fmt.Sprintf("v=0,t=%d,h=eth-personal-sign,sig=0x%s", timestamp, sigBase64)
	return

}

func (c *PolybaseClient) LisRecords(collection string, key string) (map[string]interface{}, error) {
	client := &http.Client{}
	respDecoded := make(map[string]interface{})
	path := url.QueryEscape(fmt.Sprintf("%s/%s", c.namespace, collection))
	url := fmt.Sprintf("%s/v0/collections/%s/records", c.url, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return respDecoded, err
	}
	body := make(map[string]interface{})
	bodyBytes, _ := json.Marshal(body)
	header, err := c.createAuthHeader(bodyBytes, key)
	if err != nil {
		return respDecoded, err
	}
	req.Header.Add("X-Polybase-Signature", header)
	resp, err := client.Do(req)
	if err != nil {
		return respDecoded, err
	}
	err = json.NewDecoder(resp.Body).Decode(&respDecoded)
	return respDecoded, err
}

func (c *PolybaseClient) GetRecord(collection string, key string, id string) (map[string]interface{}, error) {
	client := &http.Client{}
	respDecoded := make(map[string]interface{})
	path := url.QueryEscape(fmt.Sprintf("%s/%s", c.namespace, collection))
	url := fmt.Sprintf("%s/v0/collections/%s/records/%s", c.url, path, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return respDecoded, err
	}
	body := make(map[string]interface{})
	bodyBytes, _ := json.Marshal(body)
	header, err := c.createAuthHeader(bodyBytes, key)
	if err != nil {
		return respDecoded, err
	}
	req.Header.Add("X-Polybase-Signature", header)
	resp, err := client.Do(req)
	if err != nil {
		return respDecoded, err
	}
	err = json.NewDecoder(resp.Body).Decode(&respDecoded)
	return respDecoded, err
}

func (c *PolybaseClient) CreateRecord(collection string, args []interface{}, key string) (map[string]interface{}, error) {
	client := &http.Client{}

	request, _ := json.Marshal(map[string]interface{}{
		"args": args,
	})
	respDecoded := make(map[string]interface{})

	path := url.QueryEscape(fmt.Sprintf("%s/%s", c.namespace, collection))
	url := fmt.Sprintf("%s/v0/collections/%s/records", c.url, path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	if err != nil {
		return respDecoded, err
	}
	header, err := c.createAuthHeader(request, key)
	if err != nil {
		return respDecoded, err
	}
	req.Header.Add("X-Polybase-Signature", header)
	resp, err := client.Do(req)
	if err != nil {
		return respDecoded, err
	}
	err = json.NewDecoder(resp.Body).Decode(&respDecoded)
	return respDecoded, err

}

func NewPolybaseClient(namespace, url string) (*PolybaseClient, error) {

	return &PolybaseClient{
		namespace: namespace,
		url:       url,
	}, nil
}
