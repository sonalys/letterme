package smtp

// SMTP_TEST is an integration test for utilizing the server's interface to send mocked emails.

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	mailyak "github.com/domodwyer/mailyak/v3"
	"github.com/stretchr/testify/require"
)

func Test_SMTPServer(t *testing.T) {
	ctx := context.Background()
	sv, err := InitServerFromEnv(ctx)
	require.NoError(t, err)
	require.NotNil(t, sv)

	boilerplateMail := func() *mailyak.MailYak {
		mail, err := mailyak.NewWithTLS(sv.c.Address, nil, sv.tls)
		require.NoError(t, err)
		mail.From("a@localhost")
		mail.FromName("Bananas for Friends")
		mail.To("b@localhost")
		mail.Subject("Business proposition")
		mail.Plain().Set("my beautiful email")
		return mail
	}

	defer sv.Shutdown()
	go sv.Listen()

	t.Run("invalid from", func(t *testing.T) {
		mail := boilerplateMail()
		mail.From("bananas")
		require.Error(t, mail.Send(), "should return error")
	})

	t.Run("invalid to", func(t *testing.T) {
		mail := boilerplateMail()
		mail.To("bananas")
		require.Error(t, mail.Send(), "should return error")
	})

	t.Run("recipient outside domain", func(t *testing.T) {
		mail := boilerplateMail()
		mail.To("bananas@gmail.com")
		require.Error(t, mail.Send(), "should return error")
	})

	t.Run("second recipient outside domain", func(t *testing.T) {
		mail := boilerplateMail()
		mail.To(fmt.Sprintf("alysson@%s", sv.c.Hostname), "bananas@gmail.com")
		require.Error(t, mail.Send(), "should return error")
	})

	t.Run("email size is too big", func(t *testing.T) {
		mail := boilerplateMail()

		reader := bytes.NewReader(make([]byte, 30*MB))
		mail.AttachInline("big_attachment_30mb", reader)
		require.Error(t, mail.Send(), "should return error")
	})
}
