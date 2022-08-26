package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

var fset = token.NewFileSet()

func main() {
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

func processPackage(pkg *packages.Package) {
	if pkg.IllTyped {
		fmt.Println("IllTyped=true --> package or dependencies contain errors")
	}

	if len(pkg.Errors) > 0 {
		fmt.Printf("package %v has %v errors\n", pkg.PkgPath, len(pkg.Errors))

		for _, e := range pkg.Errors {
			var errtype string
			switch e.Kind {
			case packages.ListError:
				errtype = "listing/driver"
			case packages.ParseError:
				errtype = "parser"
			case packages.TypeError:
				errtype = "type checker"
			default:
				errtype = "unknown"
			}
			fmt.Printf("Error [%v]: %s\n", errtype, e)
		}
	}

	// This will dump the AST regardless of errors. 转储带有错误信息的AST
	for _, fileAst := range pkg.Syntax {
		ast.Print(fset, fileAst)
	}
}
