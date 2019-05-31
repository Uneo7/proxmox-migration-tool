package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Config struct {
	Address string		`json:"address"`
	Username string  	`json:"username"`
	Password string  	`json:"password"`
	Node string 	 	`json:"node"`
	Targets []string 	`json:"targets"`
	Email MailConfig 	`json:"mail"`
}

type MailConfig struct {
	Activated bool 		`json:"activated"`
	Template string 	`json:"template"`
	From string 		`json:"from"`
	Subject string 		`json:"subject"`
	Server MailServer 	`json:"smtp"`
}

type MailServer struct {
	Address string 			`json:"server"`
	Username string 		`json:"username"`
	Password string 		`json:"password"`
	Port int 			`json:"port"`
	Throttle time.Duration	`json:"throttle"`
}

func LoadConfig() (config Config) {
	data, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(data), &config)
	return
}