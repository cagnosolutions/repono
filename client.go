package repono

func _encode(op, st, k, v []byte) []byte {
	if len(op) < 1 {
		return nil
	}
	b := make([]byte, 4+len(op)+len(st)+len(k)+len(v))
	for i := 0; i < len(b); i++ {
		switch {
		case i == 0:
			b[0] = byte(len(op))
		case i < len(op):
			b[i] = op[i-1]
		}
	}
	return b
}

/*package main

import "fmt"

func main() {
	//b := encode([]byte("set"), []byte("user"), []byte("273487248"), []byte(`{"name":"greg","age":29,"active":true}`))
	//b := encode([]byte("getall"), []byte("user"), nil, nil)
	b := encode([]byte("get"), []byte("user"), []byte("273487249"), nil)

	fmt.Printf("%s\n", b)
	fmt.Printf("% x\n", b)

}*/

func encode(op, st, k, v []byte) []byte {
	if len(op) < 1 {
		return nil
	}
	b := make([]byte, 5+len(op)+len(st)+len(k)+len(v))
	args := 0
	for i := 1; i < len(b); i++ {
		switch {
		case i == 1:
			b[i] = byte(len(op))
			args++
		case i < len(op)+2:
			b[i] = op[i-2]
		case i == len(op)+2 && st != nil:
			b[i] = byte(len(st))
			args++
		case i < len(st)+len(op)+3 && st != nil:
			b[i] = st[i-(len(op)+3)]
		case i == len(st)+len(op)+3 && k != nil:
			b[i] = byte(len(k))
			args++
		case i < len(k)+len(st)+len(op)+4 && k != nil:
			b[i] = k[i-(len(st)+len(op)+4)]
		case i < len(b)-1 && v != nil:
			args = 4
			b[i] = v[i-(len(k)+len(st)+len(op)+4)]
		}
	}
	b[len(b)-1] = byte('\n')

	b[0] = byte(args)
	return b
}
