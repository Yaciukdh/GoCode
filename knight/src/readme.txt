Knowing that Go is a great language for concurrency, I wanted to translate this classic (typically in c) programming program 
into Go.

In C, it's common to choose a move then create a thread or fork then pass a board struct holding the previously travelled spaces
and the space where the piece is at which point then iterates through the next move. The boards are then returned to main,
all successful boards are tallied up, and the results are printed.

This is my goal for this program, to successfuly translate this into Go code.
