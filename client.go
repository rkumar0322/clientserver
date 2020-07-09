package main

import (
	"bufio"
	"clientserver/chat"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
)

var filePath string = "/Users/rahulkumar/go/src/clientserver"

func userInput(input string) string {
	var response string
	fmt.Print(input)
	fmt.Scanln(&response)
	return response
}

func validateUser(user string) bool {
	file1, err := os.Open(filePath + "/" + user + ".txt")
	if err != nil {
		return false
	}
	// create a scanner to read from file and split text based on lines
	scanner := bufio.NewScanner(file1)
	scanner.Split(bufio.ScanLines)
	// use Scan to iterate through the file
	for scanner.Scan() {
		if scanner.Text() == user {
			return true
		}
	}
	return false
}

func validateUserServ(c chat.ChatServiceClient, user string) bool {
	response, err := c.SayHello(context.Background(), &chat.Message{Body: "userExists," + user})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	return response.Body == "true"
}

func writeUser(user string) {
	file1, err := os.Open(filePath + "/" + user + ".txt")
	if err != nil {
		file1, err = os.Create(filePath + "/" + user + ".txt")
		if err != nil {
			panic("File Not Found")
		}
	}
	datewriter := bufio.NewWriter(file1)
	stringToWrite := user
	datewriter.WriteString(stringToWrite + "\n")
	datewriter.Flush()
	file1.Close()
}

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.128.183.117:10000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close() // this will execute last

	c := chat.NewChatServiceClient(conn)
	username := userInput("UserName: ")
	var firstPass bool
	if validateUser(username) && validateUserServ(c, username) {
		firstPass = false
	} else {
		firstPass = true
	}
	for {
		var str1 string
		var err1 error
		if firstPass == true {
			str1 = "adduser" + "," + username
			writeUser(username)
			firstPass = false
		} else {
			str1, err1 = generateInput(username)
			if err1 != nil {
				continue
			} else if str1 == "exit" {
				break
			}
		}
		response, err := c.SayHello(context.Background(), &chat.Message{Body: str1})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		log.Printf("Response from server: %s", response.Body) // display the response from the server
	}
}

func generateInput(username string) (string, error) {
	initial := userInput("What Do you want to do?: ")
	if initial == "adduser" {
		return initial + "," + userInput("Person To Add: "), nil
	} else if initial == "addchannel" {
		return initial + "," + username + "," + userInput("Channel Name: ") + "," + userInput("Public? [yes/no]: "), nil
	} else if initial == "removechannel" {
		return initial + "," + username + "," + userInput("Channel Name: "), nil
	} else if initial == "addusertochannel" {
		return initial + "," + username + "," + userInput("Person To Add: ") + "," + userInput("Channel Name: "), nil
	} else if initial == "removeuserfromchannel" {
		return initial + "," + username + "," + userInput("Person To remove: ") + "," + userInput("Channel Name: "), nil
	} else if initial == "banuserfromchannel" {
		return initial + "," + username + "," + userInput("Person To Ban: ") + "," + userInput("Channel Name: "), nil
	} else if initial == "removebanuser" {
		return initial + "," + username + "," + userInput("Person To Unban: ") + "," + userInput("Channel Name: "), nil
	} else if initial == "joinchannel" {
		return initial + "," + username + "," + userInput("Channel Name: "), nil
	} else if initial == "leavechannel" {
		return initial + "," + username + "," + userInput("Channel Name: "), nil
	} else if initial == "sendMessage" {
		return initial + "," + username + "," + userInput("Channel Name: ") + "," + userInput("Message: "), nil
	} else if initial == "showWorkspace" {
		return initial + "," + username, nil
	} else if initial == "showChannel" {
		return initial + "," + username + "," + userInput("What channel do you want to see?"), nil
	} else if initial == "exit" {
		return "exit", nil
	} else {
		return "", errors.New("Incorrect")
	}
}
