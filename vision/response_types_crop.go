package vision

import SDK "google.golang.org/api/vision/v1"

// CropEntity contains result of Crop Hints.
type CropEntity struct {
	Hints []CropHint
}

// NewCropEntity creates CropEntity from result of Web Detection.
func NewCropEntity(anno *SDK.CropHintsAnnotation) CropEntity {
	if anno == nil {
		return CropEntity{}
	}

	hints := make([]CropHint, len(anno.CropHints))
	for i, h := range anno.CropHints {
		hints[i] = NewCropHint(h)
	}
	// anno.CropHints.

	return CropEntity{
		Hints: hints,
	}
}

// CropHint is wrapper struct for SDK.CropHint.
type CropHint struct {
	Vertices           []Vertex
	Confidence         float64
	ImportanceFraction float64
}

// NewCropHint creates CropHint from SDK.CropHint.
func NewCropHint(h *SDK.CropHint) CropHint {
	return CropHint{
		Vertices:           NewVertices(h.BoundingPoly),
		Confidence:         h.Confidence,
		ImportanceFraction: h.ImportanceFraction,
	}
}
