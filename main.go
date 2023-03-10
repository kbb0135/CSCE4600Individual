package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// CLI args
	f, closeFile, err := openProcessingFile(os.Args...)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile()

	// Load and parse processes
	processes, err := loadProcesses(f)
	if err != nil {
		log.Fatal(err)
	}
	// First-come, first-serve scheduling
	FCFSSchedule(os.Stdout, "First-come, first-serve", processes)

	SJFSchedule(os.Stdout, "Shortest-job-first", processes)
	//
	SJFPrioritySchedule(os.Stdout, "Priority", processes)
	//
	RRSchedule(os.Stdout, "Round-robin", processes)

	//Saving the output in output.txt file

}

func openProcessingFile(args ...string) (*os.File, func(), error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf("%w: must give a scheduling file to process", ErrInvalidArgs)
	}
	// Read in CSV process CSV file
	f, err := os.Open(args[1])
	if err != nil {
		return nil, nil, fmt.Errorf("%v: error opening scheduling file", err)
	}
	closeFn := func() {
		if err := f.Close(); err != nil {
			log.Fatalf("%v: error closing scheduling file", err)
		}
	}

	return f, closeFn, nil
}

type (
	Process struct {
		ProcessID     int64
		ArrivalTime   int64
		BurstDuration int64
		Priority      int64
	}
	TimeSlice struct {
		PID   int64
		Start int64
		Stop  int64
	}
)
type BurstTime []Process
//sorting the process based on burst duration
func (a BurstTime) Len() int           { return len(a) }
func (a BurstTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BurstTime) Less(i, j int) bool { return a[i].BurstDuration < a[j].BurstDuration }

//Sorting the process based on priority
type Priority []Process

func (p Priority) Len() int {
	return (len(p))
}
func (p Priority) Less(i, j int) bool {
	return p[i].Priority < p[j].Priority
}
func (p Priority) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFSchedule(w io.Writer, title string, processes []Process) {
	//Implementing the process based on Burst Duration 
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
		remaining       = make([]int64, len(processes))
		completed       = make([]bool, len(processes))
		waiting         = make([]int64, len(processes))
		noComplete      int
		totalTime       int64
	)
	/*
	In the loop we check if the process is completed or not.
	The process is then scheduled based on the burst duration. After that the 
	process is put in a order and the time waiting, remaing time, and completed time are 
	calculated. The remaining time will check if the process is completed and the value
	are then printed.
	*/
	for i := range processes {
		remaining[i] = processes[i].BurstDuration
	}
	for noComplete < len(processes) {
		sort.Sort(BurstTime(processes))
		for i := 0; i < len(processes); i++ {
			p := processes[i]
			if !completed[i] {
				if remaining[i] == 0 {
					completed[i] = true
					noComplete++
					totalTime += p.BurstDuration
					totalTurnaround = float64(totalTime)
					waiting[i] = totalTime - p.BurstDuration
					waitingTime = waiting[i]
					fmt.Println("Total Wait Check =", waiting[i])
					totalWait += float64(waitingTime)
				} else {
					remaining[i]--
					totalTime++
				}
			}
			completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
			lastCompletion = float64(completion)
			start := waitingTime + processes[i].ArrivalTime
			schedule[i] = []string{
				fmt.Sprint(processes[i].ProcessID),
				fmt.Sprint(processes[i].Priority),
				fmt.Sprint(processes[i].BurstDuration),
				fmt.Sprint(processes[i].ArrivalTime),
				fmt.Sprint(waitingTime),
				fmt.Sprint(totalTurnaround),
				fmt.Sprint(completion),
			}
			serviceTime += processes[i].BurstDuration
			gantt = append(gantt, TimeSlice{
				PID:   processes[i].ProcessID,
				Start: start,
				Stop:  serviceTime,
			})
		}

	}
	count := float64(len(processes))
	fmt.Println("Count and Total Wait")
	fmt.Println(count)
	totalWait = float64(waiting[0]+waiting[1]+waiting[2])
	fmt.Println(totalWait)
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)

}

func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
	//schedule the process and updates based on the priority of the schedule
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
		remaining       = make([]int64, len(processes))	//to process the reamining time
		completed       = make([]bool, len(processes))	//for making sure if the process is completed
		waiting         = make([]int64, len(processes))	//to store the waiting time for each processes
		noComplete      int
		totalTime       int64
	)
	/*
	In the loop we check if the process is completed or not.
	The process is then scheduled based on the priority. After that the 
	process is put in a order and the time waiting, remaing time, and completed time are 
	calculated. The remaining time will check if the process is completed and the value
	are then printed.
	*/

	for i := 0; i < len(processes); i++ {
		remaining[i] = processes[i].BurstDuration
	}
	//Check if all the processes are completed
	for noComplete < len(processes) {
		sort.Sort(Priority(processes))
		for i := 0; i < len(processes); i++ {
			p := processes[i]
			if !completed[i] {
				if remaining[i] == 0 {
					completed[i] = true
					noComplete++
					totalTime += p.BurstDuration
					totalTurnaround = float64(totalTime)
					waiting[i] = totalTime - p.BurstDuration
					waitingTime = waiting[i]
					totalWait += float64(waitingTime)
				} else {
					remaining[i]--
					totalTime++
				}
			}
			completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
			lastCompletion = float64(completion)
			start := waitingTime + processes[i].ArrivalTime
			if i < len(schedule) {
				schedule[i] = []string{
					fmt.Sprint(processes[i].ProcessID),
					fmt.Sprint(processes[i].Priority),
					fmt.Sprint(processes[i].BurstDuration),
					fmt.Sprint(processes[i].ArrivalTime),
					fmt.Sprint(waitingTime),
					fmt.Sprint(totalTurnaround),
					fmt.Sprint(completion),
				}
			}

			serviceTime += processes[i].BurstDuration
			gantt = append(gantt, TimeSlice{
				PID:   processes[i].ProcessID,
				Start: start,
				Stop:  serviceTime,
			})
		}
	}
	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}
