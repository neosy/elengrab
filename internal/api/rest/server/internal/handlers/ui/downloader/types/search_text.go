package types

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type SearchText string

const minSearchTextLength = 2

func (searchText SearchText) ValidateLength() error {
	if utf8.RuneCountInString(string(searchText)) < 2 {
		return errors.New("search text must contain at least 2 characters")
	}
	return nil
}

func (searchText SearchText) IsLongEnough() bool {
	return utf8.RuneCountInString(string(searchText)) >= minSearchTextLength
}

func (text SearchText) Validate() error {
	if !text.IsLongEnough() {
		return errors.New("search text must contain at least 2 characters")
	}
	return nil
}

func (text SearchText) IsValidate() bool {
	return text.Validate() == nil
}

func (text SearchText) String() string {
	return strings.TrimSpace(string(text))
}
