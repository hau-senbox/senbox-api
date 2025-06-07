package menu

import (
	"errors"
	"strings"
)

type Direction string

const (
	Top    Direction = "top"
	Bottom Direction = "bottom"
)

func (d Direction) String() string {
	return string(d)
}

func GetDirectionFromString(typeInString string) (Direction, error) {
	switch strings.ToLower(typeInString) {
	case "top":
		return Top, nil
	case "button":
		return Bottom, nil
	default:
		return "", errors.New("invalid component type")
	}
}
