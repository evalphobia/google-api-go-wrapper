package vision

import SDK "google.golang.org/api/vision/v1"

// DocumentEntity contains result of Crop Hints.
type DocumentEntity struct {
	Text  string
	Pages []Page
}

// NewDocumentEntity creates DocumentEntity from result of Web Detection.
func NewDocumentEntity(anno *SDK.TextAnnotation) DocumentEntity {
	if anno == nil {
		return DocumentEntity{}
	}

	pages := make([]Page, len(anno.Pages))
	for i, p := range anno.Pages {
		pages[i] = NewPage(p)
	}

	return DocumentEntity{
		Text:  anno.Text,
		Pages: pages,
	}
}

// Page is wrapper struct for SDK.Page.
type Page struct {
	TextProperty
	Blocks     []Block
	Confidence float64
	Height     int64
	Width      int64
}

// NewPage creates Page from SDK.Page.
func NewPage(p *SDK.Page) Page {
	blocks := make([]Block, len(p.Blocks))
	for i, b := range p.Blocks {
		blocks[i] = NewBlock(b)
	}

	return Page{
		TextProperty: NewTextProperty(p.Property),
		Blocks:       blocks,
		Confidence:   p.Confidence,
		Height:       p.Height,
		Width:        p.Width,
	}
}

// Block is wrapper struct for SDK.Block.
type Block struct {
	TextProperty
	BlockType  string
	Confidence float64
	Paragraphs []Paragraph
	Vertices   []Vertex
}

// NewBlock creates Block from SDK.Block.
func NewBlock(b *SDK.Block) Block {
	paragraphs := make([]Paragraph, len(b.Paragraphs))
	for i, p := range b.Paragraphs {
		paragraphs[i] = NewParagraph(p)
	}

	return Block{
		TextProperty: NewTextProperty(b.Property),
		Paragraphs:   paragraphs,
		Vertices:     NewVertices(b.BoundingBox),
		Confidence:   b.Confidence,
		BlockType:    b.BlockType,
	}
}

// Language is wrapper struct for SDK.Language.
type Language struct {
	Confidence   float64
	LanguageCode string
}

// NewLanguage creates Language from SDK.Language.
func NewLanguage(l *SDK.DetectedLanguage) Language {
	return Language{
		Confidence:   l.Confidence,
		LanguageCode: l.LanguageCode,
	}
}

// Paragraph is wrapper struct for SDK.Paragraph.
type Paragraph struct {
	TextProperty
	Confidence float64
	Vertices   []Vertex
	Words      []Word
}

// NewParagraph creates Paragraph from SDK.Paragraph.
func NewParagraph(p *SDK.Paragraph) Paragraph {
	words := make([]Word, len(p.Words))
	for i, w := range p.Words {
		words[i] = NewWord(w)
	}

	return Paragraph{
		TextProperty: NewTextProperty(p.Property),
		Confidence:   p.Confidence,
		Vertices:     NewVertices(p.BoundingBox),
		Words:        words,
	}
}

// Word is wrapper struct for SDK.Word.
type Word struct {
	TextProperty
	Confidence float64
	Vertices   []Vertex
	Symbols    []Symbol
}

// NewWord creates Word from SDK.Word.
func NewWord(w *SDK.Word) Word {
	symbols := make([]Symbol, len(w.Symbols))
	for i, s := range w.Symbols {
		symbols[i] = NewSymbol(s)
	}

	return Word{
		TextProperty: NewTextProperty(w.Property),
		Confidence:   w.Confidence,
		Vertices:     NewVertices(w.BoundingBox),
		Symbols:      symbols,
	}
}

// Symbol is wrapper struct for SDK.Symbol.
type Symbol struct {
	TextProperty
	Text       string
	Confidence float64
	Vertices   []Vertex
}

// NewSymbol creates Symbol from SDK.Symbol.
func NewSymbol(s *SDK.Symbol) Symbol {
	return Symbol{
		TextProperty: NewTextProperty(s.Property),
		Text:         s.Text,
		Confidence:   s.Confidence,
		Vertices:     NewVertices(s.BoundingBox),
	}
}

// TextProperty is wrapper struct for SDK.TextProperty.
type TextProperty struct {
	Break     Break
	Languages []Language
}

// NewTextProperty creates TextProperty from SDK.TextProperty.
func NewTextProperty(p *SDK.TextProperty) TextProperty {
	if p == nil {
		return TextProperty{}
	}

	langs := make([]Language, len(p.DetectedLanguages))
	for i, l := range p.DetectedLanguages {
		langs[i] = NewLanguage(l)
	}

	return TextProperty{
		Break:     NewBreak(p.DetectedBreak),
		Languages: langs,
	}
}

// Break is wrapper struct for SDK.DetectedBreak.
type Break struct {
	IsPrefix bool
	Type     string
}

// NewBreak creates Break from SDK.DetectedBreak.
func NewBreak(b *SDK.DetectedBreak) Break {
	if b == nil {
		return Break{}
	}

	return Break{
		IsPrefix: b.IsPrefix,
		Type:     b.Type,
	}
}
