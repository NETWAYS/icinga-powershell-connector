package main

import (
	"errors"
	"fmt"
	"io"
)

const (
	Separator = ':'
)

func ParseNetstring(r io.Reader) ([]byte, error) {
	var (
		char   = make([]byte, 1)
		length int
		digit  int
	)

	// Read length from reader.
	for {
		_, err := r.Read(char)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// EOF in length means we reached the end
				return nil, nil
			}

			return nil, err
		}

		b := char[0]

		if b == 10 /* \n */ || b == 13 /* \r */ {
			// ignore line feeds in length.
			continue
		} else if b == Separator {
			break
		}

		digit = int(b - '0')
		if (digit < 0) || (digit > 9) {
			return nil, errors.New("invalid char in netstring length")
		}

		length = length*10 + digit
	}

	// Read netstring content.
	data := make([]byte, length)

	_, err := r.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed reading netstring content: %w", err)
	}

	// Read data end char and ignore possible EOF.
	_, _ = r.Read(char)

	return data, nil
}
