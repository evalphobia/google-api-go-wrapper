package vision

import SDK "google.golang.org/api/vision/v1"

// WebEntity contains result of Web Detection.
type WebEntity struct {
	Labels                []WebLabel
	FullMatchingImages    []string // URL
	PartialMatchingImages []string
	VisuallySimilarImages []string
	Entities              []Entity
	Pages                 []WebPage
}

// NewWebEntity creates WebEntity from result of Web Detection.
func NewWebEntity(anno *SDK.WebDetection) WebEntity {
	if anno == nil {
		return WebEntity{}
	}

	return WebEntity{
		Labels:                NewWebLabels(anno.BestGuessLabels),
		FullMatchingImages:    getURLListFromWebImages(anno.FullMatchingImages),
		PartialMatchingImages: getURLListFromWebImages(anno.PartialMatchingImages),
		VisuallySimilarImages: getURLListFromWebImages(anno.VisuallySimilarImages),
		Entities:              newEntitiesFromWebEntity(anno.WebEntities),
		Pages:                 NewWebPages(anno.PagesWithMatchingImages),
	}
}

func getURLListFromWebImages(list []*SDK.WebImage) []string {
	result := make([]string, len(list))
	for i, l := range list {
		result[i] = l.Url
	}
	return result
}

func newEntitiesFromWebEntity(list []*SDK.WebEntity) []Entity {
	result := make([]Entity, len(list))
	for i, l := range list {
		result[i] = Entity{
			Description: l.Description,
			Score:       l.Score,
			EntityID:    l.EntityId,
		}
	}
	return result
}

// GetMatchingURL returns all of matching image url.
func (e WebEntity) GetMatchingURL() []string {
	var result []string
	if len(e.FullMatchingImages) != 0 {
		result = append(result, e.FullMatchingImages...)
	}
	if len(e.PartialMatchingImages) != 0 {
		result = append(result, e.PartialMatchingImages...)
	}
	return result
}

// WebLabel is wrapper strcut for SDK.WebLabel.
type WebLabel struct {
	Label        string
	LanguageCode string
}

// NewWebLabels creates []WebLabel from []SDK.WebLabel.
func NewWebLabels(labels []*SDK.WebLabel) []WebLabel {
	result := make([]WebLabel, len(labels))
	for i, l := range labels {
		result[i] = WebLabel{
			Label:        l.Label,
			LanguageCode: l.LanguageCode,
		}
	}
	return result
}

// WebPage is wrapper struct for SDK.WebPage.
type WebPage struct {
	FullMatchingImages    []string
	PartialMatchingImages []string
	Title                 string
	URL                   string
}

// NewWebPages creates []WebPage from []SDK.WebPage.
func NewWebPages(labels []*SDK.WebPage) []WebPage {
	result := make([]WebPage, len(labels))
	for i, l := range labels {
		result[i] = WebPage{
			Title:                 l.PageTitle,
			URL:                   l.Url,
			FullMatchingImages:    getURLListFromWebImages(l.FullMatchingImages),
			PartialMatchingImages: getURLListFromWebImages(l.PartialMatchingImages),
		}
	}
	return result
}
