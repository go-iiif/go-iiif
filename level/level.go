package level

/*

Things I am not sure about include the relationship of level/*.go and compliance/*.go which are
very much related but somehow seem like they should be in separate namespaces. I'm not sure...
(20160912/thisisaaronland)

*/

import (
	"fmt"
	iiifcompliance "github.com/go-iiif/go-iiif/v4/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifprofile "github.com/go-iiif/go-iiif/v4/profile"
	_ "log"
)

type Level interface {
	Compliance() iiifcompliance.Compliance
	Profile() *iiifprofile.Profile
}

func NewLevelFromConfig(config *iiifconfig.Config, endpoint string) (Level, error) {

	compliance := config.Level.Compliance

	switch compliance {
	case "0":
		return NewLevel0(config, endpoint)
	case "2":
		return NewLevel2(config, endpoint)
	default:
		return nil, fmt.Errorf("Invalid compliance level '%s'", compliance)
	}
}
