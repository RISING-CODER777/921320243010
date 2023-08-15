package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type NumberResponse struct {
	Numbers []int `json:"numbers"`
}

func fetchNumbersFromURL(url string, ch chan<- []int) {
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error fetching data from", url, err)
		ch <- nil
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body from", url, err)
		ch <- nil
		return
	}

	var numResponse NumberResponse
	err = json.Unmarshal(body, &numResponse)
	if err != nil {
		fmt.Println("Error unmarshaling response from", url, err)
		ch <- nil
		return
	}

	ch <- numResponse.Numbers
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j <= high-1; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func quickSort(arr []int, low, high int) {
	if low < high {
		pivot := partition(arr, low, high)
		quickSort(arr, low, pivot-1)
		quickSort(arr, pivot+1, high)
	}
}

func mergeUniqueNumbers(numbersList ...[]int) []int {
	uniqueNumbers := make(map[int]bool)
	merged := []int{}

	for _, numbers := range numbersList {
		for _, num := range numbers {
			if !uniqueNumbers[num] {
				uniqueNumbers[num] = true
				merged = append(merged, num)
			}
		}
	}

	quickSort(merged, 0, len(merged)-1)

	return merged
}

func getMergedNumbersFromURLs(urls []string) []int {
	var wg sync.WaitGroup
	ch := make(chan []int, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			fetchNumbersFromURL(u, ch)
		}(url)
	}

	wg.Wait()
	close(ch)

	numbersList := []int{}
	for nums := range ch {
		if nums != nil {
			numbersList = append(numbersList, nums...)
		}
	}

	return mergeUniqueNumbers(numbersList)
}

func NumbersHandler(w http.ResponseWriter, r *http.Request) {
	urls, ok := r.URL.Query()["url"]
	if !ok {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	var validURLs []string
	for _, url := range urls {
		_, err := http.NewRequest(http.MethodGet, url, nil)
		if err == nil {
			validURLs = append(validURLs, url)
		}
	}

	startTime := time.Now()
	mergedNumbers := getMergedNumbersFromURLs(validURLs)
	endTime := time.Now()

	response := NumberResponse{
		Numbers: mergedNumbers,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
	fmt.Println("Time taken:", endTime.Sub(startTime))
}

func main() {
	http.HandleFunc("/numbers", NumbersHandler)

	port := ":3000"
	fmt.Println("Server is listening on port", port)
	http.ListenAndServe(port, nil)
}
