#include <stdio.h>  // for fprintf
#include <stdlib.h> // for exit and EXIT_FAILURE
#include <string.h> // for strlen
#include "stack.h"

void panic(char* message) {
    fprintf(stderr, "Panic: %s\n", message);
    exit(EXIT_FAILURE);
}

void print_tape(unsigned char TAPE[], unsigned int TAPE_SIZE) {
    printf("TAPE: [");
    for (size_t i = 0; i < TAPE_SIZE; i++)
    {
        printf("%d", TAPE[i]);
    }
    printf("]\n");
    
}

void run(char* program) {
        unsigned int length = strlen(program);
        // std::cout << "I AM CALLED WITH: " << program << '\n';
        unsigned int TAPE_SIZE = 1024 * 1024 * 4; // 4KB of Tape
        unsigned char* TAPE = (unsigned char*)calloc(TAPE_SIZE, sizeof(int));
        unsigned int POINTER = TAPE_SIZE / 2;
        unsigned int PC = 0;

        unsigned int bracketPairs[length];
        Stack* stack = create_stack(TAPE_SIZE / 2);

        // init bracketPairs
        for (unsigned int i = 0; i < length; i++) {
            char c = program[i];
            if (c == '[') {
                push(stack, i);
            } else if (c == ']') {
                unsigned int left = pop(stack);
                bracketPairs[left] = i;
                bracketPairs[i] = left;
                // printf("Added pair at %d and %d\n", left, i);
            }
        }

        if (!is_empty(stack))
            panic("Mismatching [ & ]");
        free_stack(stack);
        
        while (PC < length) {
            // print_tape(TAPE, TAPE_SIZE);
            char cmd = program[PC];
            switch(cmd) {
                // - move the pointer right
                case '>':
                    if (POINTER == TAPE_SIZE) // TAPE_SIZE < MAX_INT, 2 billion
                        panic("Pointer Out of Bounds! Increase tape size");
                    POINTER++;            
                    break;
                
                // - move the pointer left
                case '<':
                    if (POINTER == 0)
                        panic("Pointer Out of Bounds! Increase tape size");
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
                    fputc(TAPE[POINTER], stdout);
                    break;
                
                // - replace the value of the current cell with input
                case ',':
                    char ch = getchar();
                    TAPE[POINTER] = ch;
                    break; // for now ... do nothing :D
                
                // - jump to the matching ] instruction if the current value is zero
                case '[':
                    if (TAPE[POINTER] == 0) {
                        PC = bracketPairs[PC];
                        continue; // modifying the PC
                    }
                    break;
                
                // - jump to the matching [ instruction if the current value is not zero
                case ']':
                    if (TAPE[POINTER] != 0) {
                        PC = bracketPairs[PC];
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
    if (argc < 2) {
        printf("Usage: ./brainflip <path-to-bf.b>");
        return 1;
    }

    char* filename = argv[1];
    FILE* file = fopen(filename, "r");
    if(!file) {
        printf("Could not open %s", filename);
        return 1;
    }

    fseek(file, 0, SEEK_END);
    long size = ftell(file);
    rewind(file);

    char* content = malloc(size + 1);
    fread(content, 1, size, file);
    content[size] = '\0';
    fclose(file);
    run(content);
}