package gox

import "fmt"

func NotImplementedYet(fun string) error {
	return fmt.Errorf("%v not implemented yet", fun)
}


func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
