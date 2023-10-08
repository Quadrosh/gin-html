package forms

import "fmt"

// Errors holds errors by form field
type Errors map[string][]string

// ErrorText error text by tag
func (e Errors) ErrorText(tag string, param string) string {

	switch tag {
	case "required":
		return "Обязательное поле"
	case "max_length":
		return fmt.Sprintf("Максимальная длина %s симв.", param)
	}

	return tag
}
