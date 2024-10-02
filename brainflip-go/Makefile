LK=gcc
ASM=nasm

%: %-win.o runtime.o
	$(LK) -g -o $@ $< runtime.o

%-win.o: %-win.asm
	$(ASM) -o $@ -fwin64 $<

runtime.o: runtime.c
	$(LK) -c runtime.c -o runtime.o

clean:
	rm -f *.asm *.exe *.dat *.data

# %: %.asm
# 	rule: $@ and dependency: $<