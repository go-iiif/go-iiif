package vibrant

import (
	"golang.org/x/image/draw"
	"image"
	"os"
	"testing"
)

func TestScaleImageDown(t *testing.T) {
	file, err := os.Open("test_files/1.jpg")
	if err != nil {
		t.Fatal(err)
	}
	original, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	originalSize := original.Bounds().Size()
	tests := map[uint64]image.Point{
		1000000: originalSize,
		uint64(originalSize.X * originalSize.Y): originalSize,
		60000: image.Rect(0, 0, 300, 199).Size(),
	}
	for area, expected := range tests {
		scaled := ScaleImageDown(original, area, draw.NearestNeighbor)
		scaledSize := scaled.Bounds().Size()
		if scaledSize != expected {
			t.Fatalf("Expected size %v does not match %v\n", expected, scaledSize)
		}
		if scaledSize == originalSize && scaled != original {
			t.Fatal("A new image was created when it should have used the original")
		}
	}
}
