package common

import (
	"fmt"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"bytes"
	"encoding/json"	
)

type SplunkResponse struct {
	Text       string `json:"text"`      
	Code    int `json:"code"`   
}

func PostDataToSplunk(jsonData string, splunkServer string, splunkToken string) {

	// Define http request object
	postUri := "/services/collector"	
	postData := []byte(`{"event":`+ jsonData + `}`)
	url := splunkServer + postUri
	req, err := http.NewRequest("POST",url,bytes.NewBuffer(postData))

	// Set Auth Header for requesy
	authHeader := "Splunk " + splunkToken
	fmt.Println(authHeader)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client with no SSL check (internally signed certs at NAB)
	transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	// Post request
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
	}
	
    defer resp.Body.Close()

	// Output response
	data , _ := ioutil.ReadAll(resp.Body)

	// Convert json response to struct
	var response SplunkResponse
	err = json.Unmarshal(data, &response)

	// Output response
	fmt.Println("Response Code:", response.Code)
	fmt.Println("Response Text:", response.Text)
}