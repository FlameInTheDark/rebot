package web

import (
	"fmt"
	"net/url"
	"strconv"
)

// QueryParams is extension of built-in net/url.Values
// It can retrieve value by key with specific Go type
type QueryParams url.Values

// String returns a string (if exists) by key
func (qp QueryParams) String(key string) *string {
	if vals, ok := qp[key]; ok {
		return &vals[0]
	}
	return nil
}

// StringSlice returns []string (if exists) by key
func (qp QueryParams) StringSlice(key string) []string {
	return qp[key]
}

// FloatSlice returns []float64 (if exists) by key
// !!! slice may be nil !!!
func (qp QueryParams) FloatSlice(key string) (result []float64, err error) {
	if vals, ok := qp[key]; ok {
		result = make([]float64, len(vals))
		for i, val := range vals {
			if result[i], err = strconv.ParseFloat(val, 64); err != nil {
				return nil, fmt.Errorf("%q: %w", key, err)
			}
		}
	}
	return
}

// Float returns float64 (if exists) by key
func (qp QueryParams) Float(key string) (*float64, error) {
	if vals, ok := qp[key]; ok {
		f, err := strconv.ParseFloat(vals[0], 64)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", key, err)
		}
		return &f, nil
	}
	return nil, nil
}

// IntSlice returns []int (if exists) by key
func (qp QueryParams) IntSlice(key string) (result []int, err error) {
	if vals, ok := qp[key]; ok {
		result = make([]int, len(vals))
		for i, val := range vals {
			if result[i], err = strconv.Atoi(val); err != nil {
				return nil, fmt.Errorf("%q: %w", key, err)
			}
		}
	}
	return
}

// Int returns int (if exists) by key
func (qp QueryParams) Int(key string) (*int, error) {
	if vals, ok := qp[key]; ok {
		i, err := strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("%q: %w", key, err)
		}
		return &i, nil
	}
	return nil, nil
}

// Int64 returns int64 (if exists) by key
func (qp QueryParams) Int64(key string) (*int64, error) {
	if vals, ok := qp[key]; ok {
		i64, err := strconv.ParseInt(vals[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", key, err)
		}
		return &i64, nil
	}
	return nil, nil
}

// Uint returns uint (if exists) by key
func (qp QueryParams) Uint(key string) (*uint, error) {
	i, err := qp.Int(key)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", key, err)
	}
	if i != nil {
		if *i < 0 {
			return nil, fmt.Errorf("%q: cannot parse \"%d\" to unsigned integer", key, *i)
		}
		u := uint(*i)
		return &u, nil
	}
	return nil, nil
}

// Bool returns bool (if exists) by key
// Query param must be "0" or "1" to be successfully transformed to boolean type
func (qp QueryParams) Bool(key string) (*bool, error) {
	if vals, ok := qp[key]; ok {
		b := new(bool)
		switch vals[0] {
		case "1":
			*b = true
			return b, nil
		case "0":
			*b = false
			return b, nil
		default:
			return nil, fmt.Errorf("%q must be either \"0\" or \"1\"", key)
		}
	}
	return nil, nil
}
