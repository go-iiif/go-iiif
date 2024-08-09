package update

import (
	"fmt"
)

func AppendGPSPropertiesWithLatitudeAndLongitude(exif_props map[string]interface{}, lat float64, lon float64) error {

	gps_lat, err := PrepareDecimalGPSLatitudeTag(lat)

	if err != nil {
		return fmt.Errorf("Failed to prepare GPSLatitudeTag, %v", err)
	}

	gps_lat_ref, err := PrepareDecimalGPSLatitudeRefTag(lat)

	if err != nil {
		return fmt.Errorf("Failed to prepare GPSLatitudeRefTag, %v", err)
	}

	gps_lon, err := PrepareDecimalGPSLongitudeTag(lon)

	if err != nil {
		return fmt.Errorf("Failed to prepare GPSLatitudeTag, %v", err)
	}

	gps_lon_ref, err := PrepareDecimalGPSLongitudeRefTag(lon)

	if err != nil {
		return fmt.Errorf("Failed to prepare GPSLatitudeRefTag, %v", err)
	}

	exif_props["GPSLatitude"] = gps_lat
	exif_props["GPSLatitudeRef"] = gps_lat_ref
	exif_props["GPSLongitude"] = gps_lon
	exif_props["GPSLongitudeRef"] = gps_lon_ref

	return nil
}
