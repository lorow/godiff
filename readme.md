# What is this?

Godiff is basically a simple tool for debugging two requests side by side.
It's a very specific tool for a very specific use case - testing requests and comparing the output.

[todo write more about it]

# How to debug?

Since we're using bubbletea as the TUI framework, debugging is a bit awkward. The general idea is to build
the binary with all debug flags enabled, run it and then attach to the running process with a debugger

go build -gcflags="all=-N -l" -o godiff.exe && ./godiff.exe
