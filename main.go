package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gen2brain/beeep"
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
	emailFlag := flag.String("email", "", "the email")
	passwordFlag := flag.String("password", "", "the password")

	flag.Parse()

	email := os.Getenv("AD_EMAIL")
	password := os.Getenv("AD_PASSWORD")

	if *emailFlag != "" {
		email = *emailFlag
	}

	if *passwordFlag != "" {
		password = *passwordFlag
	}

	token := login(email, password)

	switch os.Args[1:][0] {
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
		msg := "Started tracking..."

		log.Println(msg)

		err := beeep.Notify("askDANTE", msg, "assets/information.png")
		if err != nil {
			log.Println(err)
		}
	}
}

func stop(email string, token AuthToken) {
	status := virtualTerminal(email, "stop", token)

	if status == http.StatusOK {
		msg := "Stopped tracking..."

		log.Println(msg)

		err := beeep.Notify("askDANTE", msg, "assets/information.png")
		if err != nil {
			log.Println(err)
		}
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
