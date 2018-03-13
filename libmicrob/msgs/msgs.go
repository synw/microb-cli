package msgs

import (
	"fmt"
	color "github.com/acmacalister/skittles"
)

func Msg(txt string) {
	msg(txt, "none")
}

func Ok(txt string) {
	msg(txt, "ok")
}

func Ready(txt string) {
	msg(txt, "ready")
}

func Warning(txt string) {
	msg(txt, "warning")
}

func Timeout(txt string) {
	msg(txt, "timeout")
}

func Error(txt string) {
	msg(txt, "error")
}

func Bold(txt string) string {
	return color.Bold(txt)
}

func msg(txt string, class string) {
	mess := txt
	if class == "warning" {
		mess = "[" + color.Magenta("Warning") + "] " + txt
	} else if class == "ok" {
		mess = "[" + color.Green("Ok") + "] " + txt
	} else if class == "ready" {
		mess = "[" + color.BoldGreen("Ready") + "] " + txt
	} else if class == "error" {
		mess = "[" + color.BoldRed("Error") + "] " + txt
	} else if class == "timeout" {
		mess = "[" + color.BoldRed("Timeout") + "] " + txt
	}
	fmt.Println(mess)
}

func Debug(obj ...interface{}) {
	for i, el := range obj {
		msg := "[" + color.BoldRed("Debug") + "]"
		fmt.Println(msg, i, el)
	}
}
