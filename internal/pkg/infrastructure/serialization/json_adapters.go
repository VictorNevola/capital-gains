package serialization

import "strconv"

type KeepZero float64

func (f KeepZero) MarshalJSON() ([]byte, error) {
	if float64(f) == float64(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 2, 32)), nil
	}

	return []byte(strconv.FormatFloat(float64(f), 'f', -2, 32)), nil
}
