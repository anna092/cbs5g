package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func unsubscribe() {
	// Specify the URL you want to send the request to
	url := "http://127.0.0.18:8000/namf-comm/v1/non-ue-n2-messages/subscriptions/"
	//url := "http://192.168.56.102:8000/namf-comm/v1/non-ue-n2-messages/subscriptions/"
	id := flag.String("id", "", "N2 Notify Subscriptions id")
	flag.Parse()
	url = url + *id
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err)
		return
	}
	fmt.Println(string(body))
}
