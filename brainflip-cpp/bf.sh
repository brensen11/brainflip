./build/bf.exe $1.b > $1.ll
llvm-as $1.ll -o $1.bc
llc -filetype=obj $1.bc -o $1.o
gcc -g -o $1 $1.o runtime.o