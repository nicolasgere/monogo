package lib

import (
	hello "hello/lib"
)

func Toto() string {
	return hello.Hello() + "toto"
}
