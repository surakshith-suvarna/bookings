package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

//Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

//Valid return true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//New initialises form struct (Initial Form load)
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

//Required checks if required fields have value
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s cannot be empty", field))
		}
	}
}

//Has checks if post field is available and its not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)

	if x == "" {
		//f.Errors.Add(field, "Field cannot be empty.")
		return false
	}
	return true
}

//MinLength checks the minimum length of string
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("The field requires minimum of %d characters", length))
		return false
	}
	return true
}

//IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email")
	}
}
