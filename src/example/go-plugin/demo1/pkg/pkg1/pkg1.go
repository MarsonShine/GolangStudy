package pkg1

import (
	"errors"
	"log"
	"plugin"
)

func init() {
	log.Println("pkg1 init")
}

type MyInterface interface {
	M1()
}

func LoadAndInvokeSomethingFromPlugin(pluginPath string) error {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}

	// 导出成员信息
	v, err := p.Lookup("V")
	if err != nil {
		return err
	}
	*v.(*int) = 7
	//	f.(func())() // prints "Hello, number 7"

	f, err := p.Lookup("F")
	if err != nil {
		return err
	}
	f.(func())()

	f1, err := p.Lookup("Foo")
	if err != nil {
		return err
	}
	i, ok := f1.(MyInterface)
	if !ok {
		return errors.New("f1 does not implement MyInterface")
	}
	i.M1()

	return nil
}
