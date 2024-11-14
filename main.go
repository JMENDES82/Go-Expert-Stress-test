package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Err        error
}

func worker(url string, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for range jobs {
		resp, err := http.Get(url)
		if err != nil {
			results <- Result{StatusCode: 0, Err: err}
			continue
		}
		resp.Body.Close()
		results <- Result{StatusCode: resp.StatusCode}
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	totalRequests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("Por favor, forneça uma URL válida usando --url")
		return
	}

	if *totalRequests <= 0 || *concurrency <= 0 {
		fmt.Println("Por favor, forneça valores positivos para --requests e --concurrency")
		return
	}

	startTime := time.Now()

	jobs := make(chan int)
	results := make(chan Result)
	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(*url, jobs, results, &wg)
	}

	var resultWg sync.WaitGroup
	resultWg.Add(1)
	var mu sync.Mutex
	total200 := 0
	statusCodes := make(map[int]int)

	go func() {
		defer resultWg.Done()
		for res := range results {
			mu.Lock()
			statusCodes[res.StatusCode]++
			if res.StatusCode == 200 {
				total200++
			}
			mu.Unlock()
		}
	}()

	go func() {
		for i := 0; i < *totalRequests; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	wg.Wait()
	close(results)

	resultWg.Wait()

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	fmt.Println("===== Relatório de Teste =====")
	fmt.Printf("Tempo total gasto: %v\n", totalTime)
	fmt.Printf("Total de requests realizados: %d\n", *totalRequests)
	fmt.Printf("Total de requests com HTTP 200: %d\n", total200)
	fmt.Println("Distribuição de códigos de status:")
	for code, count := range statusCodes {
		fmt.Printf("Status %d: %d\n", code, count)
	}
}
