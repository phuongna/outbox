package outbox

import (
	"database/sql"
	"github.com/phuongna/outbox/internal/time"
	"github.com/phuongna/outbox/internal/uuid"
)

// Publisher encapsulates the save functionality of the outbox pattern
type Publisher struct {
	store Store
	time  time.Provider
	uuid  uuid.Provider
}

// NewPublisher is the Publisher constructor
func NewPublisher(store Store) Publisher {
	return Publisher{store: store, time: time.NewTimeProvider(), uuid: uuid.NewUUIDProvider()}
}

// MessageHeader is the MessageHeader of the Message to be sent. It is used by Brokers
type MessageHeader struct {
	Key   string
	Value string
}

// Message encapsulates the contents of the message to be sent
type Message struct {
	Key     string
	Headers []MessageHeader
	Body    []byte
	Topic   string
}

// Send stores the provided Message within the provided sql.Tx
func (o Publisher) Send(msg Message, tx *sql.Tx) error {
	newID := o.uuid.NewUUID()
	record := Record{
		ID:          newID,
		Message:     msg,
		State:       PendingDelivery,
		CreatedOn:   o.time.Now().UTC(),
		LockID:      nil,
		LockedOn:    nil,
		ProcessedOn: nil,
	}

	return o.store.AddRecordTx(record, tx)
}
