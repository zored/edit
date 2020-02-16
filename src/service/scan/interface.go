package scan

type IScanner interface {
	Scan() bool
	Text() string
	Err() error
}
