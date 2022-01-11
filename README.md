# wordle-solve (or YAWS Yet-Another-Wordle-Solver)

## Usage

Define your guesses...
```
raise
.   .
```

(space = grey, `.` (dot) = yellow, anything else = green)

and then `cat guesses | go run main.go` to receive your next best guess word!

`outer`

## Test your guesses

Given the following guesses file
```
raise
. x
print
 xxx
bring
 xxx
```

then `go run main.go -t drink guesses` prints

<p style="font-family: 'Courier New'; background-color: black; color:white">
raise OK <b><span style="background-color:#b59f3b">r</span><span style="background-color:#333">a</span><span style="background-color: green">i</span><span style="background-color: #333">s</span><span style="background-color: #333">e</span><br/></b>
print OK <b><span style="background-color:#333">p</span><span style="background-color:green">r</span><span style="background-color: green">i</span><span style="background-color: green">n</span><span style="background-color: #333">t</span><br/></b>
bring OK <b><span style="background-color:#333">b</span><span style="background-color:green">r</span><span style="background-color: green">i</span><span style="background-color: green">n</span><span style="background-color: #333">g</span><br/></b>
</p>
