package common

import (
	"bytes"
	"strconv"
	"strings"
)

const SequenceDelimiter = ";"

type IntSequence []int

func (s *IntSequence) UnmarshalText(text []byte) error {
	values := bytes.Split(text, []byte(SequenceDelimiter))
	result := make([]int, len(values))

	for i, v := range values {
		r, err := strconv.Atoi(string(v))
		if err != nil {
			return err
		}

		result[i] = r
	}

	*s = result

	return nil
}

func (s IntSequence) MarshalText() ([]byte, error) {
	values := make([]string, len(s))

	for i, v := range s {
		values[i] = strconv.Itoa(v)
	}

	result := strings.Join(values, ";")

	return []byte(result), nil
}

func absDiff(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
