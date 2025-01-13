package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type IdentityResponse struct {
	IdentityId string `json:"IdentityId"`
}

type CredentialsResponse struct {
	Credentials struct {
		AccessKeyId  string `json:"AccessKeyId"`
		SecretKey    string `json:"SecretKey"`
		SessionToken string `json:"SessionToken"`
	} `json:"Credentials"`
}

type GraphQLResponse struct {
	Data map[string]struct {
		TotalItems int `json:"totalItems"`
	} `json:"data"`
}

func hmacSign(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func getAwsCredentials() (*CredentialsResponse, error) {
	identityBody := []byte(`{"IdentityPoolId":"us-east-2:6838b258-2477-4692-af42-54fea1c3a90a"}`)
	identityReq, err := http.NewRequest("POST", "https://cognito-identity.us-east-2.amazonaws.com", bytes.NewBuffer(identityBody))
	if err != nil {
		return nil, err
	}

	identityReq.Header.Set("Content-Type", "application/x-amz-json-1.1")
	identityReq.Header.Set("X-Amz-Target", "AWSCognitoIdentityService.GetId")

	client := &http.Client{}
	identityResp, err := client.Do(identityReq)
	if err != nil {
		return nil, err
	}
	defer identityResp.Body.Close()

	var identityResult IdentityResponse
	if err := json.NewDecoder(identityResp.Body).Decode(&identityResult); err != nil {
		return nil, err
	}

	credBody, _ := json.Marshal(map[string]string{"IdentityId": identityResult.IdentityId})
	credReq, err := http.NewRequest("POST", "https://cognito-identity.us-east-2.amazonaws.com", bytes.NewBuffer(credBody))
	if err != nil {
		return nil, err
	}

	credReq.Header.Set("Content-Type", "application/x-amz-json-1.1")
	credReq.Header.Set("X-Amz-Target", "AWSCognitoIdentityService.GetCredentialsForIdentity")

	credResp, err := client.Do(credReq)
	if err != nil {
		return nil, err
	}
	defer credResp.Body.Close()

	var credResult CredentialsResponse
	if err := json.NewDecoder(credResp.Body).Decode(&credResult); err != nil {
		return nil, err
	}

	return &credResult, nil
}

func sendRequest(body string, creds *CredentialsResponse) (*GraphQLResponse, error) {
	timestamp := time.Now().UTC()
	hmacDate := timestamp.Format("20060102")
	amzDate := timestamp.Format("20060102T150405Z")

	h := sha256.New()
	h.Write([]byte(body))
	bodyHash := hex.EncodeToString(h.Sum(nil))

	canonicalRequest := fmt.Sprintf("POST\n/graphql\n\nhost:graphql.studentportal.osu.edu\nx-amz-date:%s\nx-amz-security-token:%s\n\nhost;x-amz-date;x-amz-security-token\n%s",
		amzDate, creds.Credentials.SessionToken, bodyHash)

	h = sha256.New()
	h.Write([]byte(canonicalRequest))
	canonicalHash := hex.EncodeToString(h.Sum(nil))

	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s/us-east-2/appsync/aws4_request\n%s",
		amzDate, hmacDate, canonicalHash)

	kSecret := []byte("AWS4" + creds.Credentials.SecretKey)
	kDate := hmacSign(kSecret, hmacDate)
	kRegion := hmacSign(kDate, "us-east-2")
	kService := hmacSign(kRegion, "appsync")
	kSigning := hmacSign(kService, "aws4_request")
	signature := hex.EncodeToString(hmacSign(kSigning, stringToSign))

	req, err := http.NewRequest("POST", "https://graphql.studentportal.osu.edu/graphql", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/us-east-2/appsync/aws4_request, SignedHeaders=host;x-amz-date;x-amz-security-token, Signature=%s",
		creds.Credentials.AccessKeyId, hmacDate, signature))
	req.Header.Set("X-Amz-Security-Token", creds.Credentials.SessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func LookupUsers(usernames []string) ([]int, error) {
	creds, err := getAwsCredentials()
	if err != nil {
		return nil, err
	}

	var queryParts []string
	for i, username := range usernames {
		queryParts = append(queryParts, fmt.Sprintf(`user%d: getPerson(query: "%s") { totalItems, items { displayName } }`, i, username))
	}
	query := fmt.Sprintf(`{ %s }`, strings.Join(queryParts, "\n"))

	body, _ := json.Marshal(map[string]string{"query": query})

	result, err := sendRequest(string(body), creds)
	if err != nil {
		return nil, err
	}

	var notFoundIndexes []int
	for key, value := range result.Data {
		if value.TotalItems == 0 {
			index := 0
			fmt.Sscanf(key, "user%d", &index)
			notFoundIndexes = append(notFoundIndexes, index)
		}
	}

	return notFoundIndexes, nil
}
