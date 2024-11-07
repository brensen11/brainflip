./build/bf.exe $1 > $2.ll
llvm-as $2.ll -o $2.bc
llc -filetype=obj $2.bc -o $2.o
gcc -g -o $2 $2.o runtime.o