package gompress

func inBound(in, lower, upper int) bool {
	return in >= lower && in <= upper
}
