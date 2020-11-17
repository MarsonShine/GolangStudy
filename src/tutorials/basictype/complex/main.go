package main

import "fmt"

// 复数：我们把形如z=a+bi（a,b均为实数）的数称为复数，其中a称为实部，b称为虚部，i称为虚数单位。当z的虚部等于零时，常称z为实数；当z的虚部不等于零时，实部等于零时，常称z为纯虚数。复数域是实数域的代数闭包，即任何复系数多项式在复数域中总有根。
// 内建的real和imag函数分别返回复数的实部和虚部
// 复数运算法则
// 设z1=a+bi，z2=c+di是任意两个复数， z1+z2=(a+c)+(b+d) z1*z2=(a*c-b*d)+(bc+ad)i
func main() {
	var x complex128 = complex(1, 2) // 1+2i
	var y complex128 = complex(3, 4) // 3+4i
	fmt.Println(x * y)               // "(-5+10i)"
	fmt.Println(real(x * y))         // "-5" 实部
	fmt.Println(imag(x * y))         // "10" 虚部
}
