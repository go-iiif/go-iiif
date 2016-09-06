package compliance

import (
	"encoding/json"
)

var level2_spec = `{
    "image": {
    	     "region": {
	     	       "full": { "syntax": "full",        "required": true, "match": "^full$" },
		       "regionByPx": { "syntax": "x,y,w,h",     "required": true, "match": "^\\d+\\,\\d+\\,\\d+\\,\\d+$" },
		       "regionByPct": {  "syntax": "pct:x,y,w,h", "required": true, "match": "^pct\\:\\d+\\,\\d+\\,\\d+\\,\\d+$" },
		       "regionSquare": { "syntax": "square",      "required": false, "match": "^square$" }
	     },
	     "size": {
	     		"full": {              "syntax": "full",  "required": true, "supported": true, "match": "^full$" },
	     		"max": {               "syntax": "max",   "required": false, "supported": true, "match": "^max$" },
	     		"sizeByW": {           "syntax": "w,",    "required": true, "supported": true, "match": "^\\d+\\,$" },			
	     		"sizeByH": {           "syntax": ",h",    "required": true, "supported": true, "match": "^\\,\\d+$" },
	     		"sizeByPct": {         "syntax": "pct:n", "required": true, "supported": true, "match": "^pct\\:\\d+(\\.\\d+)?$" },			
	     		"sizeByConfinedWh": {  "syntax": "!w,h",  "required": true, "supported": true, "match": "" },
	     		"sizeByDistortedWh": { "syntax": "w,h",   "required": true, "supported": true, "match": "" },			
	     		"sizeByWh": {          "syntax": "w,h",   "required": true, "supported": true, "match": "" },
	     		"sizeAboveFull": {     "syntax": "",      "required": false, "supported": false, "match": "" }
	     },
	     "rotation": {
	     		"none": {              "syntax": "0",          "required": true, "supported": true, "match": "" },
	     		"rotationBy90s": {     "syntax": "90,180,270", "required": true, "supported": true, "match": "" },
	     		"rotationArbitrary": { "syntax": "",           "required": false, "supported": true, "match": "" },			
	     		"mirroring": {         "syntax": "!n",         "required": true, "supported": true, "match": "" }
	     },
	     "quality": {
	     		"default": { "syntax": "default", "required": true, "supported": true, "match": "^default$" },
	     		"color": { "syntax": "color",   "required": false, "supported": false, "match": "^colou?r$" },
	     		"gray": { "syntax": "gray",    "required": false, "supported": false, "match": "gr(?:e|a)y$" },			
	     		"bitonal": { "syntax": "bitonal", "required": true, "supported": true, "match": "^bitonal$" }
             },
	     "format": {
	     	       "jpg": { "syntax": "jpg",  "required": true, "supported": true, "match": "^jpe?g$" },
       	     	       "png": { "syntax": "png",  "required": true, "supported": true, "match": "^png$" },
       	     	       "tif": { "syntax": "tif",  "required": false, "supported": false, "match": "^tiff?$" },
      	     	       "gif": { "syntax": "gif",  "required": false, "supported": false, "match": "^gif$" },
       	     	       "pdf": { "syntax": "pdf",  "required": false, "supported": false, "match": "^pdf$" },
      	     	       "jp2": { "syntax": "jp2",  "required": false, "supported": false, "match": "^jp2$" },
       	     	       "webp": { "syntax": "webp", "required": false, "supported": false, "match": "^webp$" }
	     }	     
    }
}`

/*
    "http": [
    	    { "feature": "base URI redirects", "name": "baseUriRedirect", "required": true, "supported": 1 },
	    { "feature": "CORS", "name": "cors", "required": true, "supported": 1 },
	    { "feature": "json-ld media type", "name": "jsonldMediaType", "required": true, "supported": 1 },
	    { "feature": "profile link header", "name": "profileLinkHeader", "required": false, "supported": 0 },
	    { "feature": "canonical link header", "name": "canonicalLinkHeader", "required": false, "supported": 0 }
    ]
*/

type Level2ComplianceSpec struct {
	Image ImageCompliance `json:"image"`
}

type Level2Compliance struct {
	Compliance
	spec Level2ComplianceSpec
}

func NewLevel2Compliance() (*Level2Compliance, error) {

	spec := Level2ComplianceSpec{}
	err := json.Unmarshal([]byte(level2_spec), &spec)

	if err != nil {
		return nil, err
	}

	compliance := Level2Compliance{
		spec: spec,
	}

	return &compliance, nil
}

func (c *Level2Compliance) IsValidImageRegion(region string) bool {

	return true
}

func (c *Level2Compliance) IsValidImageSize(size string) bool {

	return true
}

func (c *Level2Compliance) IsValidImageRotation(rotation string) bool {

	return true
}

func (c *Level2Compliance) IsValidImageQuality(quality string) bool {

	return true
}

func (c *Level2Compliance) IsValidImageFormat(format string) bool {

	return true
}

func (c *Level2Compliance) Spec() ([]byte, error) {

	return json.Marshal(c.spec)
}
