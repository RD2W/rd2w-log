package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
)

// ADIFHandler defines ADIF import/export operations
type ADIFHandler interface {
	// ImportADIF imports QSOs from ADIF format
	ImportADIF(ctx context.Context, userID string, data []byte) (*model.ADIFImportResult, error)

	// ExportADIF exports QSOs to ADIF format
	ExportADIF(ctx context.Context, userID string, opts model.ADIFExportOptions) ([]byte, error)
}
