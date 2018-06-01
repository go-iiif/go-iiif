package service

/*

  http://palette.davidnewbury.com/

  "service": {
    "@context": "http://palette.davidnewbury.com/vocab/context.json",
    "profile": "http://palette.davidnewbury.com/vocab/iiifpal",
    "label": "Palette automatically generated with a IIIF Palette Server",
    "average": {
      "closest": "#696969",
      "color": "#726e51"
    },
    "palette": [
      {
        "closest": "#6b8e23",
        "color": "#887d35"
      },
      {
        "closest": "#2f4f4f",
        "color": "#332a17"
      },
      {
        "closest": "#808080",
        "color": "#968b68"
      },
      {
        "closest": "#a9a9a9",
        "color": "#a6a694"
      },
      {
        "closest": "#556b2f",
        "color": "#544d18"
      }
    ],
    "reference-closest": "css3"
  }

*/

type PaletteColor struct {
	Color   string `json:"colour"`
	Closest string `json:"closest"`
}

type PaletteService struct {
	Service `json:",omitempty"`
	Context string       `json:"@context"`
	Profile string       `json:"profile"`
	Label   string       `json:"label"`
	Average PaletteColor `json:"average,omitempty"`
	Palette []PaletteColor
}
