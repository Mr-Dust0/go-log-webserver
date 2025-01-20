package utils

// RemoveDuplicatesInSlice removes duplicates from a slice of any type which is comparable (e.g. int, string, etc.)
func RemoveDuplicatesInSlice[T comparable](s []T) []T {
	bucket := make(map[T]bool)
	var result []T
	for _, str := range s {
		// Check if the key already exists in the map if and it isnt add the value to an list that is returned at the end of the function and do nothing if the key is already present
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}
