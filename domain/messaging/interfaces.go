package messaging

type Channel interface {
	Ack(tag uint64, ackBeforeReceiving bool) error
	Nack(tag uint64, ackBeforeReceiving bool, requeue bool) error
}
