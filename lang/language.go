package lang

import (
	"golang.org/x/text/language"
)

type PreferredLanguages []language.Tag

func newLanguage(langs []string) []language.Tag {
	tags := make([]language.Tag, len(langs))
	for i, lang := range langs {
		tags[i] = language.Make(lang)
	}
	return tags
}
func NewPreferredLanguages(langs ...string) PreferredLanguages {
	return PreferredLanguages(newLanguage(langs))
}

func NewPreferredLanguagesFromAcceptLanguages(s string) (PreferredLanguages, error) {
	langs, _, err := language.ParseAcceptLanguage(s)
	if err != nil {
		return nil, err
	}

	return PreferredLanguages(langs), nil
}
