package bm25_test

import (
	"strings"
	"testing"

	"github.com/iwilltry42/bm25-go/bm25"
)

func TestNewBM25L(t *testing.T) {
	corpus := []string{"hello world", "this is a test"}
	tokenizer := func(s string) []string { return strings.Split(s, " ") }

	// Test case: Creating a new BM25L instance with negative k1
	_, err := bm25.NewBM25L(corpus, tokenizer, -1.0, 0.75, nil)
	if err == nil {
		t.Errorf("Expected an error for negative k1, but got nil")
	}

	// Test case: Creating a new BM25L instance with b outside the range [0, 1]
	_, err = bm25.NewBM25L(corpus, tokenizer, 1.2, 1.5, nil)
	if err == nil {
		t.Errorf("Expected an error for b outside the range [0, 1], but got nil")
	}

	// Test case: Creating a new BM25L instance with valid inputs
	_, err = bm25.NewBM25L(corpus, tokenizer, 1.2, 0.75, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBM25LGetScores(t *testing.T) {
	corpus := []string{"hello world", "this is a test"}
	tokenizer := func(s string) []string { return strings.Split(s, " ") }
	bm25, _ := bm25.NewBM25L(corpus, tokenizer, 1.2, 0.75, nil)

	// Test case: Getting scores for an empty query
	_, err := bm25.GetScores([]string{})
	if err == nil {
		t.Errorf("Expected an error for an empty query, but got nil")
	}

	// Test case: Getting scores for a single-term query
	scores, err := bm25.GetScores([]string{"hello"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []float64{0.8109631974066755, 0.0}
	if len(scores) != len(expected) {
		t.Errorf("Expected %d scores, but got %d", len(expected), len(scores))
	}
	for i, score := range scores {
		if score != expected[i] {
			t.Errorf("Expected score %.2f at index %d, but got %.2f", expected[i], i, score)
		}
	}

	// Test case: Getting scores for a multi-term query
	scores, err = bm25.GetScores([]string{"this", "test"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected = []float64{0.0, 1.3862943611198906}
	if len(scores) != len(expected) {
		t.Errorf("Expected %d scores, but got %d", len(expected), len(scores))
	}
	for i, score := range scores {
		if score != expected[i] {
			t.Errorf("Expected score %.2f at index %d, but got %.2f", expected[i], i, score)
		}
	}
}

func TestBM25LGetBatchScores(t *testing.T) {
	corpus := []string{"hello world", "this is a test"}
	tokenizer := func(s string) []string { return strings.Split(s, " ") }
	bm25, _ := bm25.NewBM25L(corpus, tokenizer, 1.2, 0.75, nil)

	// Test case: Getting batch scores for an empty query
	_, err := bm25.GetBatchScores([]string{}, []int{0, 1})
	if err == nil {
		t.Errorf("Expected an error for an empty query, but got nil")
	}

	// Test case: Getting batch scores for an empty document IDs slice
	_, err = bm25.GetBatchScores([]string{"hello"}, []int{})
	if err == nil {
		t.Errorf("Expected an error for an empty document IDs slice, but got nil")
	}

	// Test case: Getting batch scores for invalid document IDs
	_, err = bm25.GetBatchScores([]string{"hello"}, []int{-1, 2})
	if err == nil {
		t.Errorf("Expected an error for invalid document IDs, but got nil")
	}

	// Test case: Getting batch scores for a single-term query
	scores, err := bm25.GetBatchScores([]string{"hello"}, []int{0})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []float64{0.8109631974066755}
	if len(scores) != len(expected) {
		t.Errorf("Expected %d scores, but got %d", len(expected), len(scores))
	}
	for i, score := range scores {
		if score != expected[i] {
			t.Errorf("Expected score %.2f at index %d, but got %.2f", expected[i], i, score)
		}
	}

	// Test case: Getting batch scores for a multi-term query
	scores, err = bm25.GetBatchScores([]string{"this", "test"}, []int{1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected = []float64{1.3862943611198906}
	if len(scores) != len(expected) {
		t.Errorf("Expected %d scores, but got %d", len(expected), len(scores))
	}
	for i, score := range scores {
		if score != expected[i] {
			t.Errorf("Expected score %.2f at index %d, but got %.2f", expected[i], i, score)
		}
	}
}

func TestBM25LGetTopN(t *testing.T) {
	corpus := []string{"hello world", "this is a test"}
	tokenizer := func(s string) []string { return strings.Split(s, " ") }
	bm25, _ := bm25.NewBM25L(corpus, tokenizer, 1.2, 0.75, nil)

	// Test case: Getting top N documents for an empty query
	_, err := bm25.GetTopN([]string{}, 2)
	if err == nil {
		t.Errorf("Expected an error for an empty query, but got nil")
	}

	// Test case: Getting top N documents with n <= 0
	_, err = bm25.GetTopN([]string{"hello"}, 0)
	if err == nil {
		t.Errorf("Expected an error for n <= 0, but got nil")
	}

	// Test case: Getting top N documents for a single-term query
	topDocs, err := bm25.GetTopN([]string{"hello"}, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"hello world"}
	if len(topDocs) != len(expected) {
		t.Errorf("Expected %d top documents, but got %d", len(expected), len(topDocs))
	}
	for i, doc := range topDocs {
		if doc != expected[i] {
			t.Errorf("Expected document '%s' at index %d, but got '%s'", expected[i], i, doc)
		}
	}

	// Test case: Getting top N documents for a multi-term query
	topDocs, err = bm25.GetTopN([]string{"this", "test"}, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected = []string{"this is a test"}
	if len(topDocs) != len(expected) {
		t.Errorf("Expected %d top documents, but got %d", len(expected), len(topDocs))
	}
	for i, doc := range topDocs {
		if doc != expected[i] {
			t.Errorf("Expected document '%s' at index %d, but got '%s'", expected[i], i, doc)
		}
	}
}
