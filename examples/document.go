package main

import (
	"context"
	"fmt"
	"os"
)

// Example: Document translation (upload, check status, download).
//
// Run with:
//
//	go run examples/document.go
//
// Environment:
//
//	DEEPL_AUTH_KEY: Your DeepL API key
//
// This example demonstrates:
// - Uploading a document for translation
// - Polling translation status
// - Downloading the translated document
// - Cleaning up resources
func exampleDocument() {
	_ = os.Getenv("DEEPL_AUTH_KEY")
	_ = context.Background()

	// Note: This is a reference implementation. To run with a real document,
	// you would open an actual file:
	//
	//   file, err := os.Open("document.txt")
	//   if err != nil {
	//     log.Fatal(err)
	//   }
	//   defer file.Close()

	fmt.Println("📄 Document Translation Example")
	fmt.Println("================================")

	fmt.Println("⚠️  This is a reference implementation.")
	fmt.Println("To use with a real document:")
	fmt.Println("  1. Open your document file")
	fmt.Println("  2. Call client.DocumentUpload(ctx, file, filename, \"DE\")")
	fmt.Println("  3. Poll with client.DocumentStatus(ctx, id, key)")
	fmt.Println("  4. Download with client.DocumentDownload(ctx, id, key, output)")
	fmt.Println("Supported formats: PDF, DOCX, PPTX, XLSX, TXT, HTML, HTM, JPG, JPEG, PNG")

	fmt.Println("Example flow:")
	fmt.Println("// Upload a document")
	fmt.Println("handle, err := client.DocumentUpload(ctx, file, \"report.pdf\", \"ES\")")
	fmt.Println("")
	fmt.Println("// Poll for completion")
	fmt.Println("for {")
	fmt.Println("  status, err := client.DocumentStatus(ctx, handle.DocumentID, handle.DocumentKey)")
	fmt.Println("  if err != nil {")
	fmt.Println("    log.Fatal(err)")
	fmt.Println("  }")
	fmt.Println("  if status.Done {")
	fmt.Println("    break")
	fmt.Println("  }")
	fmt.Println("  time.Sleep(2 * time.Second)")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("// Download result")
	fmt.Println("output, err := os.Create(\"report_es.pdf\")")
	fmt.Println("defer output.Close()")
	fmt.Println("err = client.DocumentDownload(ctx, handle.DocumentID, handle.DocumentKey, output)")
}
