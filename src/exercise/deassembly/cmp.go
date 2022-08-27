package main

/*
原文节选：https://go.dev/doc/asm

关于汇编生成的CPU指令，可以从https://github.com/golang/go/files/447163/GoFunctionsInAssembly.pdf，也可以从这边文章快速了解指令的含义https://segmentfault.com/a/1190000039978109

汇编器以半抽象的形式工作，因此当看到像MOV这样的指令时，工具链实际上为该操作生成的可能根本不是一个move指令，而是一个clear或load指令。
FUNCDATA和PCDATA指令包含垃圾回收器使用的信息，这两个指令是编译器引入的

##符号

有四个预先声明的符号引用“伪寄存器”。这些寄存器不是真正的寄存器，而是由工具链维护的虚拟寄存器，比如帧指针。伪寄存器的集合对于所有架构都是相同的:
FP: Frame pointer：参数和本地变量
PC: 程序计数器，跳转和分支
SB: 基地地址：全局符号
SP: Stack pointer：本地堆栈中的最高地址
所有用户定义的符号都作为偏移量写入伪寄存器FP(参数和局部变量)和SB(全局变量)。

SB伪寄存器可以被认为是内存的源头，因此符号`foo(SB)`是名称foo在内存中的地址。此表单用于命名全局函数和数据。将<>添加到名称中，如在foo<>(SB)中，使名称仅在当前源文件中可见，就像C文件中的顶级静态声明一样。向名称添加偏移量指的是该符号地址的偏移量，因此foo+4(SB)在foo的开始位置之后的4个字节。

FP伪寄存器是一个用于引用函数参数的虚拟帧指针。编译器维护一个虚拟帧指针，并将堆栈上的参数作为与伪寄存器的偏移量来引用。因此，0(FP)是函数的第一个参数，8(FP)是第二个参数(在64位机器上)，依此类推。但是，当以这种方式引用函数参数时，有必要在开头放置一个名称，比如first_arg+0(FP)和second_arg+8(FP)。

对于带有Go原型的汇编函数，Go vet 将检查参数名称和偏移量是否匹配。在32位系统上，通过在名称中添加_lo或_hi后缀来区分64位值的高32位和低32位，如arg_lo+0(FP)或arg_hi+4(FP)。如果Go原型没有命名其结果，则预期的程序集名称为 ret。

SP伪寄存器是一个虚拟堆栈指针，用于引用框架局部变量(frame-locals)和为函数调用准备的参数。它指向本地堆栈帧内的最高地址，因此引用应该在[−framesize, 0): x-8(SP)， y-4(SP)等范围内使用负偏移量。

在具有名为SP的硬件寄存器的体系结构上，名称前缀将对虚拟堆栈指针的引用与对体系结构SP寄存器的引用区分开来。也就是说，x-8(SP)和-8(SP)是不同的内存位置:第一个指向虚拟堆栈指针伪寄存器，而第二个指向硬件的SP寄存器。

在机器上，SP和PC通常是物理的、编号的寄存器的别名，在Go汇编程序中，SP和PC的名称仍然被特殊处理；例如，对SP的引用需要一个符号，就像FP一样。要访问实际的硬件寄存器，请使用真正的R名称。例如，在ARM架构下，硬件SP和PC可以通过R13和R15访问。

分支和直接跳转总是被写入到PC的偏移量，或者作为到标签的跳转:
```
label:
	MOVW $0, R1
	JMP label
```
每个标签仅在定义它的函数中可见。因此，在一个文件中允许多个函数去定义和使用相同的标签名称。直接跳转和调用指令可以指向文本符号，如name(SB)，但不能指向符号的偏移，如name+4(SB)。

指令、寄存器和汇编指令总是大写。(例外:g寄存器在ARM上重命名。)

在Go目标文件和二进制文件中，符号的全名是包路径后面加.号和符号名:`fmt.Printf`或`math/rand.Int`。因为汇编器的解析器将`.`和斜杠视为标点符号，所以这些字符串不能直接用作标识符名称。相反，汇编器允许在标识符中使用中间的.符号U+00B7和除法斜杠U+2215，并将它们重写为普通的句号和斜杠。在汇编器源文件中，上面的符号被写成`fmt·Printf`和`math∕rand·Int`。当使用`-S`选项时，编译器生成的程序集清单会直接显示.和斜杠，而不是汇编器所需的Unicode替换。

##指令
汇编程序使用各种指令将文本和数据绑定到符号名称。例如，下面是一个简单的完整函数定义。TEXT指令声明了符号runtime·profileloop，随后的指令构成了函数体。TEXT块中的最后一条指令必须是某种跳转，通常是RET(伪)指令。(如果不是，链接器将附加一个跳转到自身的指令;在文本中没有失败。)在符号之后，参数是标志(见下文)和帧大小，一个常量(见下文):
```
TEXT runtime·profileloop(SB),NOSPLIT,$8
	MOVQ	$runtime·profileloop1(SB), CX
	MOVQ	CX, 0(SP)
	CALL	runtime·externalthreadhandler(SB)
	RET
```
在一般情况下，帧大小后面跟着参数大小，用减号分隔。(这不是减法，只是特殊的语法。)帧大小`$24-8`表示函数有一个24字节的帧，并且在调用时带有8个字节的参数，这些参数位于调用者的帧上。如果TEXT没有指定NOSPLIT，则必须提供参数size。对于带有Go原型的汇编函数，Go vet将检查参数大小是否正确。

请注意，符号名称使用中间的点来分隔组件，并被指定为与静态基伪寄存器SB的偏移量。该函数将在包运行时的Go源代码中使用简单名称profileloop调用。

全局数据符号是由一系列初始化`DATA`指令后跟一个`GLOBL`指令定义的。每个DATA指令初始化相应内存的一段。未显式初始化的内存为零。DATA指令的一般形式是:
```
DATA	symbol+offset(SB)/width, value
```
它用给定的值在给定的偏移量和宽度处初始化符号内存。给定符号的DATA指令必须以递增的偏移量写入。

GLOBL指令将一个符号声明为全局的。参数是可选的标志，数据的大小声明为一个全局，它的初始值为零，除非DATA指令初始化了它。GLOBL指令必须遵循任何相应的DATA指令。

举个例子
```
DATA divtab<>+0x00(SB)/4, $0xf4f8fcff
DATA divtab<>+0x04(SB)/4, $0xe6eaedf0
...
DATA divtab<>+0x3c(SB)/4, $0x81828384
GLOBL divtab<>(SB), RODATA, $64

GLOBL runtime·tlsoffset(SB), NOPTR, $4
```
声明并初始化divtab<>，一个 4 字节整数值的只读 64 字节表，并声明 runtime·tlsoffset，一个不包含指针的 4 字节隐式归零变量。
指令可能有一个或两个参数。 如果有两个，第一个是标志位掩码，可以写成数字表达式，加或或运算在一起，或者可以符号设置，以便于人类吸收。 它们在标准#include 文件 textflag.h 中定义的值是：

- NOPROF = 1
  (对于`TEXT` 项)不要分析被标记的函数。此标志已弃用。
- DUPOK = 2
  在一个二进制文件中有该符号的多个实例是合法的。链接器将选择一个副本来使用。
- NOSPLIT = 4
  (对于`TEXT` 项)不要插入序言来检查堆栈是否必须拆分。例程的框架，加上它所调用的任何东西，都必须能够容纳在当前堆栈段中剩余的空闲空间中。用于保护例程，如堆栈拆分代码本身。
- RODATA = 8
  (对于DATA和GLOBL项)将这些数据放在只读区域。
- NOPTR = 16
  (对于DATA和GLOBL项)该数据不包含指针，因此不需要被垃圾收集器扫描。
- WRAPPER = 32
  (对于DATA和GLOBL项)这是一个包装器函数，不应该算作禁用恢复。
- NEEDCTXT = 64
  (对于DATA和GLOBL项)这个函数是一个闭包，因此它使用传入的上下文寄存器。
- LOCAL = 128
  这个符号是动态共享对象本地的。
- TLSBSS = 256
  (对于DATA和GLOBL项)将这些数据放入线程本地存储中。
- NOFRAME = 512
  (对于`TEXT` 项) 不要插入指令来分配堆栈帧和保存/恢复返回地址，即使这不是一个叶函数。仅对声明帧大小为0的函数有效。
- TOPFRAME = 2048
  (对于`TEXT` 项) 函数是调用堆栈的最外层框架。回溯应该在此函数处停止。

## 与Go类型和常量交互
如果一个包有任何.s文件，那么go build将指示编译器发出一个名为go_asm.h的特殊头文件，然后.s文件可以#include。该文件包含符号#define常量，用于表示Go struct字段的偏移量、Go struct类型的大小，以及当前包中定义的大多数Go const声明。Go程序集应该避免对Go类型的布局做假设，而是使用这些常量。这提高了汇编代码的可读性，并保持它对Go类型定义或Go编译器使用的布局规则中的数据布局更改的健壮性。

常量的形式为const_name。例如，给定Go声明const bufSize = 1024，汇编代码可以将该常量的值引用为const_bufSize。

字段偏移量的形式为type_field。结构体大小的形式为type__size。例如，考虑下面的Go定义:
```
type reader struct {
	buf [bufSize]byte
	r   int
}
```
程序集可以将该结构体的大小称为reader__size，两个字段的偏移量称为reader_buf和reader_r。因此，如果寄存器R1包含一个指向reader的指针，程序集可以将r字段引用为reader_r(R1)。

如果这些#define名称有任何歧义(例如，带有_size字段的结构体)，#include "go_asm.h"将失败，并报出"redefinition of macro"错误。

## 协调运行时
为了让垃圾收集正确运行，运行时必须知道指针在所有全局数据和大多数堆栈帧中的位置。Go编译器在编译Go源文件时发出此信息，但是汇编程序必须明确定义它。

用NOPTR标记(见上文)的数据符号被视为不包含指向运行时分配数据的指针。带有 RODATA 标志的数据符号被分配在只读存储器中，因此被视为隐式标记的 NOPTR。总大小小于指针的数据符号也被视为隐式标记的 NOPTR。 无法在汇编源文件中定义包含指针的符号；这样的符号必须在 Go 源文件中定义。 即使没有 DATA 和 GLOBL 指令，汇编源代码仍然可以按名称引用符号。一般经验法则是在 Go 中而不是在汇编中定义所有非 RODATA 符号。

每个函数还需要注释，在其参数、结果和本地堆栈帧中给出实时指针的位置。对于没有指针结果、本地堆栈帧或没有函数调用的汇编函数，唯一的要求是在同一个包的 Go 源文件中为该函数定义一个 Go 原型。汇编函数的名称不能包含包名组件（例如，包 syscall 中的函数 Syscall 在其 TEXT 指令中应使用名称·Syscall，而不是等效名称 syscall·Syscall）。对于更复杂的情况，需要显式注释。这些注释使用标准#include funcdata.h 文件中定义的伪指令。

如果一个函数没有参数也没有结果，指针信息可以省略。这是通过TEXT指令上`$n-0`的参数大小注释表示的。否则，指针信息必须由Go源文件中的函数的Go原型提供，即使不是直接从Go调用的汇编函数也是如此。（原型也会让 go vet 检查参数引用。）在函数开始时，假定参数已初始化，但假定结果未初始化。如果结果将在调用指令期间保持实时指针，则该函数应首先将结果归零，然后执行伪指令 GO_RESULTS_INITIALIZED。该指令记录了结果现在已经初始化，并且应该在堆栈移动和垃圾回收期间进行扫描。通常更容易安排汇编函数不返回指针或不包含调用指令；标准库中没有汇编函数使用 GO_RESULTS_INITIALIZED。

如果函数没有本地堆栈帧，则可以省略指针信息。这由 TEXT 指令上 $0-n 的局部帧大小注释表示的。如果函数不包含调用指令，指针信息也可以省略。否则，本地堆栈帧不得包含指针，程序集必须通过执行伪指令 NO_LOCAL_POINTERS 来确认这一事实。因为堆栈大小调整是通过移动堆栈来实现的，所以堆栈指针可能会在任何函数调用期间发生变化：即使指向堆栈数据的指针也不能保存在局部变量中。

汇编函数应该总是被赋予 Go 原型，既可以为参数和结果提供指针信息，也可以让 `go vet` 检查用于访问它们的偏移量是否正确。

## 不支持的指令操作
汇编器被设计成支持编译器，所以并不是所有的硬件指令都是为所有架构定义的:如果编译器不生成它，它可能就不存在。如果您需要使用缺失的指令，有两种方法来继续。一种是更新汇编器以支持该指令，这很简单，但只有在该指令可能会再次使用时才值得。相反，对于简单的一次性情况，可以使用BYTE和WORD指令将显式数据放置到TEXT中的指令流中。下面是386运行时如何定义64位原子加载函数。
```
// uint64 atomicload64(uint64 volatile* addr);
// so actually
// void atomicload64(uint64 *res, uint64 volatile *addr);
TEXT runtime·atomicload64(SB), NOSPLIT, $0-12
	MOVL	ptr+0(FP), AX
	TESTL	$7, AX
	JZ	2(PC)
	MOVL	0, AX // crash with nil ptr deref
	LEAL	ret_lo+4(FP), BX
	// MOVQ (%EAX), %MM0
	BYTE $0x0f; BYTE $0x6f; BYTE $0x00
	// MOVQ %MM0, 0(%EBX)
	BYTE $0x0f; BYTE $0x7f; BYTE $0x03
	// EMMS
	BYTE $0x0F; BYTE $0x77
	RET
```
*/
import "sort"

