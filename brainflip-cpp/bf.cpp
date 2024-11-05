#include <iostream>
#include <fstream>
#include <sstream>
#include <string>

int compile(std::string program)
{
    for (char ch : program)
    {
        switch (ch) {
		case '>':
			// asm_b.Add_instr("inc     %s", TAPE_PTR)
            break;
		case '<':
			// asm_b.Add_instr("dec     %s", TAPE_PTR)
            break;
		case '+':
			// asm_b.Add_instr("inc     BYTE [%s]", TAPE_PTR)
            break;
		case '-':
			// asm_b.Add_instr("dec     BYTE [%s]", TAPE_PTR)
            break;
		case '.':
			// asm_b.Add_instr("mov     cl, BYTE [%s]", TAPE_PTR)
			// asm_b.Add_instr("call    my_putchar")
            break;
		case ',':
			// asm_b.Add_instr("call    my_getchar")
			// asm_b.Add_instr("mov     BYTE [%s], al", TAPE_PTR)
            break;
		case '[':
			// asm_b.Add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			// asm_b.Add_instr("je      right_%s", strconv.Itoa(bracket_pairs[i]))
			// asm_b.Add_label("left_%s", strconv.Itoa(i))
            break;
		case ']':
			// asm_b.Add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			// asm_b.Add_instr("jne     left_%s", strconv.Itoa(bracket_pairs[i]))
			// asm_b.Add_label("right_%s", strconv.Itoa(i))
            break;
		default:
			break; // do nothing
		}
    }
}

int main(int argc, char *argv[])
{
    if (argc <= 1)
    {
        std::cout << "Usage: ./bf <bf-program.b>";
        return 0;
    }

    std::ifstream file(argv[1]);
    if (!file)
    { // Check if the file was successfully opened
        std::cerr << "Could not open the file" << argv[1] << std::endl;
        return 1;
    }

    std::stringstream buffer;
    buffer << file.rdbuf();

    std::string program = buffer.str();

    std::cout << program << std::endl;

    return 0;
}