package response

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
)

// Pagination is a struct containing pagination metrics for a given API response.
type Pagination struct {
	// The current page of results for an API request.
	Page int `json:"page"`
	// The total number of pages of results for an API request.
	Pages int `json:"pages"`
	// The number of results, per page, for an API request.
	PerPage int `json:"perpage"`
	// The total number of results, across all pages, for an API request.
	Total int `json:"total"`
}

// Given an API response try to derive pagination metrics.
func DerivePagination(ctx context.Context, fh io.ReadSeekCloser) (*Pagination, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	page_rsp := gjson.GetBytes(body, "*.page")

	if !page_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (page) in response")
	}

	pages_rsp := gjson.GetBytes(body, "*.pages")

	if !pages_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (pages) in response")
	}

	perpage_rsp := gjson.GetBytes(body, "*.perpage")

	if !perpage_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (perpage) in response")
	}

	total_rsp := gjson.GetBytes(body, "*.total")

	if !total_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (total) in response")
	}

	pg := &Pagination{
		Page:    int(page_rsp.Int()),
		Pages:   int(pages_rsp.Int()),
		PerPage: int(perpage_rsp.Int()),
		Total:   int(total_rsp.Int()),
	}

	return pg, nil
}
