package email_wrapper

import (
	"errors"
	"io/ioutil"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
)

var (
	flxTokenRe            = regexp.MustCompile(`(?m)activationToken=(.*?)&amp`)
	ueCodeRe              = regexp.MustCompile(`(?m)<p>(\b\d{4}\b)<\/p>`)
	errFailedImapConnect  = errors.New("failed to connect to IMAP server")
	errFailedSelectInbox  = errors.New("failed to select INBOX")
	errFailedGetBody      = errors.New("failed to get body")
	errFailedReadMessage  = errors.New("failed to read message")
	errFailedParseMessage = errors.New("failed to parse message")
)

func (c *Client) Connect() error {
	var err error
	c.ImapClient, err = client.DialTLS(c.Email.ImapHost, nil)
	if err != nil {
		return errFailedImapConnect
	}
	imap.CharsetReader = charset.Reader
	return nil
}

func (c *Client) Login() error {
	return c.ImapClient.Login(c.Email.ImapEmail, c.Email.ImapPassword)
}

func (c *Client) GetOTP(site string, email string) (string, error) {
	_, err := c.ImapClient.Select("INBOX", false)
	if err != nil {
		return "nil", errFailedSelectInbox
	}

	section := &imap.BodySectionName{}

	foundEmail := false
	var code string

	for !foundEmail {
		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			seqset := new(imap.SeqSet)

			m := c.ImapClient.Mailbox().Messages

			seqset.AddRange(m, m)

			done <- c.ImapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
		}()

		msg := <-messages

		switch site {
		case "Amazon":
			if msg.Envelope.Subject == "Verify your new Amazon account" && strings.EqualFold(msg.Envelope.To[0].Address(), email) {
				foundEmail = true

				r := msg.GetBody(section)
				if r == nil {
					return "nil", errFailedGetBody
				}

				m, err := mail.ReadMessage(r)
				if err != nil {
					return "nil", errFailedReadMessage
				}

				rx := quotedprintable.NewReader(m.Body)

				body, err := ioutil.ReadAll(rx)
				if err != nil {
					return "nil", errFailedParseMessage
				}

				code = strings.Split(strings.Split(string(body), "class=\"otp\">")[1], "</p>")[0]
			}
			break
		case "FLX":
			if msg.Envelope.Subject == "Confirm your FLX account." && strings.EqualFold(msg.Envelope.To[0].Address(), email) {
				foundEmail = true

				r := msg.GetBody(section)
				if r == nil {
					return "nil", errFailedGetBody
				}

				m, err := mail.ReadMessage(r)
				if err != nil {
					return "nil", errFailedReadMessage
				}

				rx := quotedprintable.NewReader(m.Body)

				body, err := ioutil.ReadAll(rx)
				if err != nil {
					return "nil", errFailedParseMessage
				}

				code = flxTokenRe.FindStringSubmatch(string(body))[1]
			}
			break
		case "Twitter":
			break
		case "Gmail":
			break
		case "Yahoo":
			break
		case "Outlook":
			break
		case "Nike":
			break
		case "UberEats":
			if strings.EqualFold(msg.Envelope.To[0].Address(), email) && strings.Contains(msg.Envelope.Subject, "Welcome to Uber") {
				foundEmail = true

				r := msg.GetBody(section)
				if r == nil {
					return "nil", errFailedGetBody
				}

				m, err := mail.ReadMessage(r)
				if err != nil {
					return "nil", errFailedReadMessage
				}

				rx := quotedprintable.NewReader(m.Body)

				body, err := ioutil.ReadAll(rx)
				if err != nil {
					return "nil", errFailedParseMessage
				}

				code = ueCodeRe.FindStringSubmatch(string(body))[1]
			}
		}
		time.Sleep(3 * time.Second)
	}

	return code, nil
}
