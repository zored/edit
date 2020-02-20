package scanners

import (
	"fmt"
)

// Return this error to break scanner execution.
var BreakErr = fmt.Errorf("scan breaked")

// Basic usage of scanner.
func ScanAll(scanner IScanner, f func(string) error) error {
	getScannerError := scanner.Err

	for scanner.Scan() {
		text := scanner.Text()
		err := f(text)

		switch err {
		case BreakErr:
			return getScannerError()
		case nil:
			continue
		default:
			return err
		}
	}

	return getScannerError()
}
