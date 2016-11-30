package vision

import SDK "google.golang.org/api/vision/v1"

// Entity contains general result of response.
type Entity struct {
	Description string
	Score       float64
}

// SafeEntity contains result of SafeSearch.
type SafeEntity struct {
	Adult    Likelihood
	Spoof    Likelihood
	Medical  Likelihood
	Violence Likelihood
}

// NewSafeEntity creates SafeEntity from result of SafeSearchAnnotation.
func NewSafeEntity(anno *SDK.SafeSearchAnnotation) SafeEntity {
	return SafeEntity{
		Adult:    NewLikelihood(anno.Adult),
		Spoof:    NewLikelihood(anno.Spoof),
		Medical:  NewLikelihood(anno.Medical),
		Violence: NewLikelihood(anno.Violence),
	}
}

// Likelihood list of string result.
const (
	LikelihoodTextUnknown      = "UNKNOWN"
	LikelihoodTextVeryUnlikely = "VERY_UNLIKELY"
	LikelihoodTextUnlikely     = "UNLIKELY"
	LikelihoodTextPossible     = "POSSIBLE"
	LikelihoodTextLikely       = "LIKELY"
	LikelihoodTextVeryLikely   = "VERY_LIKELY"
)

// Likelihood is likelihood of detection result.
type Likelihood int

// Likelihood list.
const (
	LikelihoodError Likelihood = iota
	LikelihoodUnknown
	LikelihoodVeryUnlikely
	LikelihoodUnlikely
	LikelihoodPossible
	LikelihoodLikely
	LikelihoodVeryLikely
)

// NewLikelihood creates Likelihood.
func NewLikelihood(s string) Likelihood {
	switch s {
	case LikelihoodTextUnknown:
		return LikelihoodUnknown
	case LikelihoodTextVeryUnlikely:
		return LikelihoodVeryUnlikely
	case LikelihoodTextUnlikely:
		return LikelihoodUnlikely
	case LikelihoodTextPossible:
		return LikelihoodPossible
	case LikelihoodTextLikely:
		return LikelihoodLikely
	case LikelihoodTextVeryLikely:
		return LikelihoodVeryLikely
	default:
		return LikelihoodError
	}
}

func (l Likelihood) String() string {
	switch l {
	case LikelihoodUnknown:
		return LikelihoodTextUnknown
	case LikelihoodVeryUnlikely:
		return LikelihoodTextVeryUnlikely
	case LikelihoodUnlikely:
		return LikelihoodTextUnlikely
	case LikelihoodPossible:
		return LikelihoodTextPossible
	case LikelihoodLikely:
		return LikelihoodTextLikely
	case LikelihoodVeryLikely:
		return LikelihoodTextVeryLikely
	default:
		return ""
	}
}

// IsError checks if Likelifood is error.
func (l Likelihood) IsError() bool {
	return l == LikelihoodError
}

// Unknown checks if Likelifood is more than Unknown.
func (l Likelihood) Unknown() bool {
	return l >= LikelihoodUnknown
}

// VeryUnlikely checks if Likelifood is more than VeryUnlikely.
func (l Likelihood) VeryUnlikely() bool {
	return l >= LikelihoodVeryUnlikely
}

// Unlikely checks if Likelifood is more than Unlikely.
func (l Likelihood) Unlikely() bool {
	return l >= LikelihoodUnlikely
}

// Possible checks if Likelifood is more than Possible.
func (l Likelihood) Possible() bool {
	return l >= LikelihoodPossible
}

// Likely checks if Likelifood is more than Likely.
func (l Likelihood) Likely() bool {
	return l >= LikelihoodLikely
}

// VeryLikely checks if Likelifood is more than VeryLikely.
func (l Likelihood) VeryLikely() bool {
	return l >= LikelihoodVeryLikely
}
