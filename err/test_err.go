package err

import (
	"log"
	"testing"
)

func Test_Foo(t *testing.T) {
	x, err := Foo()
	if err != nil {
		log.Fatal("Error ", err)
	}
	_ = x
}

func Foo() (int, error) {
	return 712, nil
}
