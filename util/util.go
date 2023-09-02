package util

func CheckError(e error, legalErrors ...error) bool {
	// Permit these error, primary io.EOF while reading input
	for _, legalError := range legalErrors {
		if e == legalError {
			return true
		}
	}
	// other errors causes panic
	if e != nil {
			panic(e)
	}
	return false
}