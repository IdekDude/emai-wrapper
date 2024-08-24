package email_wrapper

import (
	"fmt"
	"testing"
)

func TestEmail(t *testing.T) {
	client := &Client{
		Email: &Email{
			ImapEmail:    "xxxxxxx@icloud.com",
			ImapPassword: "xxxxxx",
			ImapHost:     "imap.mail.me.com:993",
		},
	}
	err := client.Connect()
	if err != nil {
		t.Error(err)
	}
	err = client.Login()
	if err != nil {
		t.Error(err)
	}
	otp, err := client.GetOTP("Amazon", "xxxxxxx@icloud.com")
	if err != nil {
		fmt.Println(err)
	}
	t.Log(otp)
}
