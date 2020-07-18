package main

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	source      = "D:\\CarMusic - Prepare" // TODO: config file
	destination = "F:\\"                   // TODO: config file
)

func main() {
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("MP3 sync started ...")
	fmt.Println()

	description := color.New(color.FgGreen)
	info := color.New(color.FgYellow)

	fmt.Printf("%s %s\n", description.Sprint("Source:"), info.Sprint(source))
	fmt.Printf("%s %s\n", description.Sprint("Destination:"), info.Sprint(destination))

	fmt.Println()
	_, _ = title.Println("MP3 sync finished successful!")
}
