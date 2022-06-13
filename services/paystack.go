package services

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type PayStack struct {
	SecretKey string
}

type Data struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference        string `json:"reference"`
		AuthorizationUrl string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
	} `json:"data"`
}

func NewPaystack() *PayStack {
	secretKey := os.Getenv("PRIVATE_KEY")
	return &PayStack{
		SecretKey: secretKey,
	}
}

func (p *PayStack) SetPrivateKey() {
	secretKey := os.Getenv("PRIVATE_KEY")
	p.SecretKey = secretKey
}

func (p *PayStack) InitializePayment(info []byte) (string, error) {

	req, err := http.NewRequest(http.MethodPost, "https://api.paystack.co/transaction/initialize",
		strings.NewReader(string(info)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.SecretKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("%s", msg)
	data := Data{}
	err = json.Unmarshal(msg, &data)
	if err != nil {
		return "", err
	}

	return data.Data.AuthorizationUrl, nil
}

func (p *PayStack) VerifyReference(reference string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.paystack.co/transaction/verify/"+reference, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.SecretKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PayStack) PayStackDecodeToken(token, secret string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil

	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return claims, err
}
