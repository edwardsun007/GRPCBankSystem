package main

import "fmt"

type Employee struct {
	Name string
	Age  int
}

func (e Employee) displayInfo() {
	fmt.Printf("Name: %s, Age: %d\n", e.Name, e.Age)
}



func PublicFunc() string {
	return "PublicFunc";
}

func privateFunc() string {
	return "privateFunc";
}

func main() {

	e := Employee{
		Name: "John",
		Age:  30,
	}
	e.displayInfo()

	pubMsg := PublicFunc();
	fmt.Println(pubMsg);

	privateMsg := privateFunc();
	fmt.Println(privateMsg);


	// fmt.Println("Hello world");
	// age := 27;
	// fmt.Println(age);

	// var cities = [...]string {"Seattle", "San Jose"};
	// fmt.Println(cities[0]);
	// fmt.Println(cities[1]);

	// var userNames = [...]string{"jinzhu1", "jinzhu2",  "jinzhu3"};
	// for _, name := range userNames {
	// 	fmt.Println(name);
	// }

}