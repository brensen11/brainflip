import sys
import os
import subprocess
import asyncio

bash_args = ["C:\\msys64\\usr\\bin\\bash.exe", "--login", "-i", "-c"]
env = {
    **os.environ,
    "MSYSTEM": "UCRT64",
    "CHERE_INVOKING": "1"
}

def prog_num(filename):
    return filename[5:len(filename)-2]

def system(command : str):
    commands = bash_args + [command]
    subprocess.run(commands, shell=True, env=env)

async def run_test(path, prog):
    prog_n = prog_num(prog)
    system(f"./bf.exe -o out-{prog_n}-win.asm {os.path.join(path, prog)}")
    system(f"make out-{prog_n}")
    system(f"./out-{prog_n}.exe < {os.path.join(path, 'input.dat')} > debug/my-output-{prog_n}.dat")
    system(f"diff {os.path.join(path, 'output-' + prog_n + '.dat')} debug/my-output-{prog_n}.dat")

async def main():
    if len(sys.argv) != 2:
        print("Usage: python script.py <path>")
        print("Where <path> is the directory to the bfcheck or bftest folder")
        print("Note that this directory must include a matching output for every prog file")
        print("And must contain the file `input.dat`")
        sys.exit(1)

    path = sys.argv[1]
    if not os.path.isdir(path):
        print(f"The provided path is not a directory: {path}")
        sys.exit(1)
    
    prog_files = [f for f in os.listdir(path) if f.startswith("prog") and f.endswith(".b")]
    
    system("go build bf.go")
    system("rm -rf debug/*")

    tasks=[]
    for prog in prog_files:
        tasks.append(run_test(path, prog))
    
    await asyncio.gather(*tasks)
    
    system("make clean")

if __name__ == "__main__":
    asyncio.run(main())
