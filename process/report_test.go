package process

import (
	"fmt"
	"context"
	"testing"
	_ "embed"
)

//go:embed report_test.json
var report_test_body []byte

func TestGenerateProcessReportHTML(t *testing.T) {

	ctx := context.Background()
	html, err := GenerateProcessReportHTML(ctx, report_test_body)

	if err != nil {
		t.Fatalf("Failed to generate report HTML, %v", err)
	}

	fmt.Println(string(html))
}
