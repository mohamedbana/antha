package wtype

import (
	"errors"
	"fmt"
)

// consts for liquid handling planner errors

const (
	LH_OK = iota
	LH_ERR_NO_DECK_SPACE
	LH_ERR_NO_TIPS
	LH_ERR_NOT_IMPLEMENTED
	LH_ERR_OTHER
)

func ErrorName(code int) string {
	errornames := [...]string{"LH_OK", "LH_ERR_NO_DECK_SPACE", "LH_ERR_NO_TIPS", "LH_ERR_NOT_IMPLEMENTED", "LH_ERR_OTHER"}
	return errornames[code]
}

func ErrorDesc(code int) string {
	errorboilerplate := [...]string{"no problem", "not sufficient deck space to fit all required items; this may be due to constraints", "ran out of tips", "a required command is not implemented", ""}
	return errorboilerplate[code]
}

func LHError(code int, detail string) error {
	s := fmt.Sprintf("%d (%s): %s -- %s", code, ErrorName(code), ErrorDesc(code), detail)

	return errors.New(s)
}
