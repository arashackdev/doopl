package apimodel

// RephraseRequest is the request for a rephrase operation.
type RephraseRequest struct {
	Text  []string `json:"text"`
	Lang  string   `json:"lang"`
	Tone  string   `json:"tone,omitempty"`
	Emoji string   `json:"emoji,omitempty"`
}

// RephraseResponse is the response from a rephrase operation.
type RephraseResponse struct {
	Results []struct {
		Text string `json:"text"`
	} `json:"results"`
}
