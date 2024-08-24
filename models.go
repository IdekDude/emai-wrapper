package email_wrapper

import "github.com/emersion/go-imap/client"

type Email struct {
	ImapEmail    string
	ImapPassword string
	ImapHost     string
}

type Client struct {
	ImapClient *client.Client
	Email	   *Email
}