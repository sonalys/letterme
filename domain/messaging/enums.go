package messaging

type Queue = string

const (
	QAccountMS Queue = "account-ms"
	QEmailMS   Queue = "email-ms"
)

type Event = string

const (
	ECheckEmail Event = "check-email"
)
