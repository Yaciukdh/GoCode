This is the dining philosopher problem with go routines. There are 3 solutions and the original problem written.
set the solution number to 0 for the original problem, 1 is left hand solution, 2 is arbiter solution, 3 is chandy-ish solution.

There is a debug global variable for additional info to what steps the code is doing, set currently to 0.
Another global variable checkClaim is an additional debug variable for solution 3. 
Delay is there to ensure that there is deadlock in original problem and is set to 1. Set it to something else to remove the delay.
