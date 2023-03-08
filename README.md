# Project 1: Process Scheduler

## Description 
For this project we'll be building a simple process scheduler that takes in a file containing example processes, and outputs a schedule based on the three different schedule types:

- First Come First Serve (FCFS) [already done]
- Shortest Job First (SJF)
- SJF Priority
- Round-robin (RR)
- Assume that all processes are CPU bound (they do not block for I/O).

The scheduler will be written in [Go](https://go.dev/) (a skeleton main.go is included in the project repo).

## Steps

1. Clone down the example input/output and skeleton main.go:

   1. `git clone https://github.com/jh125486/CSCE4600`

2. Copy the `Project1` files to your own git project.

3. The processes for your scheduling algorithms are read from a file as the first argument to your program.

   1. Every line in this file includes a record with comma separated fields.

      1. The format for this record is the following: \<ProcessID>,\<Burst Duration>,\<Arrival Time>,\<Priority>.

   2. Not all fields are used by all scheduling algorithms. For example, for FCFS you only need the process IDs, arrival times, and burst durations.

   3. All processes in your input files will be provided a unique process ID. The arrival times and burst durations are integers. Process priorities have a range of [1-50]; the lower this number, the higher the priority i.e. a process with priority=1 has a higher priority than a process with priority=2.

4. Start editing the `main.go` and add the scheduling algorithms:
   1. Implement SJF (preemptive) and report average turnaround time, average waiting time, and average throughput.

   2. Implement SJF priority scheduling (preemptive) and report average turnaround time, average waiting time, and average throughput.

   3. Round-round (preemptive) and report average turnaround time, average waiting time, and average throughput.

## Grading

Code must compile and run.

Each type is worth different points:

- 30 points for implementing FCFS (already done, so you get 30 points for just submitting)
- 25 points for implementing SJF [preemptive]
- 25 points for implementing SJF Priority scheduling [preemptive]
- 20 points for implementing RR [always preemptive]

## Deliverables

A GitHub link to your project which includes:

- `README.md` <- describes anything needed to build (optional)
- `main.go` <- your scheduler

-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Compiling the code
- The code is generated in Go Language and the code was entirely modified in VS Code. To compile the code,
Do the following steps
-Download the code from the repository or clone it 
- git clone "url"
- Open the clone or download code in VS code
-In terminal where go compiler is installed,
 type "go run main.go example_processes.csv"
- This will generate the output and save the output that is generated in output.txt
-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Description on How the Code Works

1 SJF Schedule
- Shortest Job First(SJF) schedule the process based on the burst time which is the time required to complete the process. We then define BurstTime
and sort key for sorting the process based on the burst duration time. 
- In SJF schedule, we define some new parameters such as remaining, completed, and waiting to keep track of the remaining time, waiting time, and completed time. 
In the loop, the process are scheduled until the process is completed. During the time in the loop for each process, it checks the remaining time, waiting time and updates
the total waiting and turn around time. It will check if the remaining time is 0, the process is marked as completed and the average waiting time, turn around time, and different process are updated and provided in the output. Gant Scheduled is printed and output.txt is generated when the algorithm runs the code.\

=============================================================================================================================================================================

2 SJF Priority Schedule
- Shortest Job First Priority schedule the process based on the priority  which is the time required to complete the process. If the bpth process has same priority, then it is determined by the process burst time.
- We define the priority based on the process which uses len() to get the length of the process and swap method to sort the process based on the priority
- In SJF Priority function, different parameters are defined as in SJF algorithm to keep track of the waiting time, remaing times and updates. It takes in a writer to output the results, a string as the title of the schedule, and a slice of Process structs.
-In the loop, the process are scheduled based on the priority until the process is completed. During the time in the loop for each process, it checks the remaining time, waiting time and updates the total waiting and turn around time. It will check if the remaining time is 0, the process is marked as completed and the average waiting time, turn around time, and different process are updated and provided in the output. Gant Scheduled is printed and output.txt is generated when the algorithm runs the code.\

============================================================================================================================================================================

3 Round Robin Schedule
- Round Robin Schedule schedule the process in a preemetive manner. The function takes in an output writer, a title string, and a slice of processes as input arguments. Different parameters are defined such as waitinf time, remaining rime , totalTime and so on to make sure that the output can be generated correctly. Quantum is also set to 2.\
- In the RoundRobin Function, the function iterates if all the functions have been completed. If both of the conditions are satisfied, the function checks for if the remaing time is greater than quantum time. If it is, thn the quantum time is subtracted and the total time is incremented by the quantum time. If not, process is marked as completed. After that different parameters such as waitingTime, totalTime, remaining times are calculated. Once all the paramters are calculated, output is generated and the function is the terminal and saved in output.txt file