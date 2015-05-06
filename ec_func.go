package ec

import (
	"io"
	"math/big"
)

//=======================================================
//	functions which implement the interface Curve
//=======================================================

/**
 * returns the parameters of the curve
 */
func (curve *CurveParams) Params() *CurveParams {
	return curve
}

/**
 * returns true if the given point(x, y) lies on the curve
 */
func (curve *CurveParams) IsOnCurve(p *CurvePoint) bool {

	x := p.X
	y := p.Y

	// y² = x³ - 3x + b  //curve equation
	y2 := new(big.Int).Mul(y, y)
	y2.Mod(y2, curve.P)

	x3 := new(big.Int).Exp(x, big.NewInt(3), curve.P)

	threeX := new(big.Int).Lsh(x, 1)
	threeX.Add(threeX, x)

	x3.Sub(x3, threeX)
	x3.Add(x3, curve.B)
	x3.Mod(x3, curve.P)

	return x3.Cmp(y2) == 0
}

/**
 * returns sum of two points
 */
func (curve *CurveParams) Add(p1, p2 *CurvePoint) *CurvePoint {
	z1 := zForAffine(p1)
	z2 := zForAffine(p2)
	return curve.affineFromJacobian(curve.addJacobian(p1.X, p1.Y, z1, p2.X, p2.Y, z2))
}

/**
 * returns the double value of the point
 */
func (curve *CurveParams) Double(p *CurvePoint) *CurvePoint {
	z := zForAffine(p)
	return curve.affineFromJacobian(curve.doubleJacobian(p.X, p.Y, z))
}

func (curve *CurveParams) ScalarMult(p *CurvePoint, k []byte) *CurvePoint {
	Bz := new(big.Int).SetInt64(1)
	x, y, z := new(big.Int), new(big.Int), new(big.Int)

	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			x, y, z = curve.doubleJacobian(x, y, z)
			if byte&0x80 == 0x80 {
				x, y, z = curve.addJacobian(p.X, p.Y, Bz, x, y, z)
			}
			byte <<= 1
		}
	}
	return curve.affineFromJacobian(x, y, z)
}

func (curve *CurveParams) ScalarBaseMult(k []byte) *CurvePoint {
	return curve.ScalarMult(curve.G, k)
}

//=======================================================
// transformations between Jacobian coordinates and
// Affine coordinates
//=======================================================
/**
 * zForAffine returns a Jacobian Z value for the affine point(x, y)
 * If x and y are zero, it assumes that the point represents the point
 * at infinity because (0, 0) is not on any of the curves handled here
 */
func zForAffine(p *CurvePoint) *big.Int {
	z := new(big.Int)
	if p.X.Sign() != 0 || p.Y.Sign() != 0 {
		z.SetInt64(1)
	}
	return z
}

/**
 * convet the point on affine coordinates to point on Jacobian coordinates
 * a given (x, y) position on the curve, the Jacobian coordinates are (x1, y1, z1)
 * where x = x1/z1² and y = y1/z1³
 */
func (curve *CurveParams) affineFromJacobian(x, y, z *big.Int) *CurvePoint {
	if 0 == z.Sign() {
		// it represent the point at infinity when z equals 0
		return new(CurvePoint)
	}

	// inverse of z
	zinv := new(big.Int).ModInverse(z, curve.P)
	// square of inverse of z
	zinvsq := new(big.Int).Mul(zinv, zinv)

	pOut := new(CurvePoint)
	pOut.X = new(big.Int).Mul(x, zinvsq)
	pOut.X.Mod(pOut.X, curve.P)

	zinvsq.Mul(zinvsq, zinv)
	pOut.Y = new(big.Int).Mul(y, zinvsq)
	pOut.Y.Mod(pOut.Y, curve.P)

	return pOut
}

/*
 * add Jacobian takes two points in Jacobian coordinates, (x1, y1, z1) and
 * (x2, y2, z2) and returns their sum, also in Jacobian form
 */
