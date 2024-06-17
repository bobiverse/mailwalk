package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/fatih/color"
	"github.com/logrusorgru/aurora"
)

// Mailbox ..
type Mailbox struct {
	// underlying connection fields
	*client.Client

	// almost never changes
	server string
	port   int

	// mandatory login fields
	Email    string
	Password string

	// result fields
	folders []string // imap folder names for further actions
}

// EnvelopeSeqNum - every envelope/email have seq num ID to identify email
type EnvelopeSeqNum uint32

// EnvelopeSeqNums - list of seq nums
type EnvelopeSeqNums []*EnvelopeSeqNum

// String
func (mailbox *Mailbox) String() string {
	return fmt.Sprintf("%s:%d", mailbox.server, mailbox.port)
}

// NewMailbox ..
func NewMailbox(host string, port int, email, passw string) (*Mailbox, error) {
	mailbox := &Mailbox{
		server: host,
		port:   port,

		Email:    email,
		Password: passw,
	}

	// Connect to server
	c, err := client.Dial(mailbox.String())
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] Connected", mailbox.server)
	mailbox.Client = c

	// Start a TLS session
	tlsConfig := &tls.Config{
		ServerName: mailbox.server,
		MinVersion: tls.VersionTLS11,
	}
	if err := c.StartTLS(tlsConfig); err != nil {
		return nil, err
	}
	log.Printf("[%s] TLS started", mailbox.server)

	// Login
	if err := mailbox.Login(mailbox.Email, mailbox.Password); err != nil {
		return nil, err
	}
	log.Printf("[%s] Logged in as %s", mailbox.server, mailbox.Email)

	return mailbox, nil
}

// Folders - fetch folder names where messages can be found
func (mailbox *Mailbox) Folders() []string {
	if len(mailbox.folders) > 0 {
		return mailbox.folders
	}

	// PrintInfo("Fetch mailboxes/folders..")
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- mailbox.List("", "*", mailboxes)
	}()

	mailbox.folders = []string{}
	for m := range mailboxes {
		mailbox.folders = append(mailbox.folders, m.Name)
	}

	if err := <-done; err != nil {
		return nil
	}

	return mailbox.folders
}

func (mailbox *Mailbox) cleanFolderName(folderName string) string {
	return strings.TrimPrefix(folderName, "INBOX.")
}

// Unreads - seqnums of emails not seen
func (mailbox *Mailbox) Unreads(onlyFolders, ignoreFolders []string) EnvelopeSeqNums {
	var unreads EnvelopeSeqNums

	folders := mailbox.Folders()

	// only include these
	if len(onlyFolders) > 0 {
		var filterFolders []string
		for _, folder := range folders {
			for _, onf := range onlyFolders {
				if strings.Contains(mailbox.cleanFolderName(folder), onf) {
					filterFolders = append(filterFolders, folder)
				}
			}
		}
		folders = filterFolders
	}

	// Remove ignored folders
	if len(ignoreFolders) > 0 {
		var filterFolders []string
		for _, folder := range folders {
			for _, ign := range ignoreFolders {
				if !strings.Contains(mailbox.cleanFolderName(folder), ign) {
					filterFolders = append(filterFolders, folder)
				}
			}
		}
		folders = filterFolders
	}

	for _, folder := range folders {
		// log.Printf("==> [%s]", folder)
		// chose specific folder
		_, err := mailbox.Select(folder, false)
		if err != nil {
			continue
		}

		// search for unseen messages
		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}

		uids, err := mailbox.Search(criteria)
		if err != nil {
			continue
		}

		// convert to native type
		for _, uid := range uids {
			seqnum := EnvelopeSeqNum(uid)
			unreads = append(unreads, &seqnum)
		}
	}

	return unreads
}

// ReadAllMessages reads all messages from a specified folder, starting from a given UID
func (mailbox *Mailbox) ReadAllMessages(folderName string, startUID uint32) error {
	_, err := mailbox.Select(folderName, false)
	if err != nil {
		log.Printf("Failed to select folder %s: %v", folderName, err)
		return err
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Deleted"}

	// Add date criterion for today
	today := time.Now().Format("02-Jan-2006")
	criteria.Since, _ = time.Parse("02-Jan-2006", today)

	if startUID > 0 {
		criteria.Uid = new(imap.SeqSet)
		criteria.Uid.AddRange(startUID, 0) // Start from startUID to the newest message
	}
	uids, err := mailbox.UidSearch(criteria)
	if err != nil {
		log.Printf("Failed to search messages: %v", err)
		return err
	}
	if len(uids) == 0 {
		fmt.Println("No new messages to read")
		return nil // Return the startUID if no new messages
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	// Inside the ReadAllMessages function, adjust your fetch call:
	go func() {
		section := &imap.BodySectionName{Peek: true} // Use Peek to avoid marking as read
		items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope, imap.FetchUid}
		done <- mailbox.UidFetch(seqSet, items, messages)
	}()

	for msg := range messages {

		// Extract email addresses and collect them in a slice
		var froms []string
		for _, addr := range msg.Envelope.From {
			froms = append(froms, addr.Address())
		}

		isUnread := true
		if msg.Flags != nil {
			for _, flag := range msg.Flags {
				if flag == "\\Seen" {
					isUnread = false
					break
				}
			}
		}

		fmt.Printf("[%d][%s] ", aurora.Cyan(msg.Uid), msg.Envelope.Date.Format(time.DateTime))

		if isUnread {
			fmt.Printf("*")
		}
		fmt.Printf("%-30s\t `%s`\n", aurora.Blue(strings.Join(froms, ";")), color.YellowString(msg.Envelope.Subject))

		fmt.Println(strings.Repeat("-", 80))
		time.Sleep(1 * time.Second) // Delay between processing each message
	}

	if err = <-done; err != nil {
		log.Printf("Failed to fetch messages: %v", err)
		return err
	}

	return nil // Return the last UID read
}
