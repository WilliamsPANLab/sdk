package main

import (
	. "fmt"
)

func main() {
	// Create frontend printers
	printers := []Printer{
		NewBasicPrinter("go", "templates/template.go.", "dist/bridge.go"),
		NewBasicPrinter("py", "templates/template.py", "dist/python/flywheel/flywheel.py"),
		NewBasicPrinter("m", "templates/template.m", "dist/matlab/Flywheel.m"),
		NewBasicPrinter("m-binary", "templates/template.binary.m", "dist/matlab-binary/Flywheel.m"),
		NewBasicPrinter("cli", "templates/template.binary", "dist/binary/sdk.go"),
	}

	// Parse source code into useful sigs
	Println("Parsing SDK...")
	fset, sigs := GenerateSignatures("../api")
	_ = fset

	// Detect the SDK version and add that to template struct
	sigs.Version = DetectSDKVersion("../sdk.go")

	// Load templates and execute
	for _, printer := range printers {
		Println("Generating", printer.Name(), "code...")

		printer.Init()
		printer.Print(sigs)
	}
}
