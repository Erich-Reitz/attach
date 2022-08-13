package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func validateFlagsIfPrintSet(file, message string) error {
	if file != "" || message != "" {
		return fmt.Errorf("cannot use -p and -f or -m")
	}
	return nil 
}

func validateFlags(file, message string, shouldPrint bool) error{
	if (shouldPrint) {
		return validateFlagsIfPrintSet(file, message)
	}

	if len(file) == 0 {
		return fmt.Errorf("File is required")
	}
	if len(message) == 0 {
		return fmt.Errorf("Message is required")
	}

	return nil
}

type attachmentDetails struct {
	FilePath    string
	Message     string
	MessageTime string
}

func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}

	return true
}

func getFullPathOfUserSuppliedFile(user_filepath string) (string, error) {
	res, err := filepath.Abs(user_filepath)
	if err != nil {
		return "", err
	}

	return res, nil
}

func getUserAttachments() ([]attachmentDetails, error) {
	if !fileExists("./attachments.json") {
		return []attachmentDetails{}, nil
	}

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

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func attachMessageToFile(file, message string) error {

	if !fileExists(file) {
		return fmt.Errorf("File %s does not exist", file)
	}

	full_filepath, err := getFullPathOfUserSuppliedFile(file)
	if err != nil {
		return err
	}

	currentTime := getCurrentTime()
	info := attachmentDetails{full_filepath, message, currentTime}

	history, err := getUserAttachments()
	if err != nil {
		return err
	}

	json_data, err := mergeAttachments(history, info)
	if err != nil {
		return err
	}
	attachments_file, err := os.OpenFile("./attachments.json", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0644)
	defer attachments_file.Close()

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
	var print bool

	flag.StringVar(&message, "m", "", "Message to attach to file")
	flag.StringVar(&file, "f", "", "File to attach to message")
	flag.BoolVar(&print, "p", false, "Print the directory with messages attached to each file")
	

	flag.Parse()

	// parse flags
	err := validateFlags(file, message, print) 
	if err != nil {
		fmt.Println("Usage: attach [-m message] [-f file] [-p] print")
		flag.PrintDefaults()
		log.Fatal(err)
	}

	if print {
		history, err := getUserAttachments()
		if err != nil {
			log.Fatal(err)
		}
		for _, info := range history {
			fmt.Printf("%s: %s\n", info.FilePath, info.Message)
		}
	} else {
		err := attachMessageToFile(file, message)
		if err != nil {
			log.Fatal(err)
		}
	}



	return
}
