package utils

func Clamp(v int, v0 int, v1 int) int {
	if v < v0 {
		return v0
	}
	if v > v1 {
		return v1
	}
	return v
}
