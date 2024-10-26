package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type MailchimpClient struct {
	ApiKey string
	// for example, 'us9'
	Server string
	ListId string
}

func (m *MailchimpClient) Ping() {
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/ping", m.Server)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.SetBasicAuth("anystring", m.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}

func (m *MailchimpClient) CheckIfMemberOnList(email string) (bool, error) {
	subscriberHash := emailToSubscriberHash(email)
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s", m.Server, m.ListId, subscriberHash)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Mailchimp: check if member on list:", err)
		return false, nil
	}

	req.SetBasicAuth("anystring", m.ApiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Mailchimp: check if member on list:", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		log.Println("mailchimp: check if member on list: ", resp.Status)
		return false, fmt.Errorf("Unexpected status code %v", resp.StatusCode)
	}
}

// The subscriber hash is "The MD5 hash of the lowercase version of the list member's email address."
func emailToSubscriberHash(email string) string {
	lowerEmail := strings.ToLower(email)
	hash := md5.Sum([]byte(lowerEmail))
	subscriberHash := hex.EncodeToString(hash[:])
	return subscriberHash
}
