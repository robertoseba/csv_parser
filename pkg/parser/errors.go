package parser

import "errors"

var (
	ErrInvalidRow = errors.New("invalid row")
	ErrInvalidCol = errors.New("invalid column")
)
