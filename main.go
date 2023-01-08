package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// login: login/login

type LoginForm struct {
	userEmail    string
	userPassword string
	clientID     string
}

type LoginResponse struct {
	Data LoginData `json:"data"`
}

type LoginData struct {
	Token string `json:"token"`
}

type AuthToken string

const (
	BaseURL  = "https://my.askdante.com/rest/v2/"
	ClientID = "Android"
)

func main() {
	args := os.Args[1:]

	email := os.Getenv("AD_EMAIL")
	password := os.Getenv("AD_PASSWORD")

	token := login(email, password)

	switch args[0] {
	case "start":
		start(email, token)
	case "stop":
		stop(email, token)
	default:
		fmt.Println("Usage: askdante start|stop")
	}
}

func start(email string, token AuthToken) {
	status := virtualTerminal(email, "start", token)

	if status == http.StatusOK {
		log.Println("Started tracking...")
	}
}

func stop(email string, token AuthToken) {
	status := virtualTerminal(email, "stop", token)

	if status == http.StatusOK {
		log.Println("Stopped tracking...")
	}
}

func virtualTerminal(email, op string, token AuthToken) int {
	form := url.Values{
		"userEmail":     {email},
		"userAuthToken": {string(token)},
		"clientId":      {ClientID},
	}

	resp, err := http.PostForm(BaseURL+"virtual-terminal/"+op, form)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		var data any
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(data)
	}

	return resp.StatusCode
}

func login(email, password string) AuthToken {
	form := url.Values{
		"userEmail":    {email},
		"userPassword": {password},
		"clientId":     {ClientID},
	}

	resp, err := http.PostForm(BaseURL+"login/login", form)
	if err != nil {
		log.Fatal(err)
	}

	var data LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	return AuthToken(data.Data.Token)
}
