package util

func StringSplitN(msg string, sz int) []string {
	strings := []string{}

	if sz >= len(msg) {
		strings = append(strings, msg)
	} else {
		ba := []byte(msg)
		start := 0
		end := start + sz
		for end <= len(ba) {
			strings = append(strings, string(ba[start:end]))
			start = end
			end = end + sz
			if end > len(ba) {
				strings = append(strings, string(ba[start:]))
			}
		}
	}

	return strings
}
