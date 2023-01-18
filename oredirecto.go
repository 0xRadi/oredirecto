package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		originalURL := scanner.Text()
		wg.Add(1)
		go func() {
			fuzzedURLs := fuzzURL(originalURL)
			for _, url := range fuzzedURLs {
				statusCode := getStatusCode(url)
				fmt.Printf("%s %d\n", url, statusCode)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func fuzzURL(originalURL string) []string {
	u, err := url.Parse(originalURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil
	}

	queryParams := u.Query()
	var fuzzedURLs []string
	fuzzValues := []string{"FUZZ", "FUZZ2", "FUZZ3"}
	for _, fuzzValue := range fuzzValues {
		for key, value := range queryParams {
			originalValue := value[0]
			queryParams.Set(key, fuzzValue)
			u.RawQuery = queryParams.Encode()
			fuzzedURLs = append(fuzzedURLs, u.String())
			queryParams.Set(key, originalValue)
		}
	}
	return fuzzedURLs
}

func getStatusCode(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error requesting URL:", err)
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
