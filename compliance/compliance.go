package compliance

/*

Things I am not sure about include the relationship of level/*.go and compliance/*.go which are
very much related but somehow seem like they should be in separate namespaces. I'm not sure...
(20160912/thisisaaronland)

*/

type ComplianceDetails struct {
	Name      string `json:"name"`
	Syntax    string `json:"syntax"`
	Required  bool   `json:"required"`
	Supported bool   `json:"supported"`
	Match     string `json:"match,omitempty"`
}

type HTTPComplianceDetails struct {
	Name      string `json:"name"`
	Required  bool   `json:"required"`
	Supported bool   `json:"supported"`
}

type HTTPCompliance map[string]HTTPComplianceDetails

/*
   BaseUriRedirect map[string]HTTPComplianceDetails	`json:"baseUriRedirect"`
   CORS map[string]HTTPComplianceDetails	`json:"cors"`
   JsonLdMediaType map[string]HTTPComplianceDetails	`json:"jsonldMediaType"`
   ProfileLinkHeader map[string]HTTPComplianceDetails	`json:"profileLinkHeader"`
   CanonicalLinkHeader map[string]HTTPComplianceDetails	`json:"canonicalLinkHeader"`
*/

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
