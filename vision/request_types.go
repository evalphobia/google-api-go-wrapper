package vision

import SDK "google.golang.org/api/vision/v1"

// Resource is parameter struct for Annotate.
type Resource struct {
	Image        []byte
	ImageList    [][]byte
	ImageURL     string
	ImageURLList []string
	IsBase64     bool
	Type         FeatureType
	TypeList     []FeatureType
	MaxResults   int64
}

// isContainsValidImage checks if any images or url are set or not.
func (r *Resource) isContainsValidImage() bool {
	return r.isContainsImageContent() || r.isContainsImageURL()
}

// isContainsImageContent checks if any images are set or not.
func (r *Resource) isContainsImageContent() bool {
	return r.hasImage() || r.hasImageList()
}

func (r *Resource) hasImage() bool {
	return len(r.Image) != 0
}

func (r *Resource) hasImageList() bool {
	return len(r.ImageList) != 0
}

// isContainsImageURL checks if any image url are set or not.
func (r *Resource) isContainsImageURL() bool {
	return r.hasImageURL() || r.hasImageURLList()
}

func (r *Resource) hasImageURL() bool {
	return r.ImageURL != ""
}

func (r *Resource) hasImageURLList() bool {
	return len(r.ImageURLList) != 0
}

// isContainsValidType checks if any FeatureType are set or not.
func (r *Resource) isContainsValidType() bool {
	return r.hasType() || r.hasTypeList()
}

func (r *Resource) hasType() bool {
	return string(r.Type) != ""
}

func (r *Resource) hasTypeList() bool {
	return len(r.TypeList) != 0
}

// FeatureType is type for feature.
type FeatureType string

// FeatureType list
const (
	FeatureUnspecified FeatureType = "TYPE_UNSPECIFIED"
	FeatureFace        FeatureType = "FACE_DETECTION"
	FeatureLandmark    FeatureType = "LANDMARK_DETECTION"
	FeatureLogo        FeatureType = "LOGO_DETECTION"
	FeatureLabel       FeatureType = "LABEL_DETECTION"
	FeatureText        FeatureType = "TEXT_DETECTION"
	FeatureDocument    FeatureType = "DOCUMENT_TEXT_DETECTION"
	FeatureSafe        FeatureType = "SAFE_SEARCH_DETECTION"
	FeatureProperties  FeatureType = "IMAGE_PROPERTIES"
	FeatureCrop        FeatureType = "CROP_HINTS"
	FeatureWeb         FeatureType = "WEB_DETECTION"
)

func (f FeatureType) String() string {
	return string(f)
}

// Feature creates *SDK.Feature.
func (f FeatureType) Feature(max int64) *SDK.Feature {
	return &SDK.Feature{
		MaxResults: max,
		Type:       f.String(),
	}
}
