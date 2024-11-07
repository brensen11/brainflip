filename=$(basename "$1" .b)
./build/bf.exe $1 > $filename.ll
llvm-as $filename.ll -o $filename.bc
llc -filetype=obj $filename.bc -o $filename.o
gcc -g -o $filename $filename.o runtime.o