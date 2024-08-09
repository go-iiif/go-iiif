package funcs

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// Encodes 'uri' as a QR code and returns it as base-64-encoded string.
func QRCodeB64(uri string) string {

	qr, err := qrcode.New(uri, qrcode.Medium)

	if err != nil {
		// log.Printf("Failed to encode '%s' as QR code, %w", uri, err)
		return "#"
	}

	qr.DisableBorder = true

	qr_png, err := qr.PNG(200)

	if err != nil {
		// log.Printf("Failed to encode '%s' as QR code, %w", uri, err)
		return "#"
	}

	return base64.StdEncoding.EncodeToString(qr_png)
}

// Encodes 'uri' as a QR code and returns it as base-64-encoded data URI.
func QRCodeDataURI(uri string) string {

	qr_b64 := QRCodeB64(uri)
	return fmt.Sprintf("data:%s;base64,%s", "image/png", qr_b64)
}
