package helpers

func ArrayOfUint64ToUint(in *[]uint64) *[]uint {
	out := make([]uint, len(*in))
	for idx, val := range *in {
		out[idx] = uint(val)
	}

	return &out
}
