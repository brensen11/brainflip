LK=gcc
ASM=nasm

%: %-win.o runtime.o
	$(LK) -g -o $@ $< runtime.o

%-win.o: %-win.asm
	$(ASM) -o $@ -fwin64 $<

runtime.o: runtime/runtime.c
	$(LK) -g -c runtime/runtime.c -o runtime.o

clean:
	rm -f *.asm *.exe

# %: %.asm
# 	rule: $@ and dependency: $<