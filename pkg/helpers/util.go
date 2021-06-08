package helpers

// WrapErrorMap is a simple helper to fill the state+error of the default response if there is an error
func WrapErrorMap(m map[string]string, e error) (map[string]string, error) {
	if e != nil {
		m["state"] = "Errors found"
		m["error"] = e.Error()
	} else {
		m["state"] = "All checks succeeded"
	}
	return m, e
}