func bubbleUp(x sort.Interface) {
	n := x.Len()
	for i := 1; i < n; i++ {
		if x.Less(i, i-1) {
			x.Swap(i, i-1)
		}
	}
}

/*
反编译结果如下：

"".bubbleUp STEXT size=182 args=0x10 locals=0x38 funcid=0x0 align=0x0
        0x0000 00000 (.\cmp.go:5)       TEXT    "".bubbleUp(SB), ABIInternal, $56-16
        0x0000 00000 (.\cmp.go:5)       CMPQ    SP, 16(R14)
        0x0004 00004 (.\cmp.go:5)       PCDATA  $0, $-2
        0x0004 00004 (.\cmp.go:5)       JLS     152
        0x000a 00010 (.\cmp.go:5)       PCDATA  $0, $-1
        0x000a 00010 (.\cmp.go:5)       SUBQ    $56, SP
        0x000e 00014 (.\cmp.go:5)       MOVQ    BP, 48(SP)
        0x0013 00019 (.\cmp.go:5)       LEAQ    48(SP), BP
        0x0018 00024 (.\cmp.go:5)       FUNCDATA        $0, gclocals·09cf9819fc716118c209c2d2155a3632(SB)
        0x0018 00024 (.\cmp.go:5)       FUNCDATA        $1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
        0x0018 00024 (.\cmp.go:5)       FUNCDATA        $5, "".bubbleUp.arginfo1(SB)
        0x0018 00024 (.\cmp.go:5)       FUNCDATA        $6, "".bubbleUp.argliveinfo(SB)
        0x0018 00024 (.\cmp.go:5)       PCDATA  $3, $1
        0x0018 00024 (.\cmp.go:6)       MOVQ    AX, "".x+64(SP)
        0x001d 00029 (.\cmp.go:6)       MOVQ    BX, "".x+72(SP)
        0x0022 00034 (.\cmp.go:6)       PCDATA  $3, $-1
        0x0022 00034 (.\cmp.go:6)       MOVQ    24(AX), CX
        0x0026 00038 (.\cmp.go:6)       MOVQ    BX, AX
        0x0029 00041 (.\cmp.go:6)       PCDATA  $1, $0
        0x0029 00041 (.\cmp.go:6)       CALL    CX
        0x002b 00043 (.\cmp.go:6)       MOVQ    AX, "".n+24(SP)
        0x0030 00048 (.\cmp.go:6)       MOVL    $1, CX
        0x0035 00053 (.\cmp.go:7)       JMP     69
        0x0037 00055 (.\cmp.go:7)       MOVQ    "".i+32(SP), DX
        0x003c 00060 (.\cmp.go:7)       LEAQ    1(DX), CX
        0x0040 00064 (.\cmp.go:7)       MOVQ    "".n+24(SP), AX
        0x0045 00069 (.\cmp.go:7)       CMPQ    AX, CX
        0x0048 00072 (.\cmp.go:7)       JLE     142
        0x004a 00074 (.\cmp.go:7)       MOVQ    CX, "".i+32(SP)
        0x004f 00079 (.\cmp.go:8)       MOVQ    "".x+64(SP), DX
        0x0054 00084 (.\cmp.go:8)       MOVQ    32(DX), SI
        0x0058 00088 (.\cmp.go:8)       LEAQ    -1(CX), DI
        0x005c 00092 (.\cmp.go:8)       MOVQ    DI, ""..autotmp_4+40(SP)
        0x0061 00097 (.\cmp.go:8)       MOVQ    "".x+72(SP), AX
        0x0066 00102 (.\cmp.go:8)       MOVQ    CX, BX
        0x0069 00105 (.\cmp.go:8)       MOVQ    DI, CX
        0x006c 00108 (.\cmp.go:8)       CALL    SI
        0x006e 00110 (.\cmp.go:8)       TESTB   AL, AL
        0x0070 00112 (.\cmp.go:8)       JEQ     55
        0x0072 00114 (.\cmp.go:9)       MOVQ    "".x+64(SP), DX
        0x0077 00119 (.\cmp.go:9)       MOVQ    40(DX), SI
        0x007b 00123 (.\cmp.go:9)       MOVQ    "".x+72(SP), AX
        0x0080 00128 (.\cmp.go:9)       MOVQ    "".i+32(SP), BX
        0x0085 00133 (.\cmp.go:9)       MOVQ    ""..autotmp_4+40(SP), CX
        0x008a 00138 (.\cmp.go:9)       CALL    SI
        0x008c 00140 (.\cmp.go:9)       JMP     55
        0x008e 00142 (.\cmp.go:12)      PCDATA  $1, $-1
        0x008e 00142 (.\cmp.go:12)      MOVQ    48(SP), BP
        0x0093 00147 (.\cmp.go:12)      ADDQ    $56, SP
        0x0097 00151 (.\cmp.go:12)      RET
        0x0098 00152 (.\cmp.go:12)      NOP
        0x0098 00152 (.\cmp.go:5)       PCDATA  $1, $-1
        0x0098 00152 (.\cmp.go:5)       PCDATA  $0, $-2
        0x0098 00152 (.\cmp.go:5)       MOVQ    AX, 8(SP)
        0x009d 00157 (.\cmp.go:5)       MOVQ    BX, 16(SP)
        0x00a2 00162 (.\cmp.go:5)       CALL    runtime.morestack_noctxt(SB)
        0x00a7 00167 (.\cmp.go:5)       MOVQ    8(SP), AX
        0x00ac 00172 (.\cmp.go:5)       MOVQ    16(SP), BX
        0x00b1 00177 (.\cmp.go:5)       PCDATA  $0, $-1
        0x00b1 00177 (.\cmp.go:5)       JMP     0
        0x0000 49 3b 66 10 0f 86 8e 00 00 00 48 83 ec 38 48 89  I;f.......H..8H.
        0x0010 6c 24 30 48 8d 6c 24 30 48 89 44 24 40 48 89 5c  l$0H.l$0H.D$@H.\
        0x0020 24 48 48 8b 48 18 48 89 d8 ff d1 48 89 44 24 18  $HH.H.H....H.D$.
        0x0030 b9 01 00 00 00 eb 0e 48 8b 54 24 20 48 8d 4a 01  .......H.T$ H.J.
        0x0040 48 8b 44 24 18 48 39 c8 7e 44 48 89 4c 24 20 48  H.D$.H9.~DH.L$ H
        0x0050 8b 54 24 40 48 8b 72 20 48 8d 79 ff 48 89 7c 24  .T$@H.r H.y.H.|$
        0x0060 28 48 8b 44 24 48 48 89 cb 48 89 f9 ff d6 84 c0  (H.D$HH..H......
        0x0070 74 c5 48 8b 54 24 40 48 8b 72 28 48 8b 44 24 48  t.H.T$@H.r(H.D$H
        0x0080 48 8b 5c 24 20 48 8b 4c 24 28 ff d6 eb a9 48 8b  H.\$ H.L$(....H.
        0x0090 6c 24 30 48 83 c4 38 c3 48 89 44 24 08 48 89 5c  l$0H..8.H.D$.H.\
        0x00a0 24 10 e8 00 00 00 00 48 8b 44 24 08 48 8b 5c 24  $......H.D$.H.\$
        0x00b0 10 e9 4a ff ff ff                                ..J...
        rel 3+0 t=24 type.sort.Interface+96
        rel 3+0 t=24 type.sort.Interface+104
        rel 3+0 t=24 type.sort.Interface+112
        rel 41+0 t=10 +0
        rel 108+0 t=10 +0
        rel 138+0 t=10 +0
        rel 163+4 t=7 runtime.morestack_noctxt+0
go.cuinfo.packagename. SDWARFCUINFO dupok size=0
        0x0000 6d 61 69 6e                                      main
""..inittask SNOPTRDATA size=32
        0x0000 00 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
        0x0010 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        rel 24+8 t=1 sort..inittask+0
type..importpath.sort. SRODATA dupok size=6
        0x0000 00 04 73 6f 72 74                                ..sort
gclocals·09cf9819fc716118c209c2d2155a3632 SRODATA dupok size=10
        0x0000 02 00 00 00 02 00 00 00 02 00                    ..........
gclocals·69c1753bd5f81501d95132d08af04464 SRODATA dupok size=8
        0x0000 02 00 00 00 00 00 00 00                          ........
"".bubbleUp.arginfo1 SRODATA static dupok size=7
        0x0000 fe 00 08 08 08 fd ff                             .......
"".bubbleUp.argliveinfo SRODATA static dupok size=2
        0x0000 00 00
*/
