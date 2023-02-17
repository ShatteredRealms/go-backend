package helpers

func ArrayUint64ToUint(in []uint64) []uint {
	out := make([]uint, len(in))
	for x := range in {
		out[x] = uint(in[x])
	}
	return out
}
