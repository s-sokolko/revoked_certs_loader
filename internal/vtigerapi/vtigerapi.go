package vtigerapi

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func UpdateCertificatesViaAPI(apiurl, user, key string, serialsToIds map[string]int, serialsToDates map[string]string) error {
	baseUrl, err := url.Parse(apiurl)

	if err != nil {
		return err
	}

	q := url.Values{
		"operation": {"getchallenge"},
		"username":  {user},
	}

	challengeUrl := baseUrl
	baseUrl.RawQuery = q.Encode()
	resp, err := http.Get(challengeUrl.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errorText := fmt.Sprintf("Wrong HTTP status code of %d", resp.StatusCode)
		return errors.New(errorText)
	}
	defer resp.Body.Close()

	type ChallengeResponse struct {
		Success bool `json:"success"`
		Result  struct {
			Token      string `json:"token"`
			ServerTime int    `json:"serverTime"`
			ExpireTime int    `json:"expireTime"`
		} `json:"result"`
	}
	var cr ChallengeResponse
	err = json.NewDecoder(resp.Body).Decode(&cr)
	if err != nil {
		return err
	}
	tokkey := cr.Result.Token + key
	hash := md5.Sum([]byte(tokkey))
	md5sum := hex.EncodeToString(hash[:])
	values := url.Values{
		"operation": {"login"},
		"username":  {user},
		"accessKey": {md5sum}}
	encodedValues := values.Encode()
	q = url.Values{}
	baseUrl.RawQuery = q.Encode()
	resp2, err2 := http.Post(baseUrl.String(), "application/x-www-form-urlencoded", strings.NewReader(encodedValues))
	if err2 != nil {
		return err2
	}
	defer resp2.Body.Close()

	type LoginResponse struct {
		Success bool `json:"success"`
		Result  struct {
			SessionName string `json:"sessionName"`
		} `json:"result"`
	}

	var res LoginResponse

	json.NewDecoder(resp2.Body).Decode(&res)

	sessionName := res.Result.SessionName

	for certSerial, certId := range serialsToIds {
		revDate := serialsToDates[certSerial]
		element := fmt.Sprintf(`{"id": "51x%d", "cert_revoked": "1", "cert_rev_date": "%s"}`, certId, revDate)
		values = url.Values{
			"operation":   {"revise"},
			"sessionName": {sessionName},
			"element":     {element},
		}
		encodedValues = values.Encode()
		resp3, err3 := http.Post(baseUrl.String(), "application/x-www-form-urlencoded", strings.NewReader(encodedValues))
		if err3 != nil {
			return err3
		}
		defer resp3.Body.Close()

		var res2 map[string]interface{}

		json.NewDecoder(resp3.Body).Decode(&res2)
	}

	return nil
}
