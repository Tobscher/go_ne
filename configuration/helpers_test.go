package configuration

func arrayIncludes(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
