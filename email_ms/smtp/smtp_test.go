package smtp

import (
	"context"
	"testing"

	mailyak "github.com/domodwyer/mailyak/v3"
	"github.com/stretchr/testify/require"
)

func Test_Server(t *testing.T) {
	ctx := context.Background()
	sv, err := InitServerFromEnv(ctx)
	require.NoError(t, err)
	require.NotNil(t, sv)

	defer sv.Shutdown()
	go func() {
		require.NoError(t, sv.Listen())
	}()

	mail := mailyak.New("localhost:2526", nil)

	mail.From("a@localhost")
	mail.To("b@localhost")
	mail.FromName("Bananas for Friends")

	mail.Subject("Business proposition")
	mail.Plain().Set("123")
	//mail.Plain().Set("my beautiful email")
	// file, err := os.Open("server.go")
	// require.NoError(t, err)
	// mail.Attach("server.go", file)

	// mail.Send()
	require.NoError(t, mail.Send())
}
