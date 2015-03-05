package silkroad

// StringInSlice looks if a string is in a string array
func StringInSlice(array []string, item string) bool {
	for _, i := range array {
		if i == item {
			return true
		}
	}
	return false
}
