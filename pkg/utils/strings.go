package utils

import (
	"github.com/jamesits/go-bytebuilder"
	"unicode/utf16"
)

func UTF16ToStringSlice(s []uint16) (ret []string) {
	if len(s) < 2 {
		return
	}

	head := 0

	for i, v := range s {
		if v == 0 {
			ret = append(ret, string(utf16.Decode(s[head:i])))

			if s[i+1] == 0 {
				break
			} else {
				head = i + 1
			}
		}
	}

	return
}

func UTF16ByteArrayToStringSlice(s []byte) []string {
	return UTF16ToStringSlice(bytebuilder.SliceCast[byte, uint16](s))
}
