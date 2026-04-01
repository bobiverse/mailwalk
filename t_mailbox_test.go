package main

import "testing"

func TestCleanFolderName(t *testing.T) {
	m := &Mailbox{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "removes inbox dot prefix", input: "INBOX.Archive", want: "Archive"},
		{name: "keeps string without prefix", input: "Sent", want: "Sent"},
		{name: "only first prefix removed", input: "INBOX.INBOX.News", want: "INBOX.News"},
		{name: "prefix only", input: "INBOX.", want: ""},
		{name: "case sensitive", input: "inbox.Archive", want: "inbox.Archive"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.cleanFolderName(tt.input)
			if got != tt.want {
				t.Fatalf("cleanFolderName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