func (curve *CurveParams) addJacobian(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-2007-bl
	x3, y3, z3 := new(big.Int), new(big.Int), new(big.Int)
	if z1.Sign() == 0 {
		x3.Set(x2)
		y3.Set(y2)
		z3.Set(z2)
		return x3, y3, z3
	}
	if z2.Sign() == 0 {
		x3.Set(x1)
		y3.Set(y1)
		z3.Set(z1)
		return x3, y3, z3
	}

	z1z1 := new(big.Int).Mul(z1, z1)
	z1z1.Mod(z1z1, curve.P)
	z2z2 := new(big.Int).Mul(z2, z2)
	z2z2.Mod(z2z2, curve.P)

	u1 := new(big.Int).Mul(x1, z2z2)
	u1.Mod(u1, curve.P)
	u2 := new(big.Int).Mul(x2, z1z1)
	u2.Mod(u2, curve.P)
	h := new(big.Int).Sub(u2, u1)
	xEqual := h.Sign() == 0
	if h.Sign() == -1 {
		h.Add(h, curve.P)
	}
	i := new(big.Int).Lsh(h, 1)
	i.Mul(i, i)
	j := new(big.Int).Mul(h, i)

	s1 := new(big.Int).Mul(y1, z2)
	s1.Mul(s1, z2z2)
	s1.Mod(s1, curve.P)
	s2 := new(big.Int).Mul(y2, z1)
	s2.Mul(s2, z1z1)
	s2.Mod(s2, curve.P)
	r := new(big.Int).Sub(s2, s1)
	if r.Sign() == -1 {
		r.Add(r, curve.P)
	}
	yEqual := r.Sign() == 0
	if xEqual && yEqual {
		return curve.doubleJacobian(x1, y1, z1)
	}
	r.Lsh(r, 1)
	v := new(big.Int).Mul(u1, i)

	x3.Set(r)
	x3.Mul(x3, x3)
	x3.Sub(x3, j)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	x3.Mod(x3, curve.P)

	y3.Set(r)
	v.Sub(v, x3)
	y3.Mul(y3, v)
	s1.Mul(s1, j)
	s1.Lsh(s1, 1)
	y3.Sub(y3, s1)
	y3.Mod(y3, curve.P)

	z3.Add(z1, z2)
	z3.Mul(z3, z3)
	z3.Sub(z3, z1z1)
	z3.Sub(z3, z2z2)
	z3.Mul(z3, h)
	z3.Mod(z3, curve.P)

	return x3, y3, z3
}

//doubleJacobian takes a point in Jacobian coordinates (x, y, z) and
//returns its double, also in Jacobian coordinates
func (curve *CurveParams) doubleJacobian(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	delta := new(big.Int).Mul(z, z)
	delta.Mod(delta, curve.P)
	gamma := new(big.Int).Mul(y, y)
	gamma.Mod(gamma, curve.P)
	alpha := new(big.Int).Sub(x, delta)
	if alpha.Sign() == -1 {
		alpha.Add(alpha, curve.P)
	}
	alpha2 := new(big.Int).Add(x, delta)
	alpha.Mul(alpha, alpha2)
	alpha2.Set(alpha)
	alpha.Lsh(alpha, 1)
	alpha.Add(alpha, alpha2)

	beta := alpha2.Mul(x, gamma)

	x3 := new(big.Int).Mul(alpha, alpha)
	beta8 := new(big.Int).Lsh(beta, 3)
	x3.Sub(x3, beta8)
	for x3.Sign() == -1 {
		x3.Add(x3, curve.P)
	}
	x3.Mod(x3, curve.P)

	z3 := new(big.Int).Add(y, z)
	z3.Mul(z3, z3)
	z3.Sub(z3, gamma)
	if -1 == z3.Sign() {
		z3.Add(z3, curve.P)
	}
	z3.Sub(z3, delta)
	if -1 == z3.Sign() {
		z3.Add(z3, curve.P)
	}
	z3.Mod(z3, curve.P)

	beta.Lsh(beta, 2)
	beta.Sub(beta, x3)
	if -1 == beta.Sign() {
		beta.Add(beta, curve.P)
	}
	y3 := alpha.Mul(alpha, beta)

	gamma.Mul(gamma, gamma)
	gamma.Lsh(gamma, 3)
	gamma.Mod(gamma, curve.P)

	y3.Sub(y3, gamma)
	if -1 == y3.Sign() {
		y3.Add(y3, curve.P)
	}
	y3.Mod(y3, curve.P)

	return x3, y3, z3

}

var mask = []byte{0xff, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f}

// GenerateKey returns a public/private key pair. The private key is
// generated using the given reader, which must return random data.
func GenerateKey(curve Curve, rand io.Reader) (priv []byte, pubk *CurvePoint, err error) {
	bitSize := curve.Params().BitSize
	byteLen := (bitSize + 7) >> 3
	priv = make([]byte, byteLen)
	pubk = new(CurvePoint)

	for pubk.X == nil {
		_, err = io.ReadFull(rand, priv)
		if err != nil {
			return
		}
		// We have to mask off any excess bits in the case that the size of the
		// underlying field is not a whole number of bytes.
		priv[0] &= mask[bitSize%8]
		// This is because, in tests, rand will return all zeros and we don't
		// want to get the point at infinity and loop forever.
		priv[1] ^= 0x42
		pubk = curve.ScalarBaseMult(priv)
	}
	return
}

// Mashal converts a point into the form specialized in section 4.3.6 of ASNI X9.62.
func Marshal(curve Curve, p *CurvePoint) []byte {
	byteLen := (curve.Params().BitSize + 7) >> 3

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // unpressed point

	xBytes := p.X.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	yBytes := p.Y.Bytes()
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)

	return ret
}

// Unmashal converts a point, serialized by Mashal, into an x, y pair, On error, x = nil
func Unmarshal(curve Curve, data []byte) (ret *CurvePoint) {
	byteLen := (curve.Params().BitSize + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}
	if data[0] != 4 { // unpressed point
		return
	}
	ret = new(CurvePoint)
	ret.X = new(big.Int).SetBytes(data[1 : 1+byteLen])
	ret.Y = new(big.Int).SetBytes(data[1+byteLen:])
	return
}
