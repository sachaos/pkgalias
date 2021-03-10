package a

import (
	farmat "fmt" // want "invalid alias"
	"os"         // want "invalid alias"

	m "math"  // OK
	"strings" // OK
	"net"     // OK
)

func f() {
	farmat.Println("HelloWorld") // want "invalid alias"
	os.Getenv("TEST")           // want "invalid alias"

	f := os.File{} // want "invalid alias"
	f.Name()

	e := net.AddrError{}
	e.Error()

	strings.Trim("hogehoge", "h")
	m.Abs(10)
}
