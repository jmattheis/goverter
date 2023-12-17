package example

// goverter:converter
// goverter:extend ExtractFriendNames
type Converter interface {
	Convert(source []InputPerson) []OutputPerson
}

type InputPerson struct {
	Name    string
	Friends []InputPerson
}
type OutputPerson struct {
	Name    string
	Friends []string
}

func ExtractFriendNames(persons []InputPerson) []string {
	var names []string
	for _, person := range persons {
		names = append(names, person.Name)
	}
	return names
}
