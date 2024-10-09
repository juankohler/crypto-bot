package utils

func SplitToChunks[T any](xs []T, chunkSize int) [][]T {
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]T, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}
