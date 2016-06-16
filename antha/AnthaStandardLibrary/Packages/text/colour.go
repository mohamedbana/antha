// prints pre-formatted colour text using ansi codes
package text

import (
	"fmt"
	"github.com/mgutz/ansi"
)

// Prints a description highlighted in red followed by the value in unformatted text
func Print(description string, value interface{}) (fmtd string) {

	fmtd = fmt.Sprintln(ansi.Color(description, "red"), value)
	return
}

func Annotate(stringtoannotate string, colour string) (fmtd string) {
	// currently leads to undesired behaviour with string manipulation
	//fmtd = fmt.Sprint(ansi.Color(stringtoannotate, colour))
	_ = colour
	fmtd = stringtoannotate
	return
}

/*
func Printfield( value interface{}) (fmtd string) {

	switch myValue := value.(type){
		case string:
		fmt.Println(myValue)
		case Hit:
		fmt.Printf("%+v\n", myValue)
		default:
		fmt.Println("Type not handled: ", reflect.TypeOf(value))
	}

	//a := &Hsp{Len: "afoo"}
	val := reflect.Indirect(reflect.ValueOf(value))
	desc := fmt.Sprint(val.Type().Field(0).Name)

	fmtd = fmt.Sprint(ansi.Color(desc, "red"), value)
	return
}
*/
