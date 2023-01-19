package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
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
				resp := getURL(url)
				if resp != nil {
					if containsCanary(resp) != "" {
						fmt.Println(containsCanary(resp) + url)
						break
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func fuzzURL(originalURL string) []string {
	u, err := url.Parse(originalURL)
	if err != nil {
		//fmt.Println("Error parsing URL:", err)
		return nil
	}

	queryParams := u.Query()
	var fuzzedURLs []string
	fuzzValues := []string{
		"injection:/here",
		"injection:here",
		"http://ahmedradi.com",
		"http://" + u.Hostname() + ".ahmedradi.com"}
	for _, fuzzValue := range fuzzValues {
		for key, value := range queryParams {
			originalValue := value[0]
			queryParams.Set(key, fuzzValue)
			u.RawQuery = queryParams.Encode()
			//fmt.Println(u.String())
			fuzzedURLs = append(fuzzedURLs, u.String())
			queryParams.Set(key, originalValue)
		}
	}
	return fuzzedURLs
}

func getURL(url string) *http.Response {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		//fmt.Println("Error requesting URL:", err)
		return nil
	}

	return resp
}

func containsCanary(resp *http.Response) string {
	re_header := regexp.MustCompile("(?i)^https?:\\/\\/(?:[^\\/]+\\.)?ahmedradi\\.com(\\/.*)?$")
	re_body := regexp.MustCompile("(?i)=[ ]?['\"]https?:\\/\\/(?:[^\\/=\\?]+\\.)?ahmedradi\\.com(\\/.*)?['\"]|REDIRECT313|['\"]injection:[/]?here['\"]")
	found := ""
	// check the headers
	for _, headers := range resp.Header {
		for _, h := range headers {
			if re_header.MatchString(h) {
				found := re_header.FindString(h)
				result := "[Found] [Header] [" + found + "] "
				return result
			}
		}
	}
	// check the body
	body, _ := ioutil.ReadAll(resp.Body)
	if re_body.Match(body) {
		found := re_body.FindString(string(body))
		result := "[Potential] [Body] [" + found + "] "
		fmt.Print(result)
		return found
	}
	return found
}
