package main

import (
	"fmt"
	"reflect"
)

func swap[T any](a, b T) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	// Only pointers can be swapped.
	if valueA.Kind() != reflect.Pointer {
		panic("first argument must be a pointer")
	}
	if valueB.Kind() != reflect.Pointer {
		panic("second argument must be a pointer")
	}

	// Check nil pointer.
	if valueA.IsNil() {
		panic("first argument is a nil pointer")
	}
	if valueB.IsNil() {
		panic("second argument is a nil pointer")
	}

	elemA := valueA.Elem()
	elemB := valueB.Elem()

	// Check both are same type.
	if elemA.Type() != elemB.Type() {
		panic("arguments are different types")
	}

	// Check settable.
	if !elemA.CanSet() {
		panic("first argument is not settable")
	}
	if !elemB.CanSet() {
		panic("second argument is not settable")
	}

	// Swap.
	temp := reflect.ValueOf(elemA.Interface())
	elemA.Set(elemB)
	elemB.Set(temp)
}

func main() {
	a := 10
	b := 20

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)

	swap(&a, &b)

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)
}
