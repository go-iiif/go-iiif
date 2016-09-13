package compliance

type ComplianceDetails struct {
	Name      string `json:"name"`
	Syntax    string `json:"syntax"`
	Required  bool   `json:"required"`
	Supported bool   `json:"supported"`
	Match     string `json:"match,omitempty"`
}

type ImageCompliance struct {
	Region   map[string]ComplianceDetails `json:"region"`
	Size     map[string]ComplianceDetails `json:"size"`
	Rotation map[string]ComplianceDetails `json:"rotation"`
	Quality  map[string]ComplianceDetails `json:"quality"`
	Format   map[string]ComplianceDetails `json:"format"`
}

type Compliance interface {
	Formats() []string
	Qualities() []string
	IsValidImageRegion(string) (bool, error)
	IsValidImageSize(string) (bool, error)
	IsValidImageRotation(string) (bool, error)
	IsValidImageQuality(string) (bool, error)
	IsValidImageFormat(string) (bool, error)
	Spec() *Level2ComplianceSpec
}
