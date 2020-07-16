package utils

func StringSliceIntersect(as, bs []string) []string {
	result := make([]string, 0, max(len(as), len(bs)))
	for _, a := range as {
		for _, b := range bs {
			if a == b {
				result = append(result, a)
			}
		}
	}
	return result
}

func StringSliceMerge(as, bs []string) []string {
	result := make([]string, 0, max(len(as), len(bs)))
	for _, a := range as {
		for _, b := range bs {
			if a != b {
				result = append(result, a)
			}
		}
	}
	return result
}

func StringSliceDifference(from, substract []string) []string {
	// from - substract
	result := make([]string, 0, max(len(from), len(substract)))
	for _, item := range substract {
		if !StringSliceContains(from, item) {
			result = append(result, item)
		}
	}
	return result
}

func StringSliceContains(slice []string, item string) bool {
	for _, a := range slice {
		if item == a {
			return true
		}
	}
	return false
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
