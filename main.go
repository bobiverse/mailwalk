package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/araddon/dateparse"
	"github.com/fatih/color"
	"github.com/logrusorgru/aurora"
)

func main() {
	// Define flags with default values
	host := flag.String("host", "localhost", "Host (default: localhost)")
	port := flag.Int("port", 993, "Port (default: 993 for IMAP over SSL)")
	isTls := flag.Bool("tls", true, "TLS support")
	email := flag.String("email", "", "Email address (required)")
	password := flag.String("pwd", "", "Password (required)")
	folder := flag.String("folder", "", "folder to read")
	fromUID := flag.Uint("from", 0, "email uid start read from")
	command := flag.String("cmd", "", "bash command or script to execute")
	timeout := flag.Uint("timeout", 30, "timeout (seconds) for connection to host")

	dateFrom := flag.String("datefrom", "", "Date from")
	dateTo := flag.String("dateto", "", "Date to")

	flag.Parse()

	isTLSStr := aurora.Red("-- NO --")
	if *isTls {
		isTLSStr = aurora.Magenta("Yes")
	}

	isPwdStr := aurora.Red("-- NO --")
	if *password != "" {
		isPwdStr = aurora.Magenta("Yes")
	}

	cmdStr := *command
	if len(cmdStr) > 40 {
		cmdStr = cmdStr[:38] + ".."
	}

	dTimeout := time.Duration(*timeout) * time.Second

	dFrom, _ := dateparse.ParseAny(*dateFrom)
	dTo, _ := dateparse.ParseAny(*dateTo)

	fmt.Printf("%20s: %s\n", "Host", aurora.Magenta(*host))
	fmt.Printf("%20s: %d\n", "Port", aurora.Magenta(*port))
	fmt.Printf("%20s: %d\n", "TLS", aurora.Magenta(isTLSStr))
	fmt.Printf("%20s: %s\n", "Email", aurora.Magenta(*email))
	fmt.Printf("%20s: %s\n", "Password", isPwdStr)
	fmt.Printf("%20s: %s\n", "Timeout", aurora.Magenta(dTimeout))
	fmt.Printf("%20s: %s\n", "Folder", aurora.Cyan(*folder))

	if !dFrom.IsZero() {
		fmt.Printf("%20s: %s\n", "Date From", aurora.Magenta(dFrom.Format(time.DateOnly)))
	}

	if !dTo.IsZero() {
		fmt.Printf("%20s: %s\n", "Date To", aurora.Magenta(dTo.Format(time.DateOnly)))
	}

	fmt.Printf("%20s: %s\n", "Command", aurora.Yellow(cmdStr))

	//Ping host first
	if err := pingHost(*host, *port, dTimeout); err != nil {
		log.Fatal(err)
	}
	log.Printf(">> Ping host: %s", aurora.Green("OK"))

	// Mailbox
	mbox, err := NewMailbox(*host, *port, *isTls, *email, *password, dTimeout)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context that is canceled when the user presses CTRL+C
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived an interrupt, stopping...")
		cancel()
	}()
	mbox.context = ctx

	for _, mfolder := range mbox.Folders() {
		fmt.Printf("FOLDER: %-30s\t %s \n", mfolder, color.GreenString(mbox.cleanFolderName(mfolder)))
	}

	if *folder == "" {
		fmt.Println()
		log.Fatalf("ERR: %s", color.RedString("Choose folder to read: -folder=*"))
	}

	if err := mbox.ReadAllMessages(*folder, uint32(*fromUID), dFrom, dTo, *command); err != nil {
		color.Red("READ MSG: %s", err)
	}

}
