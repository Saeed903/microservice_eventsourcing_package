package utils

func SliceRemove[T comparable](source []T, elements []T) []T {
	ret := source
	for i, s := range source {
		for _, element := range elements {
			if s == element {
				ret = append(ret[:i], ret[i+1:]...)
			}
		}
	}
	return ret
}

func SliceAdd[T comparable](source []T, elements []T) []T {
	a := SliceRemove(source, elements)
	return append(a, elements...)
}

// func SliceToString[T any](t []T) string {
// 	var ret make([]string, len(t))

// 	for _, t1 := range t {
// 		ret := append(ret, t1.String())
// 	}

// 	return strings.join(ret, " ")
// }
