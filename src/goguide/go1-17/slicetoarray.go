package go117

import "fmt"

func sliceToArray() {
	s := make([]byte, 2, 4)
	s0 := (*[0]byte)(s)     //go1.17以下会报 cannot convert s (variable of type []byte) to *[0]byte
	s1 := (*[1]byte)(s[1:]) // &s1[0] == &s[1]
	s2 := (*[2]byte)(s)     // &s2[0] == &s[0]
	_ = (*[4]byte)(s)       // panics: len([4]byte) > len(s)

	var t []string
	_ = (*[0]string)(t) // t0 == nil
	_ = (*[1]string)(t) // panics: len([1]string) > len(t)

	u := make([]byte, 0)
	_ = (*[0]byte)(u) // u0 != nil
	fmt.Println(s0, s1, s2)
}
