package stringvalidator

import "strings"

type StringValidator struct {
	badWordsMap map[string]struct{}
}

func NewValidator() StringValidator {
	return StringValidator{
		badWordsMap: make(map[string]struct{}),
	}
}

func (sv *StringValidator) AddWord(word string) {
	sv.badWordsMap[word] = struct{}{}
}

func (sv *StringValidator) Clean(msg string) string {
	words := strings.Split(msg, " ")
	newMsg := ""
	for _, word := range words {
		lower := strings.ToLower(word)
		if _, ok := sv.badWordsMap[lower]; ok {
			newMsg += " ****"
		} else {
			newMsg += " " + word
		}
	}

	return newMsg
}
