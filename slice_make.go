package golangrs

func CloneSlice[S ~[]E, E any](sli S) S {
	retval := make(S, len(sli))
	copy(retval, sli)
	return retval
}
