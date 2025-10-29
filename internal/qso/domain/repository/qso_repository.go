package repository

// QSORepository combines all QSO operations
type QSORepository interface {
	QSOReader
	QSOWriter
}
