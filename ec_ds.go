package ec

import (
	"math/big"
	"strings"
)

// 4.2.1 integer to byte string
func Int2Bytes(in *big.Int) []byte {
	return in.Bytes()
}

// 4.2.2 byte string to integer
func Bytes2Int(in []byte) *big.Int {
	return new(big.Int).SetBytes(in)
}

// 4.2.3 bit string 2 byte string
func Bits2Bytes(in string) []byte {
	res := new(big.Int)
	res.SetString(in, 2)
	ret := res.Bytes()
	return ret
}

// 4.2.4 byte string to bit string
func Bytes2Bits(in []byte) string {
	len_b := len(in)
	res := make([]rune, len_b*8)
	i := 0
	for _, item := range in {
		res[i] = rune(item)&0x80>>7 + 0x30
		res[i+1] = rune(item)&0x40>>6 + 0x30
		res[i+2] = rune(item)&0x20>>5 + 0x30
		res[i+3] = rune(item)&0x10>>4 + 0x30
		res[i+4] = rune(item)&0x08>>3 + 0x30
		res[i+5] = rune(item)&0x04>>2 + 0x30
		res[i+6] = rune(item)&0x02>>1 + 0x30
		res[i+7] = rune(item)&0x01 + 0x30
		i += 8
	}
	return string(res)
}

// 字符串到bit串
func String2Bits(in string) string {
	buf := []byte(in)
	return Bytes2Bits(buf)
}

// bit串到字符串
func Bits2String(in string) string {
	buf := new(big.Int)
	buf.SetString(in, 2)
	res := buf.Bytes()
	return string(res)
}

// 字符串到字节串
func String2Bytes(in string) []byte {
	return []byte(in)
}

// 字节串到字符串
func Bytes2String(in []byte) string {
	return string(in)
}

// 字符串转换成16进制串
func String2Hex(in string) string {
	char := "0123456789ABCDEF"
	str := String2Bytes(in)
	var sb []string

	var bit uint8
	for i := 0; i < len(str); i++ {
		bit = (str[i] & 0xf0) >> 4
		sb = append(sb, string(char[bit]))
		bit = str[i] & 0x0f
		sb = append(sb, string(char[bit]))
	}
	return strings.Join(sb, "")
}

func Hex2String(in string) string {
	char := "0123456789ABCDEF"
	in = strings.Replace(in, " ", "", -1)
	str := []byte(in)
	var bytes []byte
	var bit int

	for i := 0; i < len(in)/2; i++ {
		bit = strings.IndexByte(char, str[2*i]) << 4
		bit += strings.IndexByte(char, str[2*i+1])
		bytes = append(bytes, byte(bit&0xff))
	}
	return Bytes2String(bytes)
}

func GetMessageBigs(msg string) (bits string, klen int) {
	ret := String2Bits(msg)
	return ret, len(ret)
}

// 4.2.5 field element to byte string
// func FieldElement2Bytes(a *big.Int, curve *CurveParams) string {
// 	a.Mod(a, curve.P)
// 	return Int2Bytes(a)
// }

// 4.2.6 byte string 2 field element
// func Bytes2FiledElement(in string) *big.Int {
// 	bytes :=
// 	return new(big.Int).SetBytes(in)
// }
