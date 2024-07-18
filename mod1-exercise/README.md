Exercise: Simple File Generator CLI Tool

go run main.go --content "Hello, World!" filename.txt

Objective:

Create a CLI tool that generates a text file with user-defined content and demonstrates basic command-line flag usage in Go.

Requirements:

The tool should accept the name of the file to create as a command-line argument.
It should use a flag to accept the content that will be written to the file.
If the file already exists, the tool should print an error message and not overwrite the file.
Include a flag for displaying help information that explains how to use the tool.
