package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)


func ValidateListUrls() {
	
	
	http.HandleFunc("/validateurls", validateURLs)

	http.ListenAndServe(":8080", nil)

	
}

type Urls struct {
	Urls []string `json:"urls"`
}

func validateURLs(w http.ResponseWriter, r *http.Request) {
	byteData, err := io.ReadAll(r.Body)
	urlsData := Urls{}
	if err != nil {
		fmt.Fprintln(w, "Error while reading body")
		return 
	}
	err = json.Unmarshal(byteData, &urlsData)

	if err != nil {
		fmt.Fprintln(w, "Error while unmarshalling data")
		return 
	}

	var wg sync.WaitGroup
	invalidURLsChan := make(chan string, len(urlsData.Urls))
	for _, urlStr := range urlsData.Urls {
		wg.Add(1)
		go parseRequestUrl(urlStr, invalidURLsChan, &wg)
	}
	wg.Wait()
	close(invalidURLsChan)
	for invalidURL := range invalidURLsChan {
		fmt.Printf("Channel Invalid URL : %v \n", invalidURL)
	}
	
	fmt.Fprintf(w, "Validated all urls")
}

func parseRequestUrl(rawURL string, invalidURLsChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Parse the URL
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		// Return false if the URL is invalid
		fmt.Printf("Error : Invalid URL %v \n", rawURL)
		invalidURLsChan <- rawURL
		return
		//return false
	}

	resp, err := http.Get(rawURL)

	if err != nil {
		// Return false if the URL is invalid
		fmt.Printf("Error : Failed to get the request %v \n", rawURL)
		invalidURLsChan <- rawURL
		return
		//return false
	}

	// Return true if the URL is valid
	fmt.Printf("Success: Valid URL %v \n", rawURL)

	if resp.StatusCode != 200 {
		
		fmt.Printf("Error : Did not receive http 200 response %v \n", rawURL)
		invalidURLsChan <- rawURL
		return
	}

	// bodyByte, err := io.ReadAll(resp.Body)

	// if err != nil {
	// 	// Return false if the URL is invalid
	// 	fmt.Printf("Error : Failed to read the body %v \n", rawURL)
	// 	return
	// 	//return false
	// }

	// fmt.Printf(" Success: url : %v Able to get response : %v \n ", rawURL, string(bodyByte))
	
}
