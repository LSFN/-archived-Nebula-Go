package main

import (
	"log"

	"code.google.com/p/goprotobuf/proto"
	"./example.pb"
)

func main() {
	test := &example.Test{
		Label: proto.String("hello"),
		Type:  proto.Int32(17),
		Optionalgroup: &example.Test_OptionalGroup{
			RequiredField: proto.String("good bye"),
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	newTest := new(example.Test)
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	// Now test and newTest contain the same data.
	if test.GetLabel() != newTest.GetLabel() {
		log.Fatalf("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
	}
	// etc.
}