package seed

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	iiifuri "github.com/go-iiif/go-iiif-uri"
)

type TiledImage struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func TiledImageFromString(str_uri string, no_extension bool) (*TiledImage, error) {

	ctx := context.Background()
	u, err := iiifuri.NewURI(ctx, str_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive tiled_im from URI %s, %w", str_uri, err)
	}

	return TiledImageFromURI(u, no_extension)
}

func TiledImageFromURI(u iiifuri.URI, no_extension bool) (*TiledImage, error) {

	origin := u.Origin()
	target, err := u.Target(nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive target from URI, %w", err)
	}

	if no_extension {
		target = strings.TrimSuffix(target, filepath.Ext(target))
	}

	slog.Debug("New tiled image", "origin", origin, "target", target)

	tiled_im := &TiledImage{
		Source: origin,
		Target: target,
	}

	return tiled_im, nil
}
