package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type createTokenResponse struct {
	Token string
}

type user struct {
	Email, Password string
}

type connection struct {
	Server, Port string
}

func main() {

	if len(os.Args) == 1 {
		// Make a default user and a default connection
		usr := &user{}
		usr.Email = "test@test.com"
		usr.Password = "abcdef"
		data := map[string]string{"email": usr.Email, "password": usr.Password}
		buff, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}

		conn := &connection{}
		conn.Server = "localhost"
		conn.Port = "8180"
        str := []string{"http://", conn.Server, ":", conn.Port, "/users"}
		url := strings.Join(str, "")
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(buff))

		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		// Make a token creation POST request
		tokenUrl := strings.Replace(url, "users", "tokens", -1)
		tokenResp, err := http.Post(tokenUrl, "application/json", bytes.NewBuffer(buff))

		if err != nil {
			log.Fatal(err)
		}

		token := &createTokenResponse{}
		json.NewDecoder(tokenResp.Body).Decode(token)
		tokenResp.Body.Close()

		fmt.Printf("User token: %s\n", token.Token)		
		
	} else {
		// Make a user creation POST request
		usr := &user{}
		regex := regexp.MustCompile(`[a-z0-9]*\@[a-z]*\.[a-z]*`)
		switch eval := regex.MatchString(os.Args[1]); eval {
		case true:
			usr.Email = os.Args[1]
		default:
			usr.Email = "test@test.com"
		}

		switch eval := strings.EqualFold(os.Args[2], ""); eval {
		case false:
			usr.Password = os.Args[2]
		default:
			usr.Password = "abcdef"
		}

		data := map[string]string{"email": usr.Email, "password": usr.Password}
		buff, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}

		conn := &connection{}
		switch eval := strings.EqualFold(os.Args[3], ""); eval {
		case false:
			conn.Server = os.Args[3]
		default:
			conn.Server = "localhost"
		}

		switch eval := strings.EqualFold(os.Args[4], ""); eval {
		case false:
			conn.Port = os.Args[4]
		default:
			conn.Port = "8180"
		}

		str := []string{"http://", conn.Server, ":", conn.Port, "/users"}

		url := strings.Join(str, "")
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(buff))

		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		// Make a token creation POST request
		tokenUrl := strings.Replace(url, "users", "tokens", -1)
		tokenResp, err := http.Post(tokenUrl, "application/json", bytes.NewBuffer(buff))

		if err != nil {
			log.Fatal(err)
		}

		token := &createTokenResponse{}
		json.NewDecoder(tokenResp.Body).Decode(token)
		tokenResp.Body.Close()

		fmt.Printf("User token: %s\n", token.Token)

	}	

}
