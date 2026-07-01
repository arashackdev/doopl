package apimodel

// DocumentUploadResponse is the response from a document upload.
type DocumentUploadResponse struct {
	DocumentID  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

// DocumentStatusResponse is the response from checking document status.
type DocumentStatusResponse struct {
	DocumentID       string  `json:"document_id"`
	Status           string  `json:"status"`
	SecondsRemaining *int64  `json:"seconds_remaining"`
	BilledCharacters *int64  `json:"billed_characters"`
	CharacterCount   *int64  `json:"character_count"`
	ErrorMessage     *string `json:"error_message"`
}
