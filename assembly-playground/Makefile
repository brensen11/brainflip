# Windows: nasm and ld
win-ot: hello-win-ot

hello-win-ot: hello-win-ot.o
	ld hello-win-ot.o -e start -lkernel32 -o hello-win-ot.exe

hello-win-ot.o: hello-win-ot.asm
	nasm -fwin64 hello-win-ot.asm -o hello-win-ot.o

# Windows: nasm and gcc
win: hello-win

hello-win: hello-win.o
	gcc hello-win.o -o hello-win

hello-win.o: hello-win.asm
	nasm -fwin64 hello-win.asm -o hello-win.o

jh: hello-jh

# Linux: nasm and ld
hello-jh: hello-jh.o
	ld -m elf_i386 -o hello-jh hello-jh.o

hello-jh.o: hello-jh.asm
	nasm -felf32 hello-jh.asm -o hello-jh.o

lll: hello-lll

# Linux: as and gcc
hello-lll: hello-lll.o
	gcc -nostdlib -static -o hello-lll hello-lll.o

hello-lll.o: hello-lll.asm
	as hello-lll.asm -o hello-lll.o


# Keep it clean
clean:
	rm -f hello-lll hello-jh hello-win hello-win-ot *.o *.obj


