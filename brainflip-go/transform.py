import sys

# Define the function to read the file and print byte values
def print_byte_values(file_path):
    try:
        with open(file_path, 'rb') as file:  # Open the file in binary mode
            byte = file.read(1)  # Read one byte at a time
            while byte:
                print(int.from_bytes(byte))  # Print the numeric value of the byte
                byte = file.read(1)  # Read the next byte
    except FileNotFoundError:
        print(f"The file {file_path} was not found.")
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    # Check if a filename was provided
    if len(sys.argv) != 2:
        print("Usage: python script.py <filename>")
    else:
        file_path = sys.argv[1]  # Take the first command-line argument as the filename
        print_byte_values(file_path)
