// package update provides methods for updating EXIF data in JPEG files.
// This is a thin wrapper around code in dsoprea's go-exif and go-jpeg-image-structure packages
// and includes command-line tools for updating the EXIF data JPEG files using key-value parameters
// as well as a WebAssembly (wasm) binary for updating EXIF data in JavaScript (or other languages
// that support wasm binaries).
//
// As of this writing the majority of EXIF tags are _not_ supported. Currently only EXIF tags of types `ASCII` and `BYTE` are supported. This is not ideal but I am still trying to get familiar with the requirements of the `go-exif` package. Contributions and patches for the other remaining EXIF tag types is welcomed.
package update
