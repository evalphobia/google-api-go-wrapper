package vision

import (
	"encoding/base64"
	"errors"

	SDK "google.golang.org/api/vision/v1"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log"
)

const (
	serviceName = "vision"
)

// Vision repesents Cloud Vision API client.
type Vision struct {
	service *SDK.Service
	logger  log.Logger
}

// New returns initialized *Vision.
func New(conf config.Config) (*Vision, error) {
	if len(conf.Scopes) == 0 {
		conf.Scopes = append(conf.Scopes, SDK.CloudPlatformScope)
	}
	cli, err := conf.Client()
	if err != nil {
		return nil, err
	}

	svc, err := SDK.New(cli)
	if err != nil {
		return nil, err
	}

	vision := &Vision{
		service: svc,
		logger:  log.DefaultLogger,
	}
	return vision, nil
}

// SetLogger sets internal API logger.
func (v *Vision) SetLogger(logger log.Logger) {
	v.logger = logger
}

// GetFromByte sends image to API and detects all of feature types.
func (v *Vision) GetFromByte(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image: image,
		TypeList: []FeatureType{
			FeatureFace,
			FeatureLandmark,
			FeatureLogo,
			FeatureLabel,
			FeatureText,
			FeatureSafe,
		},
	})
}

// Face sends image to API and detects face feature type.
func (v *Vision) Face(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureFace},
	})
}

// Landmark sends image to API and detects landmark feature type.
func (v *Vision) Landmark(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureLandmark},
	})
}

// Logo sends image to API and detects logo feature type.
func (v *Vision) Logo(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureLogo},
	})
}

// Label sends image to API and detects label feature type.
func (v *Vision) Label(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureLabel},
	})
}

// Text sends image to API and detects OCR feature type.
func (v *Vision) Text(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureText},
	})
}

// Document sends image to API and detects document OCR feature type.
func (v *Vision) Document(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureDocument},
	})
}

// Safe sends image to API and detects safe-search feature type.
func (v *Vision) Safe(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureSafe},
	})
}

// Crop sends image to API and detects crop hints feature type.
func (v *Vision) Crop(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureCrop},
	})
}

// Web sends image to API and detects web document feature type.
func (v *Vision) Web(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureWeb},
	})
}

// Properties sends image to API and detects propaties of images.
func (v *Vision) Properties(image []byte) (*Response, error) {
	return v.Get(&Resource{
		Image:    image,
		TypeList: []FeatureType{FeatureProperties},
	})
}

// Get executes Images.Annotate operation with parameters of Resource.
func (v *Vision) Get(r *Resource) (*Response, error) {
	// parameter validation
	switch {
	case !r.isContainsValidType():
		return nil, errors.New("cannot find FeatureType")
	case !r.isContainsValidImage():
		return nil, errors.New("cannot find valid image")
	}

	requests := createAnnotateImageRequests(r)
	return v.get(&SDK.BatchAnnotateImagesRequest{
		Requests: requests,
	})
}

// get executes Images.Annotate operation.
func (v *Vision) get(req *SDK.BatchAnnotateImagesRequest) (*Response, error) {
	resp, err := v.service.Images.Annotate(req).Do()
	if err != nil {
		v.Errorf("error on `Annotate` operation;  error=%s", err.Error())
	}
	return &Response{resp}, err
}

// Errorf logging error information.
func (v *Vision) Errorf(format string, vv ...interface{}) {
	v.logger.Errorf(serviceName, format, vv...)
}

// createAnnotateImageRequestsForContent creates SDK.AnnotateImageRequest slice.
func createAnnotateImageRequests(r *Resource) []*SDK.AnnotateImageRequest {
	if r.isContainsImageContent() {
		return createAnnotateImageRequestsForContent(r)
	}
	return createAnnotateImageRequestsForGCS(r)
}

// createAnnotateImageRequestsForContent creates SDK.AnnotateImageRequest slice with image data.
func createAnnotateImageRequestsForContent(r *Resource) []*SDK.AnnotateImageRequest {
	features := createFeatures(r)

	var images [][]byte
	switch {
	case r.hasImageList():
		images = r.ImageList
	case r.hasImage():
		images = append(images, r.Image)
	}

	imageContext := r.createImageContext()
	requests := make([]*SDK.AnnotateImageRequest, len(images))
	for i, byt := range images {
		var base64Image string
		switch {
		case r.IsBase64:
			base64Image = string(byt)
		default:
			base64Image = base64.StdEncoding.EncodeToString(byt)
		}

		requests[i] = &SDK.AnnotateImageRequest{
			Features: features,
			Image: &SDK.Image{
				Content: base64Image,
			},
			ImageContext: imageContext,
		}
	}

	return requests
}

// createAnnotateImageRequestsForGCS creates SDK.AnnotateImageRequest slice with image url of Google Cloud Strage.
func createAnnotateImageRequestsForGCS(r *Resource) []*SDK.AnnotateImageRequest {
	features := createFeatures(r)

	var list []string
	switch {
	case r.hasImageURLList():
		list = r.ImageURLList
	case r.hasImage():
		list = append(list, r.ImageURL)
	}

	requests := make([]*SDK.AnnotateImageRequest, len(list))
	for i, url := range list {
		requests[i] = &SDK.AnnotateImageRequest{
			Features: features,
			Image: &SDK.Image{
				Source: &SDK.ImageSource{
					GcsImageUri: url,
				},
			},
		}
	}
	return requests
}

// createFeatures creates []*SDK.Feature from Resource.
func createFeatures(r *Resource) (features []*SDK.Feature) {
	switch {
	case r.hasTypeList():
		size := len(r.TypeList)
		features = make([]*SDK.Feature, size)
		for i, typ := range r.TypeList {
			features[i] = typ.Feature(r.MaxResults)
		}
	case r.hasType():
		features = append(features, r.Type.Feature(r.MaxResults))
	}
	return
}
