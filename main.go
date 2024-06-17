package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fatih/color"
)

func main() {
	// Define flags with default values
	host := flag.String("host", "localhost", "Host (default: localhost)")
	port := flag.Int("port", 993, "Port (default: 993 for IMAP over SSL)")
	email := flag.String("email", "", "Email address (required)")
	password := flag.String("pwd", "", "Password (required)")
	folder := flag.String("folder", "", "folder to read")
	fromUID := flag.Uint("from", 0, "email uid start read from")
	flag.Parse()

	mbox, err := NewMailbox(*host, *port, *email, *password)
	if err != nil {
		log.Fatal(err)
	}

	for _, mfolder := range mbox.Folders() {
		fmt.Printf("FOLDER: %-30s\t %s \n", mfolder, color.GreenString(mbox.cleanFolderName(mfolder)))
	}

	if *folder == "" {
		fmt.Println()
		log.Fatalf("ERR: %s", color.RedString("Choose folder to read: -folder=*"))
	}

	if err := mbox.ReadAllMessages(*folder, uint32(*fromUID)); err != nil {
		color.Red("READ MSG: %s", err)
	}

}
