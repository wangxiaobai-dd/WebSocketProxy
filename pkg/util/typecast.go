package util

import "strconv"

func StrToUint32(s string) (uint32, error) {
	u64, err := strconv.ParseUint(s, 10, 32)
	return uint32(u64), err
}

func Uint32ToStr(u uint32) string {
	return strconv.FormatInt(int64(u), 10)
}
