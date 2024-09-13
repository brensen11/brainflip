# Compiler
CXX = g++

# Compiler flags
CXXFLAGS = -std=c++17 -Wall -Wextra #-pg -no-pie

# Source file
SRCS = brainflip.cpp

# Output executable
TARGET = brainflip

# Rule to build the executable
$(TARGET): $(SRCS)
	$(CXX) $(CXXFLAGS) $(SRCS) -o $(TARGET)

# Rule to clean up generated files
clean:
	rm -f $(TARGET)
