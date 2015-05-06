package ec

import (
	"math/big"
)

// interface
// A curve represents a short-form Weierstrass curve with a = -3
type Curve interface {
	// 返回椭圆曲线的参数
	Params() *CurveParams
	// 判断一个点是否在椭圆曲线上
	IsOnCurve(p *CurvePoint) bool
	// 椭圆曲线上的加法
	Add(p1, p2 *CurvePoint) *CurvePoint
	// 椭圆曲线上的2倍点运算
	Double(p *CurvePoint) *CurvePoint
	// 椭圆曲线上的倍点运算
	ScalarMult(p *CurvePoint, k []byte) *CurvePoint
	// 椭圆曲线的基点的倍点运算
	ScalarBaseMult(k []byte) *CurvePoint
}

// struct of CurveParams
type CurveParams struct {
	P       *big.Int    // 基域
	N       *big.Int    // 椭圆曲线的阶
	B       *big.Int    // 椭圆曲线方程的常数B y² = x³ - 3x + b  //curve equation
	G       *CurvePoint // 椭圆曲线的基点
	BitSize int         //
}

// struct of CurvePoint
type CurvePoint struct {
	X *big.Int
	Y *big.Int
}
