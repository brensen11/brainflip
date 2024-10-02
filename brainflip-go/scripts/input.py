# Define the filename
filename = 'test-input.dat'

# Open the file in binary write mode
with open(filename, 'wb') as f:
    # Write bytes 0 to 255
    f.write(bytes(range(1, 256)))
    f.write(bytes(range(257, 512)))
    f.write(bytes([0]))
