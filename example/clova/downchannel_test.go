package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

type Person struct {
	Name string
	Age  int64
}

func (person *Person) FillStruct(rawData map[string]interface{}) error {
	for key, value := range rawData {
		err := SetField(person, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Test1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData["Name"] = "Tony"
	rawData["Age"] = int64(23)

	person := &Person{}
	err := person.FillStruct(rawData)
	if err != nil {
		panic(err)
	}
	fmt.Println(person)
}
