package model

// DocumentStatus represents the status of a document translation.
type DocumentStatus string

const (
	// DocumentStatusQueued means the document is queued for processing.
	DocumentStatusQueued DocumentStatus = "queued"
	// DocumentStatusTranslating means the document is currently being translated.
	DocumentStatusTranslating DocumentStatus = "translating"
	// DocumentStatusDone means the document has been successfully translated.
	DocumentStatusDone DocumentStatus = "done"
	// DocumentStatusError means an error occurred during translation.
	DocumentStatusError DocumentStatus = "error"
)

// DocumentHandle represents a reference to an uploaded or translated document.
type DocumentHandle struct {
	// DocumentID is the unique identifier for the document.
	DocumentID string
	// DocumentKey is the secret key used to retrieve the translated document.
	DocumentKey string
}

// DocumentStatusInfo represents the current status of a document translation job.
type DocumentStatusInfo struct {
	// DocumentID is the unique identifier for the document.
	DocumentID string
	// Status is the current processing status.
	Status DocumentStatus
	// SecondsRemaining is the estimated seconds until completion (only when translating).
	SecondsRemaining *int64
	// BilledCharacters is the number of characters billed (only when done or error).
	BilledCharacters *int64
	// CharacterCount is the total number of characters in the document.
	CharacterCount *int64
	// ErrorMessage is the error description (only when status is error).
	ErrorMessage *string
}

// Done returns true if the document translation is complete (success or error).
func (d *DocumentStatusInfo) Done() bool {
	return d.Status == DocumentStatusDone || d.Status == DocumentStatusError
}
