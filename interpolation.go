// References
// - https://neuron.eng.wayne.edu/auth/ece3040/lectures/lecture14.pdf
// - http://mathsfromnothing.epizy.com/successive-parabolic-interpolation/
//
package main

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"
)

func main() {
	example1()
	// example2()
	// example3()
	// example4()
	// example5()
	// example6()
}

// knownSet은 이차다항함수 근사를 위한 세 개의 점 집합을 나타냅니다.
type knownSet struct {
	x1, y1 float64
	x2, y2 float64
	x3, y3 float64
}

// interpolate는 점 집합 w에 대해 보간한 이차곡선의 극값을 반환합니다. 보간된
// 이차곡선의 최고차항 계수는 양수거나 음수일 수 있습니다. 만약 주어진 세
// 점으로부터 이차곡선을 보간하지 못한다면 실수가 아닌 값을 반환할 것입니다.
// 추가적인 내용은 isfinite 함수를 참고하세요.
func interpolate(w *knownSet) float64 {
	z1 := w.x1 * (w.y2 - w.y3)
	z2 := w.x2 * (w.y3 - w.y1)
	z3 := w.x3 * (w.y1 - w.y2)
	return 0.5 * (w.x1*z1 + w.x2*z2 + w.x3*z3) / (z1 + z2 + z3)
}

// squeeze는 점 집합 w와 이차곡선의 극값 x가 주어졌을 때 successive parabolic
// interpolation을 수행하기 위해 점 집합을 조정합니다. 그리고 squeeze는 점 집합
// w의 어떤 원소가 x인지 알려주는 포인터와 그에 대응하는 y 포인터를 반환합니다.
// 점 집합은 반드시 정렬되어 있어야 합니다. 즉 w.x1 <= w.x2 <= w.x3 여야
// 합니다. 또한 반드시 x는 실수여야 합니다.
func squeeze(w *knownSet, x float64) (*float64, *float64) {
	nan := math.NaN()
	if x <= w.x1 {
		w.x1, w.x2, w.x3 = x, w.x1, w.x2
		w.y1, w.y2, w.y3 = nan, w.y1, w.y2
		return &w.x1, &w.y1
	} else if x <= w.x2 {
		w.x2, w.x3 = x, w.x2
		w.y2, w.y3 = nan, w.y2
		return &w.x2, &w.y2
	} else if x <= w.x3 {
		w.x1, w.x2 = w.x2, x
		w.y1, w.y2 = w.y2, nan
		return &w.x2, &w.y2
	} else {
		w.x1, w.x2, w.x3 = w.x2, w.x3, x
		w.y1, w.y2, w.y3 = w.y2, w.y3, nan
		return &w.x3, &w.y3
	}
}

// isfinite 함수는 주어진 x가 실수인지 확인합니다.  여기서 실수란 초실수계가
// 아닙니다.  참고로 양의 무한, 음의 무한, NaN은 실수가 아닙니다.
func isfinite(x float64) bool {
	return !math.IsInf(x, 0) && !math.IsNaN(x)
}

func ord3(x, y, z float64) (float64, float64, float64) {
	if x > y {
		x, y = y, z
	}
	if y > z {
		y, z = z, y
	}
	if x > y {
		x, y = y, z
	}
	return x, y, z
}

func iterativelyInterpolate(w *knownSet, z *float64) bool {
	*z = interpolate(w)
	return isfinite(*z)
}

func PrintApproximation(fn func(float64) float64, x1, x2, x3 float64) {
	o := tabwriter.NewWriter(os.Stdout, 8, 4, 1, ' ', 0)
	x1, x2, x3 = ord3(x1, x2, x3)
	w := knownSet{
		x1: x1, y1: fn(x1),
		x2: x2, y2: fn(x2),
		x3: x3, y3: fn(x3),
	}
	fmt.Fprintf(o, "#\tx\ty\tz\t\n")
	fmt.Fprintf(o, "%3d\t%.9f\t%.9f\t%.9f\t\n", 0, w.x1, w.x2, w.x3)
	var z float64
	for i := 1; i <= 100 && iterativelyInterpolate(&w, &z); i++ {
		xp, yp := squeeze(&w, z)
		*yp = fn(*xp)
		fmt.Fprintf(o, "%3d\t%.9f\t%.9f\t%.9f\t\n", i, w.x1, w.x2, w.x3)
	}
	o.Flush()
}

func F1(x float64) float64 {
	return x*x/10 - 2*math.Sin(x)
}

func F2(x float64) float64 {
	return math.Sinh(math.Sin(x))
}

func F3(x float64) float64 {
	return math.Exp(-math.Sin(x))
}

func F4(x float64) float64 {
	cos := -math.Cos(x)
	return math.Exp(-math.Sin(x)) + math.Exp(-x*cos*cos)
}

func F5(x float64) float64 {
	return math.Sin(1 / x)
}

func F6(x float64) float64 {
	const a = 0.3
	const b = 7
	z := 0.0
	for i := 0; i <= 20; i++ {
		z = math.FMA(math.Pow(a, float64(i)), math.Cos(math.Pow(b, float64(i))*math.Pi*x), z)
	}
	return z
}

func example1() {
	fmt.Println("x^2/10 - 2sin(x)")
	PrintApproximation(F1, 0, 1, 4)
	PrintApproximation(F1, 4, 5, 6)
	PrintApproximation(F1, 4, 6, 8)
	PrintApproximation(F1, 16, 20, 22)
}

func example2() {
	fmt.Println("sinh(sin(x))")
	PrintApproximation(F2, 0, 1, 2)
	PrintApproximation(F2, 3, 4, 5)
	PrintApproximation(F2, 2, 3, 4)
}

func example3() {
	fmt.Println("exp(sin(-x))")
	PrintApproximation(F3, 1, 2, 3)
	PrintApproximation(F3, 2, 3, 4)
	PrintApproximation(F3, 4, 5, 6)
}

func example4() {
	fmt.Println("exp(sin(-x)) + exp(-x * cos^2(x))")
	PrintApproximation(F4, 1, 2, 3)
	PrintApproximation(F4, 2, 3, 4)
	PrintApproximation(F4, -2, -1.2, -0.8)
	PrintApproximation(F4, 4, 5, 6)
}

func example5() {
	fmt.Println("sin(1 / x)")
	PrintApproximation(F5, 0.1, 0.2, 0.3)
	PrintApproximation(F5, 0.1, 0.3, 0.5)
	PrintApproximation(F5, 0.001, 0.002, 0.003)
	PrintApproximation(F5, 0.00001, 0.00002, 0.00003)
	PrintApproximation(F5, 1, 2, 3)
}

// https://en.wikipedia.org/wiki/Weierstrass_function
func example6() {
	fmt.Println("Weierstrass function (a = 0.3, b = 7)")
	PrintApproximation(F6, 0.5, 1, 1.5)
	PrintApproximation(F6, 0.5, 1, 1.6)
}
