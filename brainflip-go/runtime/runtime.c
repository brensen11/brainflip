#include <windows.h>

void my_putchar(unsigned char c) {
    DWORD written;
    HANDLE hConsole = GetStdHandle(STD_OUTPUT_HANDLE);  
    WriteFile(hConsole, &c, 1, &written, NULL);         
}

unsigned char my_getchar() {
    HANDLE hStdin = GetStdHandle(STD_INPUT_HANDLE);
    unsigned char buffer[1]; // Buffer to hold the byte read
    DWORD bytesRead;

    if (ReadFile(hStdin, buffer, sizeof(buffer), &bytesRead, NULL)) {
        if (bytesRead == 1) {
            return buffer[0];
        }
    }
    
    return -1;
}