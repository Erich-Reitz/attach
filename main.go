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

	if file == "" {
		return fmt.Errorf("No file specified")
	}

	if message == "" {
		return fmt.Errorf("No message specified")
	}

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

func printFileWithAttachments(attachments []attachmentDetails, file string) error {
	path, _ := filepath.Abs(file)
	for _, attachment := range attachments {
		if attachment.FilePath == path {
			fmt.Printf("%s: %s\n", attachment.FilePath, attachment.Message)
		}
	}

	return nil
}

func printAllFilesWithAttachments(attachments []attachmentDetails) error {
	for _, attachment := range attachments {
		fmt.Printf("%s: %s\n", attachment.FilePath, attachment.Message)
	}

	return nil
}

func printAllFilesWithAttachmentsInCurrentDirectory(attachments []attachmentDetails) {
	files, err := ioutil.ReadDir(".")
	path, _ := filepath.Abs(".")

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		printFileWithAttachments(attachments, filepath.Join(path, file.Name()))
	}
}

func printAttachments(allFileAttachments bool, specificFile string) error {
	history, err := getUserAttachments()
	if err != nil {
		return err
	}

	// both set
	if allFileAttachments && specificFile != "" {
		return fmt.Errorf("Cannot specify both -a and -f")
	}

	// neither set
	if !allFileAttachments && specificFile == "" {
		printAllFilesWithAttachmentsInCurrentDirectory(history)
	}

	// only -a set
	if allFileAttachments {
		printAllFilesWithAttachments(history)
	}

	// only -f set
	if specificFile != "" {
		printFileWithAttachments(history, specificFile)
	}

	return nil
}

func main() {
	attach := flag.NewFlagSet("", flag.ExitOnError)
	attachFilename := attach.String("f", "", "File to attach message to")
	message := attach.String("m", "", "Message to attach to file")

	print := flag.NewFlagSet("-p", flag.ExitOnError)
	printAll := print.Bool("a", false, "Print all attachments")
	printFilename := print.String("f", "", "Print all attachments with file")

	if len(os.Args) < 2 {
		fmt.Println("Usage: attach (-a <attach> -f <file> && -m <message>) || (-p <print> -a <print all> | -f <file>)")
		log.Fatal("No command specified")
	}

	switch os.Args[1] {
	case "-a":
		attach.Parse(os.Args[2:])
		err := attachMessageToFile(*attachFilename, *message)
		if err != nil {
			log.Fatal(err)
		}
	case "-p":
		print.Parse(os.Args[2:])
		err := printAttachments(*printAll, *printFilename)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Usage: attach (-a <attach> -f <file> && -m <message>) || (-p <print> -a <print all> | -f <file>)")
		log.Fatal("Invalid command")
	}

	return
}
