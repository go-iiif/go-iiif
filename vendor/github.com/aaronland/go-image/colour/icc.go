package colour

import (
	"fmt"
	"io"

	"github.com/mandykoh/prism/meta/autometa"
)

// ICC_DISPLAY_P3 is the ICC profile description for the Apple Display P3 profile.
const ICC_DISPLAY_P3 string = "Display P3"

// ICC_EPSON_RGB_G18 is the ICC profile description for the EPSON  Standard RGB - Gamma 1.8 profile.
const ICC_EPSON_RGB_G18 string = "EPSON  Standard RGB - Gamma 1.8"

// ICC_ADOBE_RGB_1998 is the ICC profile description for the Adobe RGB (1998) profile.
const ICC_ADOBE_RGB_1998 string = "Adobe RGB (1998)"

// ICC_SRGB_21 is the ICC profile description for the sRGB IEC61966-2.1 profile.
const ICC_SRGB_21 string = "sRGB IEC61966-2.1"

// ICC_GENERIC_GRAY is the ICC profile description for the Generic Gray Profile profile.
const ICC_GENERIC_GRAY string = "Generic Gray Profile"

// ICC_CAMERA_RGB is the ICC profile description for the Camera RGB Profile profile.
const ICC_CAMERA_RGB string = "Camera RGB Profile"

// ICCProfileDescription attempts to the derive the ICC profile description from the body of 'r'
func ICCProfileDescription(r io.Reader) (string, error) {

	md, _, err := autometa.Load(r)

	if err != nil {
		return "", fmt.Errorf("Failed to load metadata, %w", err)
	}

	pr, err := md.ICCProfile()

	if err != nil {
		return "", fmt.Errorf("Failed to derive ICC profile, %w", err)
	}

	if pr == nil {
		return "", fmt.Errorf("Missing profile")
	}

	return pr.Description()
}

/*

TBD... basically colour profiles are a giant bag of whatEVAR aren't they?

2/11, Curves
Adobe RGB
Adobe RGB (1998)
Apple Wide Color Sharing Profile
BH_DRUM_RGB_3/5/02
Camera RGB Profile
Color LCD
Custom RGB
DELL U2719DC
Display
Display P3
Dot Gain 20%
EPS MONITOR SETTINGS 10/8/98
EPSON  Gray - Gamma 1.8
EPSON  Gray - Gamma 2.2
EPSON  Standard RGB - Gamma 1.8
EPSON  sRGB
Ekta Space PS 5, J. Holmes.icm
Generic Gray Gamma 2.2 Profile
Generic Gray Profile
Generic RGB Profile
Gray Gamma 2.2
Grayscale - 20% Dot Gain
Grayscale - Gamma 2.2
HDTV
ISO Coated v2 300% (ECI)
Image Capture Custom Profile
KaiMac.icc
Nikon Adobe RGB 4.0.0.3000
Nikon Apple RGB 4.0.0.3000
Nikon sRGB 4.0.0.3001
Pixar RGB (2005)
Pixar RGB (2008)
ProPhoto RGB
SE2717H/HX
SP2200 Premium Glossy_PK
Scanner Gray Profile
Scanner RGB Profile
U.S. Web Coated (SWOP) v2
c2
c2ci
eciRGB v2
iMac
iMac-1
sRGBsRGB IEC61966-2-1 black scaled
sRGB IEC61966-2.1
sRGB Profile
sRGB Transfer with Display P3 GamutsRGB v1.31 (Canon)

*/
