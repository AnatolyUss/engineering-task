# Instructions to run the program:
1. Download and install Go compiler.
2. Git-clone or download this program.
3. Change current directory to the program's directory: <code>$ cd /path/to/program</code>.
4. Build the program: <code>$ go build main.go</code>.
5. Run the program: <code>$ ./main</code>.

# Instructions to use the program:
1. Use following commands:
<code>1</code> then hit ENTER - places a submarine at a random location.
<code>2</code> then hit ENTER - places a destroyer at a random location.
<code>3</code> then hit ENTER - places a cruiser at a random location.
<code>4</code> then hit ENTER - places a carrier at a random location.
<code>h</code> then hit ENTER - prints the commands list.
<code>q</code> then hit ENTER - exists the program.

2. Remarks:
A. Note, each vessel will be placed at a random location, vertically, horizontally or even diagonally.
B. Vessels will NOT be placed adjacent to one another.
C. Vessels coordinates are generated pseudo-randomly, 
hence if appropriate coordinates (vertex) of a vessel were not generated within fixed amount of attempts (200) 
the program will be terminated with following message:
<code>Cannot generate valid position for your vessel.</code>.
