// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
	"strings"
	"regexp"
)

// CountWordFrequency takes a string containing multiple words and returns
// a map where each key is a word and the value is the number of times that
// word appears in the string. The comparison is case-insensitive.
//
// Words are defined as sequences of letters and digits.
// All words are converted to lowercase before counting.
// All punctuation, spaces, and other non-alphanumeric characters are ignored.
//
// For example:
// Input: "The quick brown fox jumps over the lazy dog."
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}


func CountWordFrequency(text string) map[string]int {
   cleanHyphenRegex := regexp.MustCompile(`[-]`) 
   cleanEmptySpaceRegex := regexp.MustCompile(`[\t\n]`)
   cleanPunctuationRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	
	// Your implementation here
	if len(text) >= 10000 {
	    return nil
	}
	
	word_frequency := make(map[string]int, 1000)
	
	cleanedHyphenText := cleanHyphenRegex.ReplaceAllString(text, " ")
	cleanedSpaceText := cleanEmptySpaceRegex.ReplaceAllString(cleanedHyphenText, " ")
	cleanedPunctuationText := cleanPunctuationRegex.ReplaceAllString(cleanedSpaceText, "")
	
	words := strings.Split(strings.ToLower(cleanedPunctuationText), " ")
	
	for _, word := range words {
	    if len(word) > 0 {
	    _, exists := word_frequency[word]
	    if exists {
	        word_frequency[word]++
	    } else {
	        word_frequency[word] = 1
	    }
	    }
	}
	
	
	return word_frequency 
} 