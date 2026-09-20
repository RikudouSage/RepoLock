package http

import (
	"encoding/json"
	"fmt"
	"io"
)

func ParseBody[TResult any](body io.Reader) (TResult, error) {
	var result TResult

	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed parsing body: %w", err)
	}

	return result, nil
}
