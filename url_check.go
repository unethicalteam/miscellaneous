package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

var urls = []string{
	"https://unethical.team",
	"https://unethicalcdn.com/",
}

// Define Maximum Amount of Concurrent Checks
const (
	defaultMaxConcurrentChecks = 30
)

func checkURL(url string, wg *sync.WaitGroup, activeURLs *[]string, mutex *sync.Mutex) {
	defer wg.Done()

	resp, err := http.Head(url)
	if err != nil {
		mutex.Lock()
		fmt.Printf("%s %s(Error: %s)%s\n", url, ColorRed, err, ColorReset)
		mutex.Unlock()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		mutex.Lock()
		fmt.Printf("%s %s(active)%s\n", url, ColorGreen, ColorReset)
		*activeURLs = append(*activeURLs, url)
		mutex.Unlock()
	}
}

func main() {
	startTime := time.Now()
	var wg sync.WaitGroup
	var activeURLs []string
	var mutex sync.Mutex

	maxConcurrentChecks := defaultMaxConcurrentChecks

	if len(os.Args) > 1 {
		arg := os.Args[1]
		maxConcurrentChecks = defaultMaxConcurrentChecks
		n, err := fmt.Sscanf(arg, "%d", &maxConcurrentChecks)
		if n != 1 || err != nil || maxConcurrentChecks <= 0 {
			fmt.Println("Invalid argument for maximum concurrent checks. Using the default value.")
			maxConcurrentChecks = defaultMaxConcurrentChecks
		}
	}

	fmt.Printf("Checking %d URLs with a maximum of %d concurrent checks...\n", len(urls), maxConcurrentChecks)

	semaphore := make(chan struct{}, maxConcurrentChecks)

	for _, url := range urls {
		semaphore <- struct{}{}

		wg.Add(1)
		go func(url string) {
			defer func() { <-semaphore }()
			checkURL(url, &wg, &activeURLs, &mutex)
		}(url)
	}

	wg.Wait()

	fmt.Println("\nActive URLs:")
	fmt.Println("var urls = []string{")
	for _, activeURL := range activeURLs {
		fmt.Printf("\t\"%s\",\n", activeURL)
	}
	fmt.Println("}")

	fmt.Printf("\nURL check completed.\n")
	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time elapsed: %s\n", elapsedTime)
	fmt.Printf("\nSummary: %s%d%s URLs checked, %s%d%s active, %s%d%s inactive.%s\n",
		ColorYellow, len(urls), ColorReset,
		ColorGreen, len(activeURLs), ColorReset,
		ColorRed, len(urls)-len(activeURLs), ColorReset,
		ColorReset)
}