package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// collatzSteps обчислює кількість кроків для виродження числа в 1
func collatzSteps(n int) int {
	steps := 0
	for n != 1 {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}
		steps++
	}
	return steps
}

// worker обробляє числа з каналу та зберігає результати
func worker(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range jobs {
		steps := collatzSteps(num)
		results <- steps
	}
}

func main() {
	const totalNumbers = 10_000_000
	numWorkers := runtime.NumCPU()

	fmt.Printf("Обчислення гіпотези Колаца для %d чисел\n", totalNumbers)
	fmt.Printf("Кількість потоків: %d\n", numWorkers)

	startTime := time.Now()

	jobs := make(chan int, 1000)
	results := make(chan int, 1000)

	var wg sync.WaitGroup

	// Запуск воркерів
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Генерація чисел у головному потоці
	go func() {
		for i := 1; i <= totalNumbers; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Закриття каналу результатів після завершення всіх воркерів
	go func() {
		wg.Wait()
		close(results)
	}()

	// Збір результатів
	totalSteps := 0
	count := 0
	for steps := range results {
		totalSteps += steps
		count++
	}

	elapsed := time.Since(startTime)
	averageSteps := float64(totalSteps) / float64(count)

	fmt.Printf("Середня кількість кроків: %.2f\n", averageSteps)
	fmt.Printf("Час виконання: %v\n", elapsed)
}
