package model

import (
	"fmt"
	"strings"
	"time"
)

const timeLayout string = "2006-01-02T15:04:05Z07:00"

// Cookie stores the cookie information
type Cookie struct {
	Name string
	Date string
}

// ToCookie converts a givens string to Cookie
func ToCookie(s string) (Cookie, error) {
	values := strings.Split(s, ",")
	if len(values) != 2 {
		return Cookie{}, fmt.Errorf("invalid input string provided")
	}

	values[0] = strings.TrimSpace(values[0])
	values[1] = strings.TrimSpace(values[1])

	date, err := time.Parse(timeLayout, values[1])
	if err != nil {
		return Cookie{}, fmt.Errorf("invalid time format used: %w", err)
	}

	return Cookie{Name: values[0], Date: date.Format("2006-01-02")}, nil
}
