package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_checkContractNameValidate(t *testing.T) {

	cliInstance := NewCLI()
	// Only pass t into top-level Convey calls
	Convey("Given name is empty", t, func() {
		Convey("When the name is empty", func() {
			name1 := ""
			flag1 := cliInstance.checkContractNameValid(name1)

			Convey("The value should be false", func() {
				So(flag1, ShouldEqual, false)
			})
		})

		Convey("When the name is bottos", func() {
			name2 := "bottos"
			flag2 := cliInstance.checkContractNameValid(name2)

			Convey("The value should be true", func() {
				So(flag2, ShouldEqual, true)
			})
		})

		Convey("When the name is test@bottos", func() {
			name2 := "test@bottos"
			flag2 := cliInstance.checkContractNameValid(name2)

			Convey("The value should be true", func() {
				So(flag2, ShouldEqual, true)
			})
		})
	})

}
