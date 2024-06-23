package validation

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validatorEn "github.com/go-playground/validator/v10/translations/en"

	errorpkg "github.com/dostonshernazarov/resume_maker/internal/errors"
)

func Validator(s interface{}) error {
	var (
		eng      = english.New()
		uni      = ut.New(eng, eng)
		validate = validator.New()
	)

	trans, found := uni.GetTranslator("en")
	if !found {
		return errors.New("Validator translator not found")
	}

	if err := validatorEn.RegisterDefaultTranslations(validate, trans); err != nil {
		return err
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	customValidation := NewCustomValidation(validate)
	if err := validate.RegisterValidation("phone_uz", customValidation.PhoneUz); err != nil {
		return nil
	}

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		errValidation := errorpkg.NewErrValidation()
		errValidation.Err = err
		for _, e := range errs {
			errValidation.Errors[e.Field()] = strings.Replace(e.Translate(trans), e.Field(), "", 1)
		}
		return errValidation
	}
	return nil
}

type customValidation struct {
	Validate *validator.Validate
}

func NewCustomValidation(validate *validator.Validate) *customValidation {
	return &customValidation{Validate: validate}
}

func (v *customValidation) PhoneUz(fl validator.FieldLevel) bool {
	// get value
	phone := strings.TrimSpace(fl.Field().String())
	// parse our phone number
	isMatch, err := regexp.MatchString("^[9]{1}[9]{1}[8]{1}(?:77|88|93|94|90|91|95|93|99|97|98|33)[0-9]{7}$", phone)
	if err != nil {
		return false
	}
	return isMatch
}

func IsValidEmail(email string) bool {
	// Regular expression for basic email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func IsValidPassword(password string) bool {
	// Check password length
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasDigit := false

	// Check if password contains at least one letter and one digit
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
