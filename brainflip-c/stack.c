#include <stdio.h>  // for fprintf
#include <stdlib.h> // for exit and EXIT_FAILURE
#include "stack.h"

void stack_panic(const char* message) {
    fprintf(stderr, "Stack Panic: %s\n", message);
    exit(EXIT_FAILURE);
}

Stack* create_stack(unsigned int size) {
    Stack* stack = malloc(sizeof(Stack));
    unsigned int* arr = malloc(sizeof(unsigned int) * size);
    stack->arr = arr;
    stack->ptr = 0;
    stack->size = size;
    return stack;
}

void free_stack(Stack* stack) {
    free(stack->arr);
    free(stack);
}

int is_empty(Stack* stack) {
    return stack->ptr == 0;
}

unsigned int pop(Stack* stack) {
    if (stack->ptr == 0)
        stack_panic("popped with nothing left on the stack!");
    stack->ptr--;
    unsigned int value = stack->arr[stack->ptr];
    return value;
}

void push(Stack* stack, unsigned int value) {
    stack->arr[stack->ptr] = value;
    stack->ptr++;
}