package dbs

import (
	"errors"
)

var (
	// bucketNews хранит новости.
	bucketNews = []byte("news")
)

var (
	// ErrInternal — внутренняя ошибка.
	ErrInternal = errors.New("internal error")

	// ErrNotFound — запись не найдена.
	ErrNotFound = errors.New("not found")
)

type NewsPiece struct {
	Header string
	Date   int64 // Unix.
}
