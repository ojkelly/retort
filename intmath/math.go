package intmath

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
