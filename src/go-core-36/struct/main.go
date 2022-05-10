package main

import "fmt"

type (
	AnimalCategory struct {
		// ...
		species string // 种
	}

	Animal struct {
		scientificName string // 学名
		AnimalCategory
	}

	Cat struct {
		name string
		Animal
	}
)

// 方法
func (ac AnimalCategory) String() string {
	return fmt.Sprintf("%s", ac.species)
}
func (a Animal) String() string {
	return fmt.Sprintf("%s", a.scientificName)
}
func (cat Cat) String() string {
	return fmt.Sprintf("%s (category: %s, name: %q)",
		cat.scientificName, cat.Animal.AnimalCategory, cat.name)
}

// 值接收器的方法，对其值进行更改是不会体现在原值上的
// 如果对值对象的某个属性进行更改，在方法结束后其值还是原来的值
// 如下的 cat.SetName 方法，因此为了避免不必要的错误，ide也会提示用户”无效的字段分配“
func (cat Cat) SetName(name string) {
	cat.name = name
}

// 指针接收器方法对属性的更改才会体现在原值上
func (cat *Cat) SetName2(name string) {
	cat.name = name
}
func main() {
	catogory := AnimalCategory{
		species: "cat",
	}
	fmt.Printf("the animal category: %s\n", catogory)

	animal := Animal{
		scientificName: "demo1",
	}

	fmt.Printf("Animal.String() = %s", animal.String())
}
