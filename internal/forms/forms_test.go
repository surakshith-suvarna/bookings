package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/any-url", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("Got invalid when should have been valid")
	}

}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/any-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest("POST", "/any-url", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows does not have required field when it does")
	}

}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("a")
	if has {
		t.Error("form shows valid when a is blank")
	}
	postedData = url.Values{}
	postedData.Add("a", "value")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("form shows invalid when it is valid")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.MinLength("field_length", 6)
	if form.Valid() {
		t.Error("form shows valid when field does not exist")
	}
	postedData = url.Values{}
	postedData.Add("field_length", "actual length")
	form = New(postedData)
	form.MinLength("field_length", 100)
	if form.Valid() {
		t.Error("form shows valid when field length is less then Minlength of 100")
	}
	isError := form.Errors.Get("field_length")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}
	postedData = url.Values{}
	postedData.Add("field_length", "correct")
	form = New(postedData)
	form.MinLength("field_length", 1)
	if !form.Valid() {
		t.Error("form shows invalid even if the field is correct")
	}
	isError = form.Errors.Get("field_length")
	if isError != "" {
		t.Error("Should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid when the email field is empty")
	}
	postedData = url.Values{}
	postedData.Add("email", "s")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid when email address in invalid")
	}
	postedData = url.Values{}
	postedData.Add("email", "s@test.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("form shows invalid when the email id is valid")
	}
}
