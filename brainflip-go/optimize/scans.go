package optimize

import (
	lp "brainflip-go/lexparse"
	"fmt"
)

func Optimize_scans(instructions *[]lp.Instruction) {
	scans := lp.Locate_Scans(instructions)
	for i := len(scans) - 1; i >= 0; i-- {
		scan := scans[i]
		l := scan.L
		r := scan.R
		raw_instructions := fmt.Sprintf(`; Scan Code
	xor     rcx, rcx
	cmp     BYTE [rdi], 0
	je      right_vector_%d
left_vector_%d:
	movdqu  xmm0, [rdi] ; load data from tape
	movdqu  xmm1, [zeroes] ; load 0s
	pcmpeqb xmm0, xmm1 ;  zero_count = CMEQ.16B input, 0s; all 0s marked with -1 else 0

	; if ecx 16, loop
	; else adjust rdi (p[0]) by ecx

	movdqu  xmm1, [ones]
	pxor    xmm0, xmm1 ; not_zero_count = xor 0xFF.., zero_count
	movdqu  xmm1, [indices_%d]
	por     xmm0, xmm1 ; masked = ORN.16B indices, zeroes

	pmovmskb eax, xmm0 ; mov msb of xmm0 into eax, now data is in bits instead of bytes
	xor      eax, 0xFFFFFFFF
	tzcnt    ecx, eax ; count the trailing 0s of the register, which will tell me the index
	cmp      ecx, 16
	jne      end_vector_%d
	add      rdi, 16
	jmp      left_vector_%d
end_vector_%d:
	add      rdi, rcx
right_vector_%d:
; Scan Code End
		`, l, l, abs(scan.Moves), l, l, l, l)
		var replace []lp.Instruction
		replace = append(replace, lp.Raw{raw_instructions})

		*instructions = lp.Instructions_replace(*instructions, l, r+1, replace)
	}
}
