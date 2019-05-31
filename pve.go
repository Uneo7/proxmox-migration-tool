package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type PVE struct {
	Token string
	Ticket string
	Username string
	Server string
}

func (pve PVE) Login(address string, username string, password string) PVE {

	serverUrl := "https://" + address + "/api2/json/access/ticket"
	data := "username=" + username + "&password=" + password

	var response, err = http.Post(serverUrl, "application/x-www-form-urlencoded", strings.NewReader(data))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		response, _ := ioutil.ReadAll(response.Body)

		var data map[string]interface{}
		err = json.Unmarshal([]byte(response), &data)
		if err != nil {
			fmt.Printf("JSON Unmarshal error %s\n", err)
		}

		ticket := data["data"].(map[string]interface{})

		pve.Ticket = ticket["ticket"].(string)
		pve.Token = ticket["CSRFPreventionToken"].(string)
		pve.Username = ticket["username"].(string)
		pve.Server = "https://" + address + "/api2/json"
		return pve
	}

	return PVE{}
}

func (pve PVE) Query(method string, resource string, values url.Values) (b []byte, err error) {

	client := &http.Client{}
	request, err := http.NewRequest(method, pve.Server + resource, strings.NewReader(values.Encode()))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	request.AddCookie(&http.Cookie{
		Name: "PVEAuthCookie",
		Value: pve.Ticket,
	})

	if method != "GET" {
		request.Header.Set(
			"CSRFPreventionToken",
			pve.Token,
		)
	}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	return ioutil.ReadAll(response.Body)
}