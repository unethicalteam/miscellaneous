package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
    "time"
)

const (
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorReset  = "\033[0m"
    MaxRetries  = 3 // Maximum number of retries for checking a URL
    maxConcurrentChecks = 30 // Maximum URLs to check at once
)

func readURLsFromFile(filePath string) ([]string, error) {
    var urls []string
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &urls)
    return urls, err
}

func checkURL(url string, wg *sync.WaitGroup, activeURLs *[]string, mutex *sync.Mutex) {
    defer wg.Done()

    for attempts := 0; attempts < MaxRetries; attempts++ {
        resp, err := http.Head(url)
        if err == nil {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusOK {
                mutex.Lock()
                fmt.Printf("%s %s(active)%s\n", url, ColorGreen, ColorReset)
                *activeURLs = append(*activeURLs, url)
                mutex.Unlock()
                return
            }
        }
        time.Sleep(time.Second)
    }

    mutex.Lock()
    fmt.Printf("%s %s(dead)%s\n", url, ColorRed, ColorReset)
    mutex.Unlock()
}

func processURLs(urls []string, maxConcurrentChecks int) []string {
    var wg sync.WaitGroup
    var activeURLs []string
    var mutex sync.Mutex

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
    return activeURLs
}


func writeURLsToFile(urls []string, filePath string) error {
    data, err := json.MarshalIndent(urls, "", "    ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filePath, data, 0644)
}

func main() {
    startTime := time.Now()

    urls, err := readURLsFromFile("urls.json")
    if err != nil {
        fmt.Println("Error reading URLs:", err)
        return
    }

    activeURLs := processURLs(urls, maxConcurrentChecks)

    err = writeURLsToFile(activeURLs, "urls.json")
    if err != nil {
        fmt.Println("Error writing active URLs:", err)
        return
    }

    elapsedTime := time.Since(startTime)
    fmt.Printf("\nURL check completed in %s. Active URLs updated.\n", elapsedTime)
}
