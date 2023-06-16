package pinata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type AddHeaderTransport struct {
	T           http.RoundTripper
	AccessToken string
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", adt.AccessToken))
	req.Close = true
	return adt.T.RoundTrip(req)
}

type PinanaAPI struct {
	client *http.Client
}

type hashResponse struct {
	IpfsHash  string    `json:"IpfsHash"`
	PinSize   int       `json:"PinSize"`
	Timestamp time.Time `json:"Timestamp"`
}

func (p PinanaAPI) PinImage(data []byte, imageName string) (cid string, err error) {
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)

	//File path
	fw1, _ := bw.CreateFormFile("file", imageName)
	fw1.Write(data)

	//Pinata options
	opts, _ := json.Marshal(map[string]int{"cidVersion": 1})
	p1w, _ := bw.CreateFormField("pinataOptions")
	p1w.Write(opts)
	//pinataMetadata
	meta, _ := json.Marshal(map[string]interface{}{"name": imageName, "keyvalues": map[string]string{"company": "test"}})
	p2w, _ := bw.CreateFormField("pinataMetadata")
	p2w.Write(meta)
	bw.Close()

	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", bw.FormDataContentType())

	resp, err := p.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Bad status. StatusCode = %v Data %s", resp.StatusCode, string(body))
		return
	}
	hashResp := hashResponse{}
	err = json.Unmarshal(body, &hashResp)
	cid = hashResp.IpfsHash
	return
}
func New(token string) *PinanaAPI {
	tripper := AddHeaderTransport{T: http.DefaultTransport, AccessToken: token}
	client := &http.Client{
		Transport: &tripper,
	}
	return &PinanaAPI{client: client}
}
