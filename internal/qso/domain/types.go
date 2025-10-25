package domain

import "time"

// QSOLists содержит списки QSO с метаданными пагинации
type QSOLists struct {
	QSOs       []*QSO `json:"qsos"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	HasNext    bool   `json:"has_next"`
}

// ADIFImportResult содержит результат импорта ADIF файла
type ADIFImportResult struct {
	ImportedCount int      `json:"imported_count"`
	ErrorCount    int      `json:"error_count"`
	Errors        []string `json:"errors"`
}

// ADIFExportOptions содержит параметры для экспорта ADIF
type ADIFExportOptions struct {
	UserID     string    `json:"user_id"`
	DateFrom   time.Time `json:"date_from"`
	DateTo     time.Time `json:"date_to"`
	IncludeQSL bool      `json:"include_qsl"`
}

// BandModeStatistics содержит статистику по диапазонам и режимам
type BandModeStatistics struct {
	Band  string `json:"band"`
	Mode  string `json:"mode"`
	Count int    `json:"count"`
}
