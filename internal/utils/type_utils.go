package utils

type Scannable interface {
	Scan(dest ...interface{}) error
}
