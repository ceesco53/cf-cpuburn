package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	numBurn        int
	updateInterval int
	timeout        int
)


// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
}

func cpuBurn() {
	for {
		for i := 0; i < 2147483647; i++ {
		}
		runtime.Gosched()
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Something hit the health check endpoint")
	w.Header().Set("cpuburn", "cpuburn health check")
	w.WriteHeader(200)
}

func init() {
	flag.IntVar(&numBurn, "n", 0, "number of cores to burn (0 = all)")
	flag.IntVar(&updateInterval, "u", 10, "seconds between updates (0 = don't update)")
	flag.Parse()
	if numBurn <= 0 {
		numBurn = runtime.NumCPU()
	}
}

func main() {
	timeout := getEnvAsInt("CPUBURN_TIMEOUT", 600)
	if timeout > 600 {
		timeout = 600 // limit timeout of cpuburn in seconds
	}

	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8080"
	}
	go func() {
		http.HandleFunc("/", handler)
		http.HandleFunc("/health", healthCheck)
		fmt.Println("Listening on port", port)
		http.ListenAndServe(":"+port, nil)
	}()

	ch := make(chan bool, 1)
	defer close(ch)

	go func() {
		ch <- true
	}()

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	runtime.GOMAXPROCS(numBurn)
	fmt.Printf("Burning %d CPUs/cores\n", numBurn)
	for i := 0; i < numBurn; i++ {
		go cpuBurn()
	}
	if updateInterval > 0 {
		t := time.Tick(time.Duration(updateInterval) * time.Second)
		for secs := updateInterval; ; secs += updateInterval {
			<-t
			fmt.Printf("%d seconds\n", secs)
		}
	} else {
		select {} // wait forever
	}
}

