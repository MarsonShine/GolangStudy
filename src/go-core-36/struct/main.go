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
)

// 方法
func (ac AnimalCategory) String() string {
	return fmt.Sprintf("%s", ac.species)
}
func main() {
	catogory := AnimalCategory{
		species: "cat",
	}
	fmt.Printf("the animal category: %s\n", catogory)

	animal := Animal{
		scientificName: "demo1",
	}
}
