package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sort"
	"sync"
	"time"

	"github.com/fatih/color"
)

// worker function that attempts to connect to a given port and sends the result back
func worker(address string, ports, results chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, p))
		if err != nil {
			results <- 0
		} else {
			conn.Close()
			results <- p
		}
	}
}

func main() {
	// Get address from command line argument
	address := flag.String("address", "scanme.nmap.org", "Address to scan")
	flag.Parse()

	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int
	var wg sync.WaitGroup

	// Start the worker goroutines
	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go worker(*address, ports, results, &wg)
	}

	// Loading message
	go func() {
		color.Set(color.FgGreen)
		defer color.Unset()

		loadingChars := `-\|/`
		i := 0
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				fmt.Printf("\rScanning ports %c", loadingChars[i%len(loadingChars)])
				i++
			}
		}
	}()

	// Send port numbers to the ports channel
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
		close(ports)
	}()

	// Capture interrupt signal to handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		close(results)
		fmt.Println("\nScan interrupted")
		os.Exit(1)
	}()

	// Collect the results
	go func() {
		for port := range results {
			if port != 0 {
				openports = append(openports, port)
			}
		}
	}()

	wg.Wait()
	close(results)

	sort.Ints(openports)
	fmt.Println("\nScan complete. Open ports:")
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
