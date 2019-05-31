package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Container struct {
	Id string `json:"vmid"`
	Name string `json:"name"`
	Email string
	Service int
}

type Containers struct {
	Data []Container `json:"data"`
}

var pve PVE


func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	config := LoadConfig()

	pve = pve.Login(config.Address, config.Username, config.Password)
	if pve.Server == "" {
		panic("Can't authenticate")
	}

	containers := getContainers()
	migrate(containers)
	//sendMails(config, containers)


}

func migrate(containers Containers) {

	for key, value := range containers.Data {

		vps := "vps3"

		if (key % 2) == 0 {
			vps = "vps4"
		}

		fmt.Println("Migrating " + value.Id + " on " + vps)

		params := url.Values{}

		params.Add("target", vps)
		params.Add("restart", "1")

		resp, err := pve.Query("POST", "/nodes/vps5/lxc/" + value.Id + "/migrate", params)

		fmt.Println(string(resp))

		fmt.Println()
		if err != nil {
			fmt.Println("Error while migrating " + value.Id)
		}

	}

}

func sendMails(config Config, containers Containers) {

	d := gomail.NewDialer(config.Email.Server.Address, config.Email.Server.Port, config.Email.Server.Username, config.Email.Server.Password)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	for _, container := range containers.Data {
		fmt.Println("Sending mail for: " + container.Email)

		m := gomail.NewMessage()

		m.SetHeader("From", config.Email.From)
		m.SetHeader("To", container.Email)
		m.SetHeader("Subject", config.Email.Subject)

		message, _ := ioutil.ReadFile(config.Email.Template)
		m.SetBody("text/html", string(message))

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", container.Email, err)
		}
		m.Reset()

		time.Sleep(config.Email.Server.Throttle * time.Second)
	}
}

func getContainers() (containers Containers){

	request, _ := pve.Query("GET", "/nodes/vps5/lxc", url.Values{})
	_ = json.Unmarshal([]byte(request), &containers)
	for key, value := range containers.Data {
		request, _ = pve.Query("GET", "/nodes/vps5/lxc/" + value.Id + "/config", url.Values{})
		var data map[string]interface{}

		json.Unmarshal([]byte(request), &data)
		config := data["data"].(map[string]interface{})

		description := strings.Split(config["description"].(string), "\n")

		if !strings.HasPrefix(description[1], "Email: ") {
			panic("Can't exctract email from VM : " + value.Id)
		}

		if !strings.HasPrefix(description[2], "Service ID: ") {
			panic("Can't exctract service id from VM : " + value.Id)
		}


		containers.Data[key].Email = strings.TrimPrefix(description[1], "Email: ")
		containers.Data[key].Service, _ = strconv.Atoi(strings.TrimPrefix(description[2], "Service ID: "))
	}

	return
}