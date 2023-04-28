import sys
import os
import glob
import chess.pgn
import subprocess

# Check that the correct number of command line arguments have been provided
if len(sys.argv) != 3:
    print("Usage: python program.py <directory> <output_file>")
    sys.exit(1)

# Get the directory name and the output file name from the command line arguments
directory = sys.argv[1]
output_file = sys.argv[2]

# Use glob to get a list of all the PGN files in the directory
pgn_files = glob.glob(os.path.join(directory, "*.pgn"))

# Open a pipe to the sort command for writing
sort_pipe = subprocess.Popen(["sort", "-u", "-o", output_file], stdin=subprocess.PIPE)

# Loop through each PGN file and parse its contents using the chess library
for pgn_file in pgn_files:
    with open(pgn_file, "r") as f:
        # Use the chess.pgn.read_game() function to read the PGN file and create a chess.Game object
        game = chess.pgn.read_game(f)
        board = game.board()
        # Loop through each move in the game and write the floating point score representation of the result in square brackets and its FEN representation to the output file
        for move in game.mainline_moves():
            board.push(move)
            result = game.headers["Result"]
            if result == "1-0":
                score = 1.0
            elif result == "0-1":
                score = 0.0
            else:
                score = 0.5
            sort_pipe.stdin.write(f"[{score:.1f}] {board.fen()}\n".encode())

# Close the pipe to the sort command
sort_pipe.stdin.close()

# Wait for the sort command to finish
sort_pipe.wait()

# Run the uniq and wc commands using the subprocess module and capture their output
uniq_output = subprocess.check_output(["uniq", output_file])
wc_output = subprocess.check_output(["wc", "-l", output_file])

# Print the output of the uniq and wc commands
print(uniq_output.decode())
print(wc_output.decode().strip())
