package utils

func DeepCopy2D[T comparable](src [][]T) [][]T {
	dst := make([][]T, len(src))

	for i := 0; i < len(dst); i++ {
		dst[i] = DeepCopy1D(src[i])
	}

	return dst
}

func DeepCopy1D[T comparable](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)

	return dst
}

func DeepEqual2D[T comparable](s1, s2 [][]T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if eq := DeepEqual1D(s1[i], s2[i]); !eq {
			return false
		}
	}

	return true
}

func DeepEqual1D[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func MapDeepCopy[K, V comparable](m map[K]V) map[K]V {
	nm := make(map[K]V, len(m))

	for k, v := range m {
		nm[k] = v
	}

	return nm
}