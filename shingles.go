package shingles

import (
	"fmt"
	"hash/crc32"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const (
	UNIGRAM = iota + 1
	BIGRAM
	TRIGRAM
	FOURGRAM
	FIVEGRAM
	SIXGRAM
	SEVENGRAM
	EIGHTGRAM
)

var stopwords = []string{"i", "me", "my", "myself", "we", "our", "ours", "ourselves",
	"you", "your", "yours", "yourself", "yourselves", "he", "him", "his", "himself", "she",
	"her", "hers", "herself", "it", "its", "itself", "they", "them", "their", "theirs",
	"themselves", "what", "which", "who", "whom", "this", "that", "these", "those", "am",
	"is", "are", "was", "were", "be", "been", "being", "have", "has", "had", "having", "do",
	"does", "did", "doing", "a", "an", "the", "and", "but", "if", "or", "because", "as",
	"until", "while", "of", "at", "by", "for", "with", "about", "against", "between", "into",
	"through", "during", "before", "after", "above", "below", "to", "from", "up", "down", "in",
	"out", "on", "off", "over", "under", "again", "further", "then", "once", "here", "there",
	"when", "where", "why", "how", "all", "any", "both", "each", "few", "more", "most", "other",
	"some", "such", "no", "nor", "not", "only", "own", "same", "so", "than", "too", "very", "s",
	"t", "can", "will", "just", "don", "should", "now"}

// NGram is a single n-gram, including it's frequency in the corpora
type NGram struct {
	ngram       string
	occurrences int
}

// Shingles holds a number of NGrams representing a given corpora.
type Shingles struct {
	mtx    sync.Mutex
	n      int
	ngrams map[uint32]*NGram
	hashes []uint32
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func hash(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

func frequency(count int, total int) (frequency float32) {
	return (float32)(count) / (float32)(total)
}

// Extract all sentences from a given input
func sentences(s string) (sentences []string) {
	re := regexp.MustCompile(`(?m)\b([^.,;\-\?\!]{1,})\b`) // A sentence extractor... works surprisingly well, but only for latin characters.
	for _, match := range re.FindAllString(s, -1) {
		sentences = append(sentences, match)
	}
	return
}

// Extract all words from a given input, removing any english stopwords if normalizing
func words(s string, normalize bool) (words []string) {
	re := regexp.MustCompile(`(?m)\b([\w'-]{1,})\b`) // A word extractor... works surprisingly well, but only for latin characters.
	for _, match := range re.FindAllString(s, -1) {
		switch {
		case normalize:
			if contains(stopwords, strings.ToLower(match)) {
				continue
			}
			words = append(words, strings.ToLower(match))
		default:
			words = append(words, match)
		}
	}
	return
}

// Compute n-grams for a set of words given as a slice of strings
func computeNGrams(words []string, length int) (ngrams []string) {
	max := len(words) - length
	for i := 0; i <= max; i++ {
		ngrams = append(ngrams, strings.Join(words[i:i+length], " "))
	}
	return
}

// Add a given n-gram in the form of a string array
func (sh *Shingles) add(ngrams []string) {
	for _, ngram := range ngrams {
		hash := hash(ngram)
		sh.mtx.Lock()
		existing, exists := sh.ngrams[hash]
		switch {
		case exists: // Update existing n-gram
			existing.occurrences = existing.occurrences + 1
		default:
			sh.ngrams[hash] = &NGram{
				ngram:       ngram,
				occurrences: 1,
			}
			sh.hashes = append(sh.hashes, hash)
		}
		sh.mtx.Unlock()
	}

}

//
// The following implements the necessary methods for the sort interface
// This allows us to return a list of ngrams sorted by *occurrences*.
//

// Len returns the number of n-grams
func (sh *Shingles) Len() int {
	sh.mtx.Lock()
	defer sh.mtx.Unlock()
	return len(sh.hashes)
}

// Less returns true when the value at index i is less than the value at index j
func (sh *Shingles) Less(i, j int) bool {
	sh.mtx.Lock()
	defer sh.mtx.Unlock()
	return sh.ngrams[sh.hashes[i]].occurrences < sh.ngrams[sh.hashes[j]].occurrences
}

// Swap two values
func (sh *Shingles) Swap(i, j int) {
	sh.mtx.Lock()
	sh.hashes[i], sh.hashes[j] = sh.hashes[j], sh.hashes[i]
	sh.mtx.Unlock()
}

// Initialize the shingles
func (sh *Shingles) Initialize(n int) {
	sh.mtx.Lock()
	sh.n = n
	sh.ngrams = make(map[uint32]*NGram, 10000)
	sh.mtx.Unlock()
}

// Count the number of ngrams
func (sh *Shingles) Count() int {
	return len(sh.hashes)
}

// Incorporate extracted n-grams. The corpora will be any text, from which sentences and words are extracted.
func (sh *Shingles) Incorporate(corpora string, normalize bool) {
	for _, sentence := range sentences(corpora) {
		sh.add(computeNGrams(words(sentence, normalize), sh.n))
	}
}

// Walk and print all n-grams
func (sh *Shingles) Walk() {
	for hash, ngram := range sh.ngrams {
		fmt.Printf("hash:%v, ngram:'%v', count:%v, frequency:%v\n", hash, ngram.ngram, ngram.occurrences, frequency(ngram.occurrences, sh.Count()))
	}
}

// SortedWalk walks and prints all n-grams, sorted by decreasing frequency
func (sh *Shingles) SortedWalk() {
	sort.Sort(sort.Reverse(sh))
	for _, hash := range sh.hashes {
		fmt.Printf("hash:%v, ngram:'%v', count:%v, frequency:%v\n", hash, sh.ngrams[hash].ngram, sh.ngrams[hash].occurrences, frequency(sh.ngrams[hash].occurrences, sh.Count()))
	}
}
