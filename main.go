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

func userSuppliedFileExists(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return err
	}

	return nil
}

func getFullPathOfUserSuppliedFile(user_filepath string) (string, error) {
	res, err := filepath.Abs(user_filepath)
	if err != nil {
		return "", err
	}

	return res, nil
}

func getUserAttachments(file *os.File) ([]attachmentDetails, error) {

	defer file.Close()
	// read data from file
	content, err := ioutil.ReadFile("./attachments.json")
	if err != nil {
		return nil, err
	}
	var history []attachmentDetails
	if len(content) != 0 {
		err = json.Unmarshal(content, &history)

		if err != nil {
			return nil, err
		}
	}

	return history, nil
}

func mergeAttachments(attachments []attachmentDetails, newInfo attachmentDetails) ([]byte, error) {
	attachments = append(attachments, newInfo)

	// write to file
	json_data, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}

	return json_data, nil
}

func attachMessageToFile(file, message string) error {

	if err := userSuppliedFileExists(file); err != nil {
		return err
	}

	full_filepath, err := getFullPathOfUserSuppliedFile(file)
	if err != nil {
		return err
	}

	info := attachmentDetails{full_filepath, message}
	attachments_file, err := os.OpenFile("./attachments.json", os.O_RDONLY, 0644)

	history, err := getUserAttachments(attachments_file)
	if err != nil {
		return err
	}

	json_data, err := mergeAttachments(history, info)
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
