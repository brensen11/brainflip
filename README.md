# Brainflip

A few toy bf interpreters, compiler, and a collection of benches.

## Build and Compile
This compiler is meant for a windows x64 machine.
You must have nasm, gcc, and go installed.

**Build the Compiler**
```
cd brainflip-go
go build compiler/bfc.go
```

**Compile a program to assembly**
```
./bfc <path/to/program_name.b>
```

**Link and Assembly the program**
```
make <program_name>
```

Full Example
```
cd brainflip-go
go build compiler/bfc.go
./bfc ../benches/hello.b
make hello
./hello.exe
```
