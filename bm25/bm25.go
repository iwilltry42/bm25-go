package bm25

import (
	"errors"
	"fmt"
	"log"
	"math"
)

// BM25 is an interface that defines the common methods for all BM25 variants.
type BM25 interface {
	CorpusSize() int
	AvgDocLen() float64
	DocLengths() []int
	IDF(term string) (float64, error)
	GetScores(query []string) ([]float64, error)
	GetBatchScores(query []string, docIDs []int) ([]float64, error)
	GetTopN(query []string, n int) ([]string, error)
}

// Bm25Base is a base struct that holds common fields and methods for all BM25 variants.
type Bm25Base struct {
	corpus     [][]string
	corpusSize int
	avgDocLen  float64
	docLengths []int
	termFreqs  map[string]int
	idfCache   map[string]float64
	tokenizer  func(string) []string
	logger     *log.Logger
}

// NewBM25Base creates a new instance of the Bm25Base struct.
func NewBM25Base(corpus []string, tokenizer func(string) []string, logger *log.Logger) (*Bm25Base, error) {
	if len(corpus) == 0 {
		return nil, errors.New("corpus cannot be empty")
	}

	if tokenizer == nil {
		return nil, errors.New("tokenizer function cannot be nil")
	}

	base := &Bm25Base{
		corpus:    make([][]string, len(corpus)),
		termFreqs: make(map[string]int),
		idfCache:  make(map[string]float64),
		tokenizer: tokenizer,
		logger:    logger,
	}

	var totalDocLen int
	for i, doc := range corpus {
		tokens := tokenizer(doc)
		if len(tokens) == 0 {
			return nil, fmt.Errorf("tokenizer function returned an empty slice for document at index %d", i)
		}
		base.corpus[i] = tokens
		base.docLengths = append(base.docLengths, len(tokens))
		totalDocLen += len(tokens)

		// Use a map or set to ensure each term is only counted once per document
		seenTokens := make(map[string]struct{})
		for _, token := range tokens {
			if _, seen := seenTokens[token]; !seen {
				base.termFreqs[token]++
				seenTokens[token] = struct{}{}
			}
		}
	}

	base.corpusSize = len(corpus)
	base.avgDocLen = float64(totalDocLen) / float64(base.corpusSize)

	if base.logger != nil {
		base.logger.Printf("Corpus size: %d, Average document length: %.2f", base.corpusSize, base.avgDocLen)
	}

	return base, nil
}

// CorpusSize returns the size of the corpus.
func (b *Bm25Base) CorpusSize() int {
	return b.corpusSize
}

// AvgDocLen returns the average document length in the corpus.
func (b *Bm25Base) AvgDocLen() float64 {
	return b.avgDocLen
}

// DocLengths returns the lengths of all documents in the corpus.
func (b *Bm25Base) DocLengths() []int {
	return b.docLengths
}

// IDF returns the inverse document frequency (IDF) of the given term.
func (b *Bm25Base) IDF(term string) (float64, error) {
	if term == "" {
		return 0, errors.New("term cannot be empty")
	}

	if idf, ok := b.idfCache[term]; ok {
		return idf, nil
	}

	termFreq, ok := b.termFreqs[term]
	if !ok {
		b.idfCache[term] = 0.0
		return 0.0, nil
	}

	if termFreq == 0 {
		// Term does not appear in any document, set IDF to 0
		b.idfCache[term] = 0.0
		return 0.0, nil
	}

	if termFreq == b.corpusSize {
		// Term appears in all documents, set IDF to a small positive value
		idf := math.Log(0.5 / (float64(termFreq) + 0.5)) // This will give a small negative value; you can return 0 instead
		b.idfCache[term] = idf
		return idf, nil
	}

	idf := math.Log(((float64(b.corpusSize) - float64(termFreq) + 0.5) / (float64(termFreq) + 0.5)) + 1.0)
	b.idfCache[term] = idf

	if b.logger != nil {
		b.logger.Printf("IDF for term '%s': %.2f", term, idf)
	}

	return idf, nil
}

// GetScores returns the BM25 scores for the given query.
func (b *Bm25Base) GetScores(query []string) ([]float64, error) {
	return nil, errors.New("not implemented")
}

// GetBatchScores returns the BM25 scores for the given query and a subset of documents.
func (b *Bm25Base) GetBatchScores(query []string, docIDs []int) ([]float64, error) {
	return nil, errors.New("not implemented")
}

// GetTopN returns the top N documents for the given query.
func (b *Bm25Base) GetTopN(query []string, n int) ([]string, error) {
	return nil, errors.New("not implemented")
}
