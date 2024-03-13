// package client provides interfaces for common methods for accessing the Flickr API.
// Currently there is only a single client interface that calls the Flickr API using the OAuth1
// authentication and authorization scheme but it is assumed that eventually there will be at least
// one other when OAuth1 is superseded.
package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/response"
	"io"
	"net/url"
	"strconv"
	"time"
)

// The default Flickr API endpoint.
const API_ENDPOINT string = "https://api.flickr.com/services/rest"

// The default Flickr API endpoint for uploading images.
const UPLOAD_ENDPOINT string = "https://up.flickr.com/services/upload/"

// The default Flickr API endpoint for replacing images.
const REPLACE_ENDPOINT string = "https://up.flickr.com/services/replace/"

// Client is the interface that defines common methods for all Flickr API Client implementations.
// Currently there is only a single implementation that calls the Flickr API using the OAuth1
// authentication and authorization scheme but it is assumed that eventually there will be at least
// one other when OAuth1 is superseded.
type Client interface {
	// Return a new Client instance that uses the credentials included in the auth.AccessToken instance.
	WithAccessToken(context.Context, auth.AccessToken) (Client, error)
	// Call the Flickr API and create a new request token as part of the token authorization flow.
	GetRequestToken(context.Context, string) (auth.RequestToken, error)
	// Generate the URL using a request token and permissions string used to redirect a user to in order to authorize a token request.
	GetAuthorizationURL(context.Context, auth.RequestToken, string) (string, error)
	// Call the Flickr API to exchange a request and authorization token for a permanent access token.
	GetAccessToken(context.Context, auth.RequestToken, auth.AuthorizationToken) (auth.AccessToken, error)
	// Execute a Flickr API method.
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
	// Upload an image using the Flickr API.
	Upload(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
	// Replace an image using the Flickr API.
	Replace(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
}

// ExecuteMethodPaginatedCallback is the interface for callback functions passed to the
// ExecuteMethodPaginatedWithClient method.
type ExecuteMethodPaginatedCallback func(context.Context, io.ReadSeekCloser, error) error

// ExecuteMethodPaginatedWithClient invokes the Flickr API using a Client instance and then continues
// to invoke that method as many times as necessary to paginate through all of the results. Each result
// is passed to the ExecuteMethodPaginatedCallback for processing.
func ExecuteMethodPaginatedWithClient(ctx context.Context, cl Client, args *url.Values, cb ExecuteMethodPaginatedCallback) error {

	page := 1
	pages := -1

	if args.Get("page") == "" {
		args.Set("page", strconv.Itoa(page))
	} else {

		p, err := strconv.Atoi(args.Get("page"))

		if err != nil {
			return fmt.Errorf("Invalid page number '%s', %v", args.Get("page"), err)
		}

		page = p
	}

	for {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		fh, err := cl.ExecuteMethod(ctx, args)

		err = cb(ctx, fh, err)

		if err != nil {
			return err
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return fmt.Errorf("Failed to rewind response, %v", err)
		}

		if pages == -1 {

			pagination, err := response.DerivePagination(ctx, fh)

			if err != nil {
				return err
			}

			pages = pagination.Pages
		}

		page += 1

		if page <= pages {
			args.Set("page", strconv.Itoa(page))
		} else {
			break
		}
	}

	return nil
}

// UploadAsyncWithClient invokes the Flickr API using a Client instance to upload an image asynchronously
// and then waits, invoking the CheckTicketWithClient method at regular intervals, until the upload is
// complete.
func UploadAsyncWithClient(ctx context.Context, cl Client, fh io.Reader, args *url.Values) (int64, error) {

	args.Set("async", "1")

	rsp, err := cl.Upload(ctx, fh, args)

	if err != nil {
		return 0, err
	}

	return checkAsyncResponseWithClient(ctx, cl, rsp)
}

// ReplaceAsyncWithClient invokes the Flickr API using a Client instance to replace an image asynchronously
// and then waits, invoking the CheckTicketWithClient method at regular intervals, until the replacement is
// complete.
func ReplaceAsyncWithClient(ctx context.Context, cl Client, fh io.Reader, args *url.Values) (int64, error) {

	args.Set("async", "1")

	rsp, err := cl.Replace(ctx, fh, args)

	if err != nil {
		return 0, err
	}

	return checkAsyncResponseWithClient(ctx, cl, rsp)
}

func checkAsyncResponseWithClient(ctx context.Context, cl Client, rsp_fh io.ReadSeekCloser) (int64, error) {

	ticket, err := response.UnmarshalUploadTicketResponse(rsp_fh)

	if err != nil {
		return 0, err
	}

	if ticket.Error != nil {
		return 0, ticket.Error
	}

	if ticket.TicketId == "" {
		return 0, fmt.Errorf("Missing ticket ID")
	}

	return CheckTicketWithClient(ctx, cl, ticket)
}

// CheckTicketWithClient calls with Flickr API with Client and reponse.UploadTicket at regular intervals
// (every 2 seconds) to check the status of an upload ticket. If successful it will return the photo ID
// assigned to the upload.
func CheckTicketWithClient(ctx context.Context, cl Client, ticket *response.UploadTicket) (int64, error) {

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return 0, nil
		case <-ticker.C:

			args := &url.Values{}
			args.Set("method", "flickr.photos.upload.checkTickets")
			args.Set("tickets", ticket.TicketId)

			check_rsp, err := cl.ExecuteMethod(ctx, args)

			if err != nil {
				return 0, err
			}

			check_ticket, err := response.UnmarshalCheckTicketResponse(check_rsp)

			if err != nil {
				return 0, err
			}

			for _, t := range check_ticket.Uploader.Tickets {

				if t.TicketId != ticket.TicketId {
					continue
				}

				if t.Complete != 1 {
					continue
				}

				// Because the Flickr API returns strings
				// for photo IDs

				str_id := t.PhotoId

				id, err := strconv.ParseInt(str_id, 10, 64)

				if err != nil {
					return 0, err
				}

				return id, nil
			}
		}
	}

}
