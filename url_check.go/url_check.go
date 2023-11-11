package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
    "time"
)

// constants for output formatting
const (
    colorRed    = "\033[31m"
    colorGreen  = "\033[32m"
    colorReset  = "\033[0m"
)

// global variable for max retries
var maxRetries int

// readURLsFromFile reads urls from a json file. it's simple: give it a file path, and it gives you urls.
func readURLsFromFile(filePath string) ([]string, error) {
    var urls []string
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &urls)
    return urls, err
}

// checkURL checks if a url is active. it tries a few times (based on maxRetries) before giving up.
func checkURL(url string, wg *sync.WaitGroup, activeURLs *[]string, mutex *sync.Mutex) {
    defer wg.Done()

    for attempt := 0; attempt < maxRetries; attempt++ {
        resp, err := http.Head(url)
        if err == nil {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusOK {
                mutex.Lock()
                fmt.Printf("%s %s(active)%s\n", url, colorGreen, colorReset)
                *activeURLs = append(*activeURLs, url)
                mutex.Unlock()
                return
            }
        }
    }

    mutex.Lock()
    fmt.Printf("%s %s(dead)%s\n", url, colorRed, colorReset)
    mutex.Unlock()
}

// handleSingleURL is a helper to keep things neat. it just manages one url check.
func handleSingleURL(url string, semaphore chan struct{}, wg *sync.WaitGroup, activeURLs *[]string, mutex *sync.Mutex) {
    defer func() { <-semaphore }() // this line makes sure we play nice with others
    checkURL(url, wg, activeURLs, mutex)
}

// processURLs is where the magic happens. it checks all the urls, but not all at once.
func processURLs(urls []string, maxConcurrentChecks int) []string {
    var wg sync.WaitGroup
    var activeURLs []string
    var mutex sync.Mutex

    semaphore := make(chan struct{}, maxConcurrentChecks)

    for _, url := range urls {
        semaphore <- struct{}{} // got a permit? go ahead.
        wg.Add(1)
        go handleSingleURL(url, semaphore, &wg, &activeURLs, &mutex)
    }

    wg.Wait() // hold on till everyone's done
    return activeURLs
}

// writeURLsToFile writes urls to a json file. nothing fancy, just saving stuff.
func writeURLsToFile(urls []string, filePath string) error {
    data, err := json.MarshalIndent(urls, "", "    ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filePath, data, 0644)
}

func main() {
    // let's grab some flags from the command line
    maxConcurrentChecksFlag := flag.Int("concurrent", 10, "how many urls we check at once")
    maxRetriesFlag := flag.Int("retries", 3, "how many times we try a url before giving up")
    urlsFileFlag := flag.String("urlsFile", "urls.json", "where your urls are stored")

    flag.Parse()

    // using the flags
    maxConcurrentChecks := *maxConcurrentChecksFlag
    maxRetries = *maxRetriesFlag
    urlsFile := *urlsFileFlag

    // starting the clock
    startTime := time.Now()

    // read urls from the file
    urls, err := readURLsFromFile(urlsFile)
    if err != nil {
        fmt.Println("whoops! couldn't read the urls:", err)
        return
    }

    // let's get to work!
    activeURLs := processURLs(urls, maxConcurrentChecks)

    // saving the active urls back to the file
    err = writeURLsToFile(activeURLs, urlsFile)
    if err != nil {
        fmt.Println("hmm, couldn't write the active urls:", err)
        return
    }

    // all done, let's see how long it took
    elapsedTime := time.Since(startTime)
    fmt.Printf("all done! it took %s. active urls are now saved.\n", elapsedTime)
}
