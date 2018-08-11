package utils

func ContainsMany(src, content []string) bool {
	for _, contentElem := range content {
		if !Contains(src, contentElem) {
			return false
		}
	}

	return true
}

func Contains(src []string, elem string) bool {

	for _, e := range src {
		if e == elem {
			return true
		}
	}

	return false
}
