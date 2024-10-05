package optimize

import (
	lp "brainflip-go/lexparse"
	"brainflip-go/utils"
)

func Optimize_scans(instructions *[]lp.Instruction) {
	scans := lp.Locate_Scans(instructions)
	for i := len(scans) - 1; i >= 0; i-- {
		scan := scans[i]
		l := scan.L
		r := scan.R
		var dir_instr string
		if scan.Moves > 0 {
			dir_instr = "add"
		} else {
			dir_instr = "sub"
		}
		var ri utils.Builderf
		ri.Add_instr("; Scan Code")
		ri.Add_instr("xor     rcx, rcx")
		ri.Add_instr("cmp     BYTE [rdi], 0")
		ri.Add_instr("je      right_vector_%d", l)
		ri.Add_instr("movdqu  xmm1, [ones]")
		ri.Add_instr("movdqu  xmm2, [zeroes] ; load 0s")
		ri.Add_instr("movdqu  xmm3, [indices_%d]", abs(scan.Moves))

		ri.Add_label("left_vector_%d", l)
		ri.Add_instr("movdqu  xmm0, [rdi] ; load data from tape")
		ri.Add_instr("pcmpeqb xmm0, xmm2 ;  zero_count = CMEQ.16B input, 0s; all 0s marked with -1 else 0")
		ri.Add_instr("pxor    xmm0, xmm1 ; not_zero_count = xor 0xFF.., zero_count")
		ri.Add_instr("por     xmm0, xmm3 ; masked = ORN.16B indices, zeroes")
		ri.Add_instr("pmovmskb eax, xmm0 ; mov msb of xmm0 into eax, now data is in bits instead of bytes")
		ri.Add_instr("xor      eax, 0xFFFFFFFF")
		ri.Add_instr("tzcnt    ecx, eax ; count the trailing 0s of the register, which will tell me the index")
		ri.Add_instr("cmp      ecx, 16")
		ri.Add_instr("jne      end_vector_%d", l)
		ri.Add_instr("%s      rdi, 16", dir_instr)
		ri.Add_instr("jmp      left_vector_%d", l)

		ri.Add_label("end_vector_%d", l)
		ri.Add_instr("%s      rdi, rcx", dir_instr)

		ri.Add_label("right_vector_%d", l)
		ri.Add_instr("; Scan Code End")

		var replace []lp.Instruction
		replace = append(replace, lp.Raw{ri.String()})

		*instructions = lp.Instructions_replace(*instructions, l, r+1, replace)
	}
}
