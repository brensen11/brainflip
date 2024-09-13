import os
import subprocess
import platform

# Directory to work in
directory = "../benches"
cmd = "./brainflip" if platform.system() != "Windows" else "./brainflip.exe"

# Loop through each file in the directory
for filename in os.listdir(directory):
    file_path = os.path.join(directory, filename)
    if os.path.isfile(file_path):  # Check if it is a file
        print(f"Benching {file_path}")
        output_file = f"./bench-out/{filename.split('.')[0]}.out"
        with open(output_file, 'w') as outfile:
            subprocess.run([cmd, file_path], stdout=outfile)
