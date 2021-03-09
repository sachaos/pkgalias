package a

import (
	format "fmt" // want "invalid alias"
	"os"         // want "use alias"

	m "math"  // OK
	"strings" // OK
)

func f() {
	format.Println("HelloWorld")
	os.Getenv("TEST")

	strings.Trim("hogehoge", "h")
	m.Abs(10)
}
