package main

import "fmt"

type struct1 struct {
	a []*struct2
	b []*struct2
}

type struct2 struct {
	f float64
}

func main() {
	struct1 := struct1{}
	struct1.a = append(struct1.a, &struct2{0.5})
	struct1.b = append(struct1.b, &struct2{0.5})

	struct2 := struct2{0.1}

	struct1.a = append(struct1.a, &struct2)
	struct1.b = append(struct1.b, &struct2)

	fmt.Printf("%p\n", struct1.a[0])
	fmt.Printf("%p\n", struct1.b[0])
	fmt.Println("---------------")
	fmt.Printf("%p\n", struct1.a[1])
	fmt.Printf("%p\n", struct1.b[1])
	fmt.Println("---------------")
	fmt.Printf("%p\n", &struct2)
}
