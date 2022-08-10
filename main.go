package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func validateFlags(file, message string) error {
	if len(file) == 0 {
		return fmt.Errorf("File is required")
	}
	if len(message) == 0 {
		return fmt.Errorf("Message is required")
	}

	return nil
}

type attachmentDetails struct {
	FilePath string
	Message  string
}

func attachMessageToFile(file, message string) error {
	// check if file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}

	// get full path of file
	filePath, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	info := attachmentDetails{filePath, message}

	attachments_file, err := os.OpenFile("./attachments.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer attachments_file.Close()
	// read data from file
	content, err := ioutil.ReadFile("./attachments.json")
	if err != nil {
		return err
	}
	var history []attachmentDetails
	if len(content) != 0 {
		err = json.Unmarshal(content, &history)

		if err != nil {
			return err
		}
	}

	history = append(history, info)
	fmt.Println(history)
	// write to file
	json_data, err := json.Marshal(history)
	fmt.Println(string(json_data))
	if err != nil {
		return err
	}

	_, err = attachments_file.Write(json_data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// initialize flags
	var message string
	var file string

	flag.StringVar(&message, "m", "", "Message to attach to file")
	flag.StringVar(&file, "f", "", "File to attach to message")

	flag.Parse()

	// parse flags
	err := validateFlags(file, message)
	if err != nil {
		fmt.Println("Usage: attach [-m message] [-f file]")
		flag.PrintDefaults()
		log.Fatal(err)
	}

	// attach message to file
	err = attachMessageToFile(file, message)
	if err != nil {
		log.Fatal(err)
	}

	return
}
