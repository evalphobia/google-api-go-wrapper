package vision

import SDK "google.golang.org/api/vision/v1"

// Response contains response from Annotate API.
type Response struct {
	*SDK.BatchAnnotateImagesResponse
}

// Len returns size of response.
// This coinsides the number of images in request.
func (r *Response) Len() int {
	if r.BatchAnnotateImagesResponse == nil {
		return 0
	}
	return len(r.BatchAnnotateImagesResponse.Responses)
}

// FaceResult returns all of detection score(0-1) from face detection results.
func (r *Response) FaceResult() []float64 {
	if r.Len() == 0 {
		return nil
	}

	var list []float64
	for _, resp := range r.Responses {
		for _, anno := range resp.FaceAnnotations {
			list = append(list, anno.DetectionConfidence)
		}
	}
	return list
}

// LandmarkResult returns all of scores from landmark detection results.
func (r *Response) LandmarkResult() []Entity {
	if r.Len() == 0 {
		return nil
	}

	var list []Entity
	for _, resp := range r.Responses {
		for _, anno := range resp.LandmarkAnnotations {
			list = append(list, Entity{
				Description: anno.Description,
				Score:       anno.Score,
			})
		}
	}
	return list
}

// LogoResult returns all of scores from logo detection results.
func (r *Response) LogoResult() []Entity {
	if r.Len() == 0 {
		return nil
	}

	var list []Entity
	for _, resp := range r.Responses {
		for _, anno := range resp.LogoAnnotations {
			list = append(list, Entity{
				Description: anno.Description,
				Score:       anno.Score,
			})
		}
	}
	return list
}

// LabelResult returns all of scores from label detection results.
func (r *Response) LabelResult() []Entity {
	if r.Len() == 0 {
		return nil
	}

	var list []Entity
	for _, resp := range r.Responses {
		for _, anno := range resp.LabelAnnotations {
			list = append(list, Entity{
				Description: anno.Description,
				Score:       anno.Score,
			})
		}
	}
	return list
}

// TextResult returns all of text from OCR results.
func (r *Response) TextResult() []string {
	if r.Len() == 0 {
		return nil
	}

	var list []string
	for _, resp := range r.Responses {
		for _, anno := range resp.TextAnnotations {
			list = append(list, anno.Description)
		}
	}
	return list
}

// SafeResult returns all of text from safe-search results.
func (r *Response) SafeResult() []SafeEntity {
	if r.Len() == 0 {
		return nil
	}

	list := make([]SafeEntity, r.Len())
	for i, resp := range r.Responses {
		list[i] = NewSafeEntity(resp.SafeSearchAnnotation)
	}
	return list
}
