package main_test

import (
	"testing"
	"go-crapy/message"
)

func Test_main(t *testing.T) {
	err := message.SendMsg("aaa")
	if err != nil {
		panic(err)
	}
}
