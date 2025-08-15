package matcha

func Conditional[T any](test bool, happy, sad T) T {
	if test {
		return happy
	}
	return sad
}

func Map[T any, K any](slice []T, iterator func(index int, element T) K) []K {
	output := make([]K, len(slice))

	for i, e := range slice {
		output = append(output, iterator(i, e))
	}
	return output
}
