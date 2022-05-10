package main

import "fmt"

type (
	Pet interface {
		// SetName(name string)
		Name() string
		Category() string
	}

	Dog struct {
		name string
		// category string
	}

	// 接口的嵌套，组合
	Animal interface {
		SpecificName() string
		Category() string
		Name() string
	}

	Pet1 interface {
		Animal
		Name() string
	}
)

func (dog Dog) Name() string {
	return dog.name
}
func (dog Dog) Category() string {
	return ""
}
func (dog *Dog) SetName(name string) {
	dog.name = name
}

func main() {
	dog := Dog{"斑点狗"}
	// var pet Pet = dog	//值对象无法隐式转换为pet，因为没有完全实现 Pet 接口的SetName方法
	// var pet Pet = &dog
	dog2 := dog
	dog.name = "藏獒"
	fmt.Printf("dog.name=%s, dog2.name=%s", dog.name, dog2.name) // 值彼此不会受到影响

	var dog1 *Dog
	fmt.Printf("the first dog is nil.")
	dog3 := dog1
	fmt.Printf("the second dog is nil.")
	var pet Pet = dog3
	if pet == nil { // 编译器提示报错，this comparison is never true; the lhs of the comparison has been assigned a concretely typed value，上述的 pet = dog3 与 pet = nil 不一样，dog3虽然为nil，但是它携带了类型类信
		fmt.Printf("the pet is nul.")
	} else {
		fmt.Printf("the pet is not null.")
	}
}
