package main

import (
	"strconv"
)

func parseFloat(x *string) *float64 {
	if len(*x) == 0 {
		return nil
	}
	if f, err := strconv.ParseFloat(*x, 64); err != nil {
		panic(err)
	} else {
		return &f
	}
}

func parseBool(x *string) *bool {
	if len(*x) == 0 {
		return nil
	}
	if f, err := strconv.ParseBool(*x); err != nil {
		panic(err)
	} else {
		return &f
	}
}

func parseStringSlice(x *string) []string {
	if len(*x) == 0 {
		return nil
	}

	return []string{*x}
}
