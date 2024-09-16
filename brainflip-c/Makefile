CC = gcc
PROFILE_FLAGS = -fprofile-arcs -ftest-coverage -pg -no-pie
CFLAGS = -Wall -Wextra -O2
TARGET = brainflip

all: $(TARGET)

$(TARGET): stack.o brainflip.o
	$(CC) $(CFLAGS) -o $(TARGET) stack.o brainflip.o

stack.o: stack.c stack.h
	$(CC) $(CFLAGS) -c stack.c

brainflip.o: brainflip.c stack.h
	$(CC) $(CFLAGS) -c brainflip.c

clean:
	rm -f $(TARGET) *.o
