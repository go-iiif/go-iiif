package tiffstructure

import (
	"os"
	"path"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath = ""
)

// getModuleRootPath returns our source-path when running from source during
// tests.
func getModuleRootPath() string {
	moduleRootPath := os.Getenv("TIFF_MODULE_ROOT_PATH")
	if moduleRootPath != "" {
		return moduleRootPath
	}

	currentWd, err := os.Getwd()
	log.PanicIf(err)

	currentPath := currentWd
	visited := make([]string, 0)

	for {
		tryStampFilepath := path.Join(currentPath, ".MODULE_ROOT")

		_, err := os.Stat(tryStampFilepath)
		if err != nil && os.IsNotExist(err) != true {
			log.Panic(err)
		} else if err == nil {
			break
		}

		visited = append(visited, tryStampFilepath)

		currentPath = path.Dir(currentPath)
		if currentPath == "/" {
			log.Panicf("could not find module-root: %v", visited)
		}
	}

	return currentPath
}

// getTestAssetsPath returns the test-asset path.
func getTestAssetsPath() string {
	if assetsPath == "" {
		moduleRootPath := getModuleRootPath()
		assetsPath = path.Join(moduleRootPath, "assets")
	}

	return assetsPath
}

// getTestExifImageFilepath returns the file-path of the default EXIF image
// asset.
func getTestExifImageFilepath() string {
	assetsPath := getTestAssetsPath()
	return path.Join(assetsPath, "file_example_TIFF_1MB-scaled.tiff")
}
