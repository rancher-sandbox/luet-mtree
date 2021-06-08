package helpers

// WrapErrorMap is a simple helper to fill the state+error of the default response if there is an error
func WrapErrorMap(e error) (map[string]string, error) {
	data := map[string]string{
		"data":  "",
		"error": "",
		"state": "",
	}
	if e != nil {
		data["state"] = "Errors found"
		data["error"] = e.Error()
	} else {
		data["state"] = "All checks succeeded"
	}
	return data, e
}
