package model

// WriteTone represents the desired tone for rephrasing.
type WriteTone string

const (
	// WriteToneFormal rephrases for formal, professional language.
	WriteToneFormal WriteTone = "formal"
	// WriteToneInformal rephrases for casual, conversational language.
	WriteToneInformal WriteTone = "informal"
)

// WriteEmojiMode controls emoji usage in rephrased text.
type WriteEmojiMode string

const (
	// WriteEmojiKeep preserves existing emojis.
	WriteEmojiKeep WriteEmojiMode = "keep"
	// WriteEmojiAdd adds emojis to the rephrased text.
	WriteEmojiAdd WriteEmojiMode = "add"
	// WriteEmojiRemove removes emojis from the rephrased text.
	WriteEmojiRemove WriteEmojiMode = "remove"
)

// WriteResult represents text that has been rephrased by the Write API.
type WriteResult struct {
	// Text is the rephrased text.
	Text string
}
