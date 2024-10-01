LK=gcc
ASM=nasm

%: %-win.o
	$(LK) -g -o $@ $<

%-win.o: %-win.asm
	$(ASM) -o $@ -fwin64 $<

clean:
	rm -f *.asm *.exe

# %: %.asm
# 	rule: $@ and dependency: $<