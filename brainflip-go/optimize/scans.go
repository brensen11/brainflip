package optimize

import (
	lp "brainflip-go/lexparse"
	"brainflip-go/utils"
)

func Optimize_scans(instructions *[]lp.Instruction) {
	scans := lp.Locate_Scans(instructions)
	for i := len(scans) - 1; i >= 0; i-- {
		scan := scans[i]
		right := scan.Moves > 0
		l := scan.L
		r := scan.R
		var scan_dir string
		if right {
			scan_dir = "add"
		} else {
			scan_dir = "sub"
		}
		var ri utils.Builderf
		ri.Add_instr("; Scan Code")
		ri.Add_instr("xor     rcx, rcx")
		ri.Add_instr("cmp     BYTE [rdi], 0")
		ri.Add_instr("je      right_vector_%d", l)
		ri.Add_instr("movdqu  xmm0, [zeroes] ; load 0s")
		if !right {
			ri.Add_instr("sub     rdi, 15")
		}

		ri.Add_label("left_vector_%d", l)
		ri.Add_instr("movdqu  xmm1, [rdi] ; load data from tape")
		ri.Add_instr("pcmpeqb xmm1, xmm0 ;  zero_count = CMEQ.16B input, 0s; all 0s marked with -1 else 0")
		ri.Add_instr("pmovmskb eax, xmm1 ; mov msb of xmm1 into eax, now data is in bits instead of bytes")
		ri.Add_instr("and      eax, [mask_%d]", abs(scan.Moves))
		ri.Add_instr("movzx    eax, ax")
		if right {
			ri.Add_instr("tzcnt    ecx, eax ; count the trailing 0s of the register, which will tell me the index")
		} else {
			ri.Add_instr("lzcnt    ecx, eax ; count the trailing 0s of the register, which will tell me the index")
			ri.Add_instr("sub      ecx, 16 ; subtract 16 to account for 32 bit register since we are counting leading 0's, not trailing")
		}
		ri.Add_instr("cmp      ecx, 16")
		ri.Add_instr("jl      end_vector_%d", l)
		ri.Add_instr("%s      rdi, 16", scan_dir)
		ri.Add_instr("jmp      left_vector_%d", l)
		ri.Add_label("end_vector_%d", l)
		if right {
			ri.Add_instr("add     rdi, rcx") // Jump from RDI to index offset
		} else {
			ri.Add_instr("mov     rdx, 15 ; 15 because we're IN the zeroes and we don't want to move OUT")
			ri.Add_instr("sub     rdx, rcx")
			ri.Add_instr("add     rdi, rdx") // duplicate but, more clear idk
		}

		ri.Add_label("right_vector_%d", l)
		ri.Add_instr("; Scan Code End")

		var replace []lp.Instruction
		replace = append(replace, lp.Raw{ri.String()})

		*instructions = lp.Instructions_replace(*instructions, l, r+1, replace)
	}
}
