package ec

import (
	"crypto/md5"
	"errors"
	"io"
	"math"
	"regexp"
	"strings"
)

const (
	v = 128 // 杂凑函数采用golang md5 128bit
)

/**
 * 密钥派生算法
 * @param {[type]} z    string [description]
 * @param {[type]} klen int)   (string,      error [description]
 */
func KDF(bit_x1x2 string, klen int) (string, error) {
	if klen > (2<<31-1)*v {
		return "", errors.New("invalid length")
	}
	ct := []byte{0, 0, 0, 1}
	hlen := int(math.Ceil(float64(klen) / v))
	var hashstring []string
	for i := 0; i < hlen; i++ {
		hashstring = append(hashstring, Bytes2Bits(H128(bit_x1x2+Bytes2Bits(ct))))
		ct[3] += 1
	}
	var tmp string
	var index int
	if 0 != klen%v {
		tmp = hashstring[hlen-1]
		index = klen - v*int(math.Floor(float64(klen)/v))
		tmp = tmp[:index]
		hashstring[hlen-1] = tmp
	}
	res := strings.Join(hashstring, "")
	re, err := regexp.Compile(`^0*$`)
	if err != nil {
		return "", errors.New("正则错误")
	}
	if re.MatchString(res) {
		return "", errors.New("生成的t全为0")
	}
	return res, nil
}

func H128(in string) []byte {
	t := md5.New()
	io.WriteString(t, in)
	return t.Sum(nil)
}
