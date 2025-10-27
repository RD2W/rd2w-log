package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestNewQSLInfo(t *testing.T) {
	qsoID := uuid.New()
	qslInfo := model.NewQSLInfo(qsoID)

	assert.NotNil(t, qslInfo)
	assert.NotEmpty(t, qslInfo.ID)
	assert.Equal(t, qsoID, qslInfo.QSOID)
	assert.False(t, qslInfo.Sent)
	assert.False(t, qslInfo.Received)
	assert.Empty(t, qslInfo.SentVia)
	assert.Empty(t, qslInfo.ReceivedVia)
	assert.True(t, qslInfo.SentDate.IsZero())
	assert.True(t, qslInfo.ReceivedDate.IsZero())
	assert.WithinDuration(t, time.Now(), qslInfo.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), qslInfo.UpdatedAt, time.Second)
}

func TestQSLInfo_MarkSent(t *testing.T) {
	qsoID := uuid.New()
	qslInfo := model.NewQSLInfo(qsoID)

	oldUpdatedAt := qslInfo.UpdatedAt
	sentDate := time.Now()
	via := "Bureau"

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	qslInfo.MarkSent(via, sentDate)

	assert.True(t, qslInfo.Sent)
	assert.Equal(t, via, qslInfo.SentVia)
	assert.Equal(t, sentDate, qslInfo.SentDate)
	assert.True(t, qslInfo.UpdatedAt.After(oldUpdatedAt))
}

func TestQSLInfo_MarkReceived(t *testing.T) {
	qsoID := uuid.New()
	qslInfo := model.NewQSLInfo(qsoID)

	oldUpdatedAt := qslInfo.UpdatedAt
	receivedDate := time.Now()
	via := "Direct"

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	qslInfo.MarkReceived(via, receivedDate)

	assert.True(t, qslInfo.Received)
	assert.Equal(t, via, qslInfo.ReceivedVia)
	assert.Equal(t, receivedDate, qslInfo.ReceivedDate)
	assert.True(t, qslInfo.UpdatedAt.After(oldUpdatedAt))
}

func TestQSLInfo_Validate(t *testing.T) {
	t.Run("valid empty QSL", func(t *testing.T) {
		qsoID := uuid.New()
		qslInfo := model.NewQSLInfo(qsoID)
		assert.NoError(t, qslInfo.Validate())
	})

	t.Run("valid sent QSL", func(t *testing.T) {
		qsoID := uuid.New()
		qslInfo := model.NewQSLInfo(qsoID)
		qslInfo.MarkSent("Bureau", time.Now())
		assert.NoError(t, qslInfo.Validate())
	})

	t.Run("invalid sent without via", func(t *testing.T) {
		qsoID := uuid.New()
		qslInfo := model.NewQSLInfo(qsoID)
		qslInfo.Sent = true
		qslInfo.SentVia = ""
		assert.Error(t, qslInfo.Validate())
		assert.Contains(t, qslInfo.Validate().Error(), "invalid QSL sent via")
	})

	t.Run("invalid received without via", func(t *testing.T) {
		qsoID := uuid.New()
		qslInfo := model.NewQSLInfo(qsoID)
		qslInfo.Received = true
		qslInfo.ReceivedVia = ""
		assert.Error(t, qslInfo.Validate())
		assert.Contains(t, qslInfo.Validate().Error(), "invalid QSL received via")
	})
}
