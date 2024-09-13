#ifndef STACK_H
#define STACK_H

// Definition of Stack structure
typedef struct {
    unsigned int* arr;
    unsigned int ptr;
    unsigned int size;
} Stack;

// Function declarations
void stack_panic(const char* message);
Stack* create_stack(unsigned int size);
void free_stack(Stack* stack);
unsigned int pop(Stack* stack);
void push(Stack* stack, unsigned int value);
int is_empty(Stack* stack);

#endif // STACK_H
