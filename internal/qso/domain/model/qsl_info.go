package model

import (
	"time"

	"github.com/google/uuid"
)

type QSLInfo struct {
	ID           uuid.UUID `json:"id" db:"id"`
	QSOID        uuid.UUID `json:"qso_id" db:"qso_id"`
	Sent         bool      `json:"sent" db:"sent"`
	Received     bool      `json:"received" db:"received"`
	SentVia      string    `json:"sent_via" db:"sent_via"`
	ReceivedVia  string    `json:"received_via" db:"received_via"`
	SentDate     time.Time `json:"sent_date" db:"sent_date"`
	ReceivedDate time.Time `json:"received_date" db:"received_date"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewQSLInfo(qsoID uuid.UUID) *QSLInfo {
	return &QSLInfo{
		ID:        uuid.New(),
		QSOID:     qsoID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (q *QSLInfo) MarkSent(via string, sentDate time.Time) {
	q.Sent = true
	q.SentVia = via
	q.SentDate = sentDate
	q.UpdatedAt = time.Now()
}

func (q *QSLInfo) MarkReceived(via string, receivedDate time.Time) {
	q.Received = true
	q.ReceivedVia = via
	q.ReceivedDate = receivedDate
	q.UpdatedAt = time.Now()
}

func (q *QSLInfo) Validate() error {
	if q.Sent && q.SentVia == "" {
		return ErrInvalidQSLSentVia
	}

	if q.Received && q.ReceivedVia == "" {
		return ErrInvalidQSLReceivedVia
	}

	return nil
}
