package compliance

type ComplianceDetails struct {
	Name      string `json:"name"`
	Syntax    string `json:"syntax"`
	Required  bool   `json:"required"`
	Supported bool   `json:"supported"`
	Match     string `json:"match,omitempty"`
}

type Compliance interface {
     IsValidImageRegion(string) bool
     IsValidImageSize(string) bool
     IsValidImageRotation(string) bool
     IsValidImageQuality(string) bool
     IsValidImageFormat(string) bool
     Spec() string
}