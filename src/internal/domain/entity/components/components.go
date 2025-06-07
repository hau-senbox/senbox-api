package components

import (
	"errors"
	"strings"
)

type ComponentType string

const (
	ButtonURL  ComponentType = "button_url"
	ButtonForm ComponentType = "button_form"
)

func (c ComponentType) String() string {
	return string(c)
}

func GetComponentTypeFromString(typeInString string) (ComponentType, error) {
	switch strings.ToLower(typeInString) {
	case "button_url":
		return ButtonURL, nil
	case "button_form":
		return ButtonForm, nil
	default:
		return "", errors.New("invalid component type")
	}
}
