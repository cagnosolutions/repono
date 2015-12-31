package repono

import "bytes"

func formatList(bb [][]byte) []byte {
	if bb != nil {
		bb[0] = append([]byte{'['}, bb[0]...)
		bb[len(bb)-1] = append(bb[len(bb)-1], ']')
		return bytes.Join(bb, []byte{','})
	}
	return nil
}