/*
RRSchedule function schedules the processes in a preemptive manner and 
outputs the schedule, gantt chart, and summary statistics of the processes 
to the output writer.
*/
func RRSchedule(w io.Writer, title string, processes []Process) {
	//defining different parameters to schedule, queue, and storing and 
	//calculating the waiting time, remainingtime, no of Completete process
	//and so on

	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
		remaining       = make([]int64, len(processes))
		waiting         = make([]int64, len(processes))
		noComplete      int
		totalTime       int64
	)
	var quantum int64 = 2 //Defining Quantum time for Round Robin Algorithm
	//fmt.Print("Entering quantum time);
	//fmt.Scanln(&quantum)
	fmt.Println("Quantum time set as :  ", quantum)
	//Loop Until if all the processes has been completed
	for i := range processes {
		remaining[i] = processes[i].BurstDuration
	}
	//For each process, check if it has arrived and 
	//still has remaining burst time
	for noComplete < len(processes) {
		for i := 0; i < len(processes) && i < len(remaining); i++ {
			p := &processes[i]
			if p.ArrivalTime > totalTime || remaining[i] == 0 {
				continue
			}
	//if the reamining time is greater than quantum time, we subtract the 
	//quantum time and totalTime on quantumTime else we will update the 
	//total time, remaining time, waitingTime and totalTurnAround Time are
	//updated and the total completion is completed
			if remaining[i] > quantum {
				remaining[i] -= quantum
				totalTime += quantum
			} else {
				totalTime += remaining[i]
				remaining[i] = 0
				noComplete++
				waiting[i] = totalTime - p.BurstDuration - p.ArrivalTime
				waitingTime = waiting[i]
				totalTurnaround = float64(totalTime)
				totalWait += float64(waitingTime)
			}
			completion := totalTime
			lastCompletion = float64(completion)
			start := p.ArrivalTime + waitingTime
			schedule[i] = []string{
				fmt.Sprint(processes[i].ProcessID),
				fmt.Sprint(processes[i].Priority),
				fmt.Sprint(processes[i].BurstDuration),
				fmt.Sprint(processes[i].ArrivalTime),
				fmt.Sprint(waitingTime),
				fmt.Sprint(totalTurnaround),
				fmt.Sprint(completion),
			}
			serviceTime += p.BurstDuration
			gantt = append(gantt, TimeSlice{
				PID:   processes[i].ProcessID,
				Start: start,
				Stop:  serviceTime,
			})
		}
	}
	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

//endregion

//region Output helpers

func outputTitle(w io.Writer, title string) {
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
	_, _ = fmt.Fprintln(w, strings.Repeat(" ", len(title)/2), title)
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
}

func outputGantt(w io.Writer, gantt []TimeSlice) {
	_, _ = fmt.Fprintln(w, "Gantt schedule")
	_, _ = fmt.Fprint(w, "|")
	for i := range gantt {
		pid := fmt.Sprint(gantt[i].PID)
		padding := strings.Repeat(" ", (8-len(pid))/2)
		_, _ = fmt.Fprint(w, padding, pid, padding, "|")
	}
	_, _ = fmt.Fprintln(w)
	for i := range gantt {
		_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Start), "\t")
		if len(gantt)-1 == i {
			_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Stop))
		}
	}
	_, _ = fmt.Fprintf(w, "\n\n")
}

func outputSchedule(w io.Writer, rows [][]string, wait, turnaround, throughput float64) {
	_, _ = fmt.Fprintln(w, "Schedule table")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"ID", "Priority", "Burst", "Arrival", "Wait", "Turnaround", "Exit"})
	table.AppendBulk(rows)
	table.SetFooter([]string{"", "", "", "",
		fmt.Sprintf("Average\n%.2f", wait),
		fmt.Sprintf("Average\n%.2f", turnaround),
		fmt.Sprintf("Throughput\n%.2f/t", throughput)})
	table.Render()
}

//endregion

//region Loading processes.

var ErrInvalidArgs = errors.New("invalid args")

func loadProcesses(r io.Reader) ([]Process, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: reading CSV", err)
	}

	processes := make([]Process, len(rows))
	for i := range rows {
		processes[i].ProcessID = mustStrToInt(rows[i][0])
		processes[i].BurstDuration = mustStrToInt(rows[i][1])
		processes[i].ArrivalTime = mustStrToInt(rows[i][2])
		if len(rows[i]) == 4 {
			processes[i].Priority = mustStrToInt(rows[i][3])
		}
	}

	return processes, nil
}

func mustStrToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return i
}
//endregion
