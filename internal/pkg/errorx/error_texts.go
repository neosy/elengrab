package errorx

import (
	"strings"
)

// textErrorsSeparatorDefault, separator for combining error texts into one
var textErrorsSeparatorDefault = "; "

type ErrorTexts struct {
	texts     []string
	textsMap  map[string]struct{}
	separator string
}

// SetTextErrorsSeparatorDefault set default separator value
func SetTextErrorsSeparatorDefault(sep string) {
	textErrorsSeparatorDefault = sep
}

// SetSeparator set separator value
func (errTxts *ErrorTexts) SetSeparator(sep string) {
	errTxts.separator = sep
}

// NewErrorTexts creating an ErrorTexts object
func NewErrorTexts() (errTexts *ErrorTexts) {
	texts := make([]string, 0, 4)
	textsMap := make(map[string]struct{})

	errTexts = &ErrorTexts{
		texts:     texts,
		textsMap:  textsMap,
		separator: textErrorsSeparatorDefault,
	}

	return errTexts
}

// Add adding text
func (errTxts *ErrorTexts) Add(txt string) *ErrorTexts {
	if errTxts.Contains(txt) {
		return errTxts
	}

	if txt != "" {
		errTxts.texts = append(errTxts.texts, txt)
		errTxts.textsMap[txt] = struct{}{}
	}

	return errTxts
}

// AddErr adding text from error
func (errTxts *ErrorTexts) AddErr(err error) *ErrorTexts {
	if err != nil {
		errTxts.Add(err.Error())
	}

	return errTxts
}

// AddErrs adding text from errors
func (errTxts *ErrorTexts) AddErrs(errs ...error) *ErrorTexts {
	if len(errs) == 0 {
		return errTxts
	}

	for _, err := range errs {
		if err == nil {
			continue
		}

		errTxts.AddErr(err)
	}

	return errTxts
}

// AddUnwrapErr analyze the errors and then add one by one
func (errTxts *ErrorTexts) AddUnwrapErr(err error) *ErrorTexts {
	if err != nil {
		errs := UnwrapUnique(err)
		for _, e := range errs {
			errTxts.AddErr(e)
		}
	}

	return errTxts
}

// Join combining texts from errors into one line with a separator
func (errTxts *ErrorTexts) Join() (text string) {
	if len(errTxts.texts) > 0 {
		text = strings.Join(errTxts.texts, errTxts.separator)
	}

	return
}

// Contains returns true if the text is found
func (errTxts *ErrorTexts) Contains(txt string) (exists bool) {
	_, exists = errTxts.textsMap[txt]

	return
}

// ContainsErr returns true if the text from the error is found
func (errTxts *ErrorTexts) ContainsErr(err error) (exists bool) {
	if err != nil {
		exists = errTxts.Contains(err.Error())
	}

	return
}
