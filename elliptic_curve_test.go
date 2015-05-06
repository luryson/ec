package ec

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestOnGeneratekey(t *testing.T) {
	// p521 := P384()
	p521 := P521()
	if !p521.IsOnCurve(p521.Params().G) {
		t.Errorf("Fail")
	}
	_, pubk, err := GenerateKey(p521, rand.Reader)
	if err != nil {
		t.Errorf("Generate Fail")
	}
	if !p521.IsOnCurve(pubk) {
		t.Errorf("pubk not on curve!")
	}

	if !p521.IsOnCurve(p521.Double(pubk)) {
		t.Error("Double point not on curve")
	}

	if !p521.IsOnCurve(p521.Add(p521.Params().G, pubk)) {
		t.Error("Add point not on curve")
	}

	if !p521.IsOnCurve(p521.ScalarMult(p521.Params().G, []byte{1, 3, 4})) {
		t.Error("Scalar mul point not on curve")
	}
}

func TestOnKDF(t *testing.T) {
	// 加密明文
	msg := "xxx"
	msg = Hex2String(String2Hex(msg))
	t.Errorf("%v\n", msg)

	// 椭圆曲线参数
	p521 := P521()
	priv, pubk, err := GenerateKey(p521, rand.Reader)
	if err != nil {
		t.Error("fail to generate key of B")
	}
	t.Errorf("B's priv key is :\n%v\n", priv)
	t.Errorf("B's pub key is :\n%v, \n%v\n", pubk.X, pubk.Y)

	k, _, err1 := GenerateKey(p521, rand.Reader)
	if err1 != nil {
		t.Error("fail to generate rand k")
	}
	t.Errorf("rand k is :\n%v\n", k)

	c1 := Bytes2Bits(Marshal(p521, p521.ScalarBaseMult(k)))
	t.Errorf("c1 is :\n%v\n", c1)

	tmp := p521.ScalarMult(p521.Params().G, k)
	t.Errorf("x2, y2 is :\n%v\n%v\n", tmp.X, tmp.Y)

	x2 := String2Hex(tmp.X.String())
	y2 := String2Hex(tmp.Y.String())
	t1, _ := KDF(x2+y2, len(String2Bits(msg)))
	t.Errorf("t is %v\n", t1)

	t.Errorf("%v\n", String2Bits(msg))
	m, _ := new(big.Int).SetString(String2Hex(msg), 16)
	t2, _ := new(big.Int).SetString(t1, 2)

	t.Errorf("%v, %v\n", m, t2)
	c2 := new(big.Int).And(m, t2)
	t.Errorf("c2 is :\n%v\n", c2)

	c3 := H128(x2 + m.String() + y2)
	t.Errorf("c3 is :\n%v\n", String2Hex(Bytes2String(c3)))

	cipher := String2Hex(Bytes2String(Bits2Bytes(c1))) + String2Hex(Bytes2String(Int2Bytes(c2))) + String2Hex(Bytes2String(c3))
	t.Errorf("%v\n", cipher)

	// decrypt
	cipherReceived, _ := new(big.Int).SetString(cipher, 16)
	b_cipherReceived := cipherReceived.Bytes()
	t.Errorf("%v\n", b_cipherReceived)
	l_c1 := (p521.Params().BitSize + 7) >> 3
	t.Errorf("%v\n", l_c1)

	c1_received := b_cipherReceived[:1+(2*l_c1)]
	t.Errorf("%v\n%v\n", Unmarshal(p521, c1_received).X, Unmarshal(p521, c1_received).Y)

	if !p521.IsOnCurve(Unmarshal(p521, c1_received)) {
		t.Errorf("P is not on curver")
	}

	p2 := p521.ScalarMult(Unmarshal(p521, c1_received), priv)
	x2 = Bytes2Bits(Int2Bytes(p2.X))
	y2 = Bytes2Bits(Int2Bytes(p2.Y))
	t.Errorf("x2:%v\ny2:%v\n", x2, y2)
	klen := len(b_cipherReceived) - (1 + (2 * l_c1)) - (128 >> 3)
	t.Errorf("%v\n", klen)
	t3, _ := KDF(x2+y2, klen<<3)
	t.Errorf("t:%v\n", t3)

	c2b := b_cipherReceived[1+(2*l_c1) : 1+(2*l_c1)+klen]
	t.Errorf("%v\n", c2b)

	t3b, _ := new(big.Int).SetString(t3, 2)
	c3t := new(big.Int).SetBytes(c2b)
	res_int := new(big.Int).And(t3b, c3t)
	t.Errorf("%v\n", res_int)
	t.Errorf("%v\n", Bytes2String(res_int.Bytes()))
}

// func TestOnMashal(t *testing.T) {
// 	p521 := P384()
// 	priv, pubk, err := GenerateKey(p521, rand.Reader)
// 	if err != nil {
// 		t.Error("fail to generate key of B")
// 	}
// 	t.Errorf("B's priv key is :\n%v\n", priv)
// 	t.Errorf("B's pub key is :\n%v, \n%v\n", pubk.x, pubk.y)

// 	pre := Marshal(p521, pubk)
// 	t.Errorf("%v\n", pre)

// 	t.Errorf("%v\n", (p521.Params().BitSize+7)>>3)
// 	t.Errorf("%v\n", len(pre))

// 	p := Unmarshal(p521, pre)
// 	t.Errorf("%v\n%v\n", p.x, p.y)
// }
