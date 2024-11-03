package main

import (
	"ReflectSample/StructMarshaller"
	"fmt"
)

type SampleStruct struct {
	Name        string
	Age         int
	AverageMark float32
	comment     string
}

func main() {
	s := SampleStruct{
		Name:        "Вася",
		Age:         20,
		AverageMark: 4.8,
		comment:     "Молодец",
	}

	fmt.Printf("%+v \n", s)
	var s2 SampleStruct

	marshalledFields, err := StructMarshaller.StructMarshall(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(marshalledFields)
	marshalledFields["comment"] = "Молодец"

	err = StructMarshaller.StructUnmarshall(&s2, marshalledFields)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v \n", s2)
}
