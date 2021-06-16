package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type user struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Pass        string `json:"pass"`
	Description string `json:"description"`
}

type message struct {
	Body string `json:"body"`
}

var client *http.Client = &http.Client{}

func main() {
	for {
		var rStdin *bufio.Reader = bufio.NewReader(os.Stdin)

		fmt.Println("0. Exit")
		fmt.Println("1. Add user")
		fmt.Println("2. Show users")
		fmt.Println("3. Search user")

		reply, _, err := rStdin.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		switch string(reply) {
		case "0":

			fmt.Println("E X I T I N G . . .")
			os.Exit(0)

		case "1":

			var user user

			fmt.Println("Enter a name")
			reply, _, err := rStdin.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			user.Name = string(reply)

			fmt.Println("Enter a username")
			reply, _, err = rStdin.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			user.Username = string(reply)

			fmt.Println("Enter a password")
			reply, _, err = rStdin.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			user.Pass = string(reply)

			fmt.Println("Enter a description")
			reply, _, err = rStdin.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			user.Description = string(reply)

			dataJSON, err := json.Marshal(user)
			if err != nil {
				log.Fatal(err)
			}

			request, err := http.NewRequest("POST", "http://localhost:8080/add", bytes.NewBuffer(dataJSON))
			if err != nil {
				log.Fatal(err)
			}

			request.Header.Set("Accept", "application/json")

			response, err := client.Do(request)
			if response.StatusCode != 200 {
				log.Fatal(response.Status)
			}
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			var data message
			err = json.NewDecoder(response.Body).Decode(&data)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(data.Body)
			fmt.Println()

		case "2":

			var users []user

			request, err := http.NewRequest("GET", "http://localhost:8080/show", nil)
			if err != nil {
				log.Fatal(err)
			}

			request.Header.Set("Accept", "application/json")

			response, err := client.Do(request)
			if response.StatusCode != 200 {
				log.Fatal(response.Status)
			}
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			err = json.NewDecoder(response.Body).Decode(&users)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("|%-7s|%-15s|%-15s|%-15s|%-15s|\n", "id", "Name", "Username", "Pass", "Desc")
			fmt.Println("_______________________________________________________________________")
			for i := 0; i < len(users); i++ {
				fmt.Printf("|%-7d|%-15s|%-15s|%-15s|%-15s|\n", users[i].Id, users[i].Name, users[i].Username, users[i].Pass, users[i].Description)
			}
			fmt.Println()

		case "3":

			var userSearch user

			fmt.Println("Enter id to search")
			reply, _, err := rStdin.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			userSearch.Id, _ = strconv.Atoi(string(reply))

			dataJSON, err := json.Marshal(userSearch)
			if err != nil {
				log.Fatal(err)
			}

			request, err := http.NewRequest("GET", "http://localhost:8080/search", bytes.NewBuffer(dataJSON))
			if err != nil {
				log.Fatal(err)
			}

			request.Header.Set("Accept", "application/json")

			response, err := client.Do(request)
			if response.StatusCode != 200 {
				log.Fatal(response.Status)
			}
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			if strings.Contains(string(body), "id") {

				err = json.Unmarshal(body, &userSearch)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("|%-7s|%-15s|%-15s|%-15s|%-15s|\n", "id", "Name", "Username", "Pass", "Desc")
				fmt.Println("_______________________________________________________________________")
				fmt.Printf("|%-7d|%-15s|%-15s|%-15s|%-15s|\n", userSearch.Id, userSearch.Name, userSearch.Username, userSearch.Pass, userSearch.Description)
				fmt.Println()

			} else {

				var mess message
				err = json.Unmarshal(body, &mess)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(mess.Body)
				fmt.Println()

			}
		}
	}
}
