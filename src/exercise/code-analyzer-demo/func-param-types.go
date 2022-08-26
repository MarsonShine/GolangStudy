package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

// 如何分析一个代码的成员信息？
// 类型参数、返回类型、方法名、参数类型等

var fset = token.NewFileSet()

func main() {
	/*
		第一步：初始化配置，读取的文件路径等
	*/
	const mode packages.LoadMode = packages.NeedName |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Expecting a single argument: directory of module")
	}

	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: flag.Args()[0]}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		processPackage(pkg)
	}
}

/*
	找到指定包的所有成员信息
*/
func processPackage(pkg *packages.Package) {

	for _, fileAst := range pkg.Syntax { // 获取包内所有语法树文件
		ast.Inspect(fileAst, func(n ast.Node) bool {
			// 找到函数申明
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				processFuncDeclare(funcDecl, pkg.TypesInfo)
			}
			return true
		})
	}
}

/*
	分析函数参数信息
*/
func processFuncDeclare(fd *ast.FuncDecl, tinfo *types.Info) {
	fmt.Println("=== Function", fd.Name)
	// 遍历形参列表并检查每个形参的类型
	for _, field := range fd.Type.Params.List {
		var names []string
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
		fmt.Println("param:", strings.Join(names, ", "))
		processTypeExpr(field.Type, tinfo)
	}
}

/*
	查找struct类型字段信息
*/
// 函数递归地剥离指针和数组/切片的类型
func processTypeExpr(e ast.Expr, tinfo *types.Info) {
	switch tyExpr := e.(type) {
	case *ast.StarExpr:
		fmt.Println(" pointer to ...")
		processTypeExpr(tyExpr.X, tinfo)
	case *ast.ArrayType:
		fmt.Println(" slice or arrat of ...")
		processTypeExpr(tyExpr, tinfo)
	default:
		switch ty := tinfo.Types[e].Type.(type) {
		case *types.Basic: // 基本类型 如 string
			fmt.Println(" basic =", ty.Name())
		case *types.Named: // 命名类型 如 MyType，也可以为命名类型（匿名类型）struct {x int}
			fmt.Println(" name =", ty)
			uty := ty.Underlying()
			fmt.Println("  type =", ty.Underlying())
			if sty, ok := uty.(*types.Struct); ok {
				fmt.Println("  fields:")
				// 处理结构参数类型信息
				processStructParamType(sty)
			}
			fmt.Println("  pos =", fset.Position(ty.Obj().Pos()))
		default:
			fmt.Println("  unnamed type =", ty)
			if sty, ok := ty.(*types.Struct); ok {
				fmt.Println("  fields:")
				processStructParamType(sty)
			}
		}
	}
}

/*
	处理结构字段类型
*/
func processStructParamType(sty *types.Struct) {
	for i := 0; i < sty.NumFields(); i++ {
		field := sty.Field(i)
		fmt.Println("     ", field.Name(), field.Type())
	}
}
