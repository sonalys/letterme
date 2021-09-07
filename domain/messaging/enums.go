package messaging

type Queue string

const (
	AccountMS Queue = "account-ms"
	QEmailMS  Queue = "email-ms"
)

type Event string

const (
	FetchEmailPublicInfo Event = "check-email"
)
