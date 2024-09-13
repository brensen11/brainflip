#include <iostream>
#include <string>

#include <streambuf>
#include <filesystem>

#include <fstream>

#include <unordered_map>
#include <stack>

#include <stdio.h>  // for fprintf
#include <stdlib.h> // for exit and EXIT_FAILURE

void panic(std::string message) {
    std::cerr << "Panic: " << message << '\n';
    exit(EXIT_FAILURE);
}

void run(std::string program) {
        // std::cout << "I AM CALLED WITH: " << program << '\n';
        unsigned int TAPE_SIZE = 1024*4; // 4KB of Tape
        unsigned char TAPE[TAPE_SIZE];
        unsigned int POINTER = 0;
        unsigned int PC = 0;

        std::unordered_map<unsigned int, unsigned int> bracketPairs;
        std::stack<unsigned int> bracketLefts;

        // init bracketPairs
        for (unsigned int i = 0; i < program.length(); i++) {
            char c = program[i];
            if (c == '[') {
                bracketLefts.push(i);
            } else if (c == ']') {
                unsigned int left = bracketLefts.top();
                bracketLefts.pop();
                bracketPairs[left] = i;
                bracketPairs[i] = left;
            }
        }

        if (!bracketLefts.empty())
            panic("Mismatching [ & ]");
        
        while (PC < program.length()) {
            char cmd = program[PC];
            switch(cmd) {
                // - move the pointer right
                case '>':
                    if (POINTER == TAPE_SIZE) // TAPE_SIZE < MAX_INT, 2 billion
                        panic("Pointer Out of Bounds!");
                    POINTER++;            
                    break;
                
                // - move the pointer left
                case '<':
                    if (POINTER == 0)
                        panic("Pointer Out of Bounds!");
                    POINTER--;
                    break;
                
                // - increment the current cell
                case '+':
                    TAPE[POINTER]++;
                    break;
                
                // - decrement the current cell
                case '-':
                    TAPE[POINTER]--;
                    break;
                
                // - output the value of the current cell
                case '.':
                    std::cout << TAPE[POINTER];
                    break;
                
                // - replace the value of the current cell with input
                case ',':
                    break; // for now ... do nothing :D
                
                // - jump to the matching ] instruction if the current value is zero
                case '[':
                    if (TAPE[POINTER] == 0) {
                        PC = bracketPairs.at(PC);
                        continue; // modifying the PC
                    }
                    break;
                
                // - jump to the matching [ instruction if the current value is not zero
                case ']':
                    if (TAPE[POINTER] != 0) {
                        PC = bracketPairs.at(PC);
                        continue; // modifying the PC
                    }
                    break;
                
                // do nothing on other characters
                default:
                    break;
            }
            PC++;
        }
}

int main(int argc, char** argv) {

    std::string currpath = std::filesystem::current_path().string();
    std::string arg1 = argv[1];
    // std::string arg1 = "/bench.b";

    std::string path = currpath + arg1;
    std::ifstream file(currpath + arg1);
    if(!file.is_open())
        panic("Could not open file " + path);

// https://stackoverflow.com/questions/2602013/read-whole-ascii-file-into-c-stdstring
    std::string program( (std::istreambuf_iterator<char>(file)),
                          std::istreambuf_iterator<char>() );
    run(program);
    file.close();

    // system("pause");
}