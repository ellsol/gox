package utilx

import "fmt"


const(
	MaxInt64 = 9223372036854775807
)

func NotImplementedYet(fun string) error {
	return fmt.Errorf("%v not implemented yet", fun)
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
