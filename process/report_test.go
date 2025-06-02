package process

import (
	_ "context"
	_ "embed"
	_ "fmt"
	"testing"
)

// go:embed report_test.json
// var report_test_body []byte

func TestGenerateProcessReportHTML(t *testing.T) {

	t.Skip()

	/*
		ctx := context.Background()
		html, err := GenerateProcessReportHTML(ctx, report_test_body)

		if err != nil {
			t.Fatalf("Failed to generate report HTML, %v", err)
		}

		fmt.Println(string(html))
	*/
}
