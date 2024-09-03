## Current Implementation CPU profile
This addresses the initial implementation referenced through this
[git commit](https://github.com/ercross/wheel/commit/47898660878d8355bcdb1880362d0ff2c70360a8)

See [cpu profile dump image](pprof001.svg)

## Nailing one of the problems (goroutine starvation)
The initial implementation cpu profile shows that there’s an inefficient cpu utilization, 
as the runtime function call stack was consuming most of my execution time. 
Here’s the invocation stack trace:
```
runtime.mcall -> runtime.park_m -> runtime.schedule -> runtime.findRunnable -> 
 -> runtime.stealWork -> runtime.runqsteal -> runtime.runqgrab -> runtime.usleep
```
This stack trace indicates that a significant amount of CPU time is being spent on the Go runtime's scheduler, 
particularly in the process of finding work for goroutines to execute.
To simplify, Go runtime is spending a lot of time trying to find work for the CPU to do, 
but it's often failing to find runnable goroutines and ends up idling or sleeping.
The culprit that came to my mind was no other than
Goroutine Starvation
That’s the monster I’ve built and I’m in the process of putting down the monster by optimizing the program.
Essentially, there’s a goroutine starvation because of an inefficient work/job distribution 
causing some processors to be underutilized while others are overworked.

## My proposed solution - Reduce invocation of runtime scheduler through Worker Pool
In the current implementation, the scheduler is invoked quite a lot, 
leading to overhead incurring in finding runnable goroutines, balancing workload, and handling context switches. 
A worker pool might help with predictable scheduling. 
With a fixed number of goroutines, the scheduler has fewer decisions to make regarding which goroutine to run next. 
This predictability reduces the need for the scheduler to frequently intervene, leading to fewer invocations.
Asides a controlled goroutine management through predictable scheduling, other ways the worker pool might improve this program are:
1. Reduced Work-Stealing: the runtime.stealWork function is often invoked when there are idle processors with no runnable goroutines.   
   A worker pool would ensure that each worker is busy or waiting on a well-defined queue, 
   reducing the need for the scheduler to steal work across processors.
2. Minimized Idle Time: Worker pools keep the processors engaged with a steady stream of tasks.  
   This minimizes the time the scheduler spends idling (runtime.usleep), as the workers are constantly pulling jobs from the queue, keeping the system busy.