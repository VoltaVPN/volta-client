package main

import "github.com/voltavpn/volta-client/internal/gui"

// main wires the executable to the GUI layer.
// All UI details live in the internal/gui package.
func main() {
	gui.Run()
}
