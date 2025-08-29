package main

import "fmt"

type Human struct {
	Name    string
	Surname string
	Age     int
}

func (h Human) Talk() {
	fmt.Println("Hello!")
}

func (h Human) Sleep() {
	fmt.Println("i am going to sleep")
}

type Action struct {
	Human
}

func (a Action) Dance() {
	fmt.Println("I am dancing!")
}

func (a Action) Sing() {
	fmt.Println("La la la")
}

func main() {
	john := Action{
		Human: Human{
			Name:    "John",
			Surname: "Ivanov",
			Age:     25,
		},
	}

	john.Talk()
	fmt.Println("My name is", john.Name)
	john.Dance()

}
