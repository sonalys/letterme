package smtp

import (
	"context"
	"fmt"
	"testing"

	mailyak "github.com/domodwyer/mailyak/v3"
	"github.com/stretchr/testify/require"
)

func boilerplateMail() *mailyak.MailYak {
	mail := mailyak.New("localhost:2526", nil)
	mail.From("a@localhost")
	mail.FromName("Bananas for Friends")
	mail.To("b@localhost")
	mail.Subject("Business proposition")
	mail.Plain().Set("my beautiful email")
	return mail
}

func Test_Server(t *testing.T) {
	ctx := context.Background()
	sv, err := InitServerFromEnv(ctx)
	require.NoError(t, err)
	require.NotNil(t, sv)

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

	// TODO: test tls upgrade
	t.Run("email size is too big", func(t *testing.T) {
		// mail := boilerplateMail()
		// require.Error(t, mail.Send(), "should return error")
	})
}
