package tile

// Tile defines an interface for IIIF tiles
type Tile interface {
	// Height returns the height in pixels of `Tile`
	Height() int
	// Width returns the width in pixels of `Tile`
	Width() int
}
