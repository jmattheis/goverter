package example

// goverter:converter
// goverter:extend ConvertAnimals
type Converter interface {
	Convert(source Input) Output

	// used only in extend method
	ConvertDogs([]Dog) []Animal
	ConvertCats([]Cat) []Animal
}

type Input struct {
	Animals InputAnimals
}
type InputAnimals struct {
	Cats []Cat
	Dogs []Dog
}
type Output struct {
	Animals []Animal
}

type Cat struct{ Name string }
type Dog struct{ Name string }

type Animal struct{ Name string }

func ConvertAnimals(c Converter, input InputAnimals) []Animal {
	dogs := c.ConvertDogs(input.Dogs)
	cats := c.ConvertCats(input.Cats)
	return append(dogs, cats...)
}
