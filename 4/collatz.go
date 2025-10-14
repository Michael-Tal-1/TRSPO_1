package main

import (
	"fmt"
	"runtime"
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

// workerSequential обробляє послідовний діапазон чисел
func workerSequential(start, end int, result *int64, done chan struct{}) {
	localSum := int64(0)
	for i := start; i <= end; i++ {
		steps := collatzSteps(i)
		localSum += int64(steps)
	}
	*result = localSum // Кожен воркер пише у свій слот - немає race condition
	done <- struct{}{}
}

// workerInterleaved обробляє числа з інтерлівінгом
func workerInterleaved(workerID, totalNumbers, numWorkers int, result *int64, done chan struct{}) {
	localSum := int64(0)
	for i := workerID + 1; i <= totalNumbers; i += numWorkers {
		steps := collatzSteps(i)
		localSum += int64(steps)
	}
	*result = localSum // Кожен воркер пише у свій слот - немає race condition
	done <- struct{}{}
}

func runSequential(totalNumbers int) (float64, time.Duration) {
	numWorkers := runtime.NumCPU()
	startTime := time.Now()

	results := make([]int64, numWorkers) // Кожен воркер має свій слот
	done := make(chan struct{}, numWorkers)

	chunkSize := totalNumbers / numWorkers
	for i := 0; i < numWorkers; i++ {
		start := i*chunkSize + 1
		end := (i + 1) * chunkSize
		if i == numWorkers-1 {
			end = totalNumbers
		}
		go workerSequential(start, end, &results[i], done)
	}

	for i := 0; i < numWorkers; i++ {
		<-done
	}

	// Підсумовуємо результати після завершення всіх воркерів
	var totalSteps int64
	for _, val := range results {
		totalSteps += val
	}

	elapsed := time.Since(startTime)
	averageSteps := float64(totalSteps) / float64(totalNumbers)
	return averageSteps, elapsed
}

func runInterleaved(totalNumbers int) (float64, time.Duration) {
	numWorkers := runtime.NumCPU()
	startTime := time.Now()

	results := make([]int64, numWorkers) // Кожен воркер має свій слот
	done := make(chan struct{}, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go workerInterleaved(i, totalNumbers, numWorkers, &results[i], done)
	}

	for i := 0; i < numWorkers; i++ {
		<-done
	}

	// Підсумовуємо результати після завершення всіх воркерів
	var totalSteps int64
	for _, val := range results {
		totalSteps += val
	}

	elapsed := time.Since(startTime)
	averageSteps := float64(totalSteps) / float64(totalNumbers)
	return averageSteps, elapsed
}

func workerWithSync(jobs <-chan int, results chan<- int, done chan struct{}) {
	for num := range jobs {
		steps := collatzSteps(num)
		results <- steps
	}
	done <- struct{}{}
}

func runWithSync(totalNumbers int) (float64, time.Duration) {
	numWorkers := runtime.NumCPU()

	startTime := time.Now()

	jobs := make(chan int, 1000)
	results := make(chan int, 1000)
	done := make(chan struct{}, numWorkers)

	// Запуск воркерів
	for i := 0; i < numWorkers; i++ {
		go workerWithSync(jobs, results, done)
	}

	// Генерація чисел
	go func() {
		for i := 1; i <= totalNumbers; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Закриття каналу результатів після завершення всіх воркерів
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
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

	return averageSteps, elapsed
}

func main() {
	const totalNumbers = 10_000_000
	numWorkers := runtime.NumCPU()

	fmt.Printf("Обчислення гіпотези Колаца для %d чисел\n", totalNumbers)
	fmt.Printf("Кількість потоків: %d\n\n", numWorkers)

	// 1. Sequential (лінійний розподіл)
	fmt.Println("1. Sequential (лінійний розподіл)")
	avgSteps1, elapsed1 := runSequential(totalNumbers)
	fmt.Printf("Середня кількість кроків: %.2f\n", avgSteps1)
	fmt.Printf("Час виконання: %v\n\n", elapsed1)

	// 2. Interleaved (інтерлівінг)
	fmt.Println("2. Interleaved (інтерлівінг)")
	avgSteps2, elapsed2 := runInterleaved(totalNumbers)
	fmt.Printf("Середня кількість кроків: %.2f\n", avgSteps2)
	fmt.Printf("Час виконання: %v\n\n", elapsed2)

	// 3. З каналами (попередня версія)
	fmt.Println("3. З каналами (попередня версія)")
	avgSteps4, elapsed3 := runWithSync(totalNumbers)
	fmt.Printf("Середня кількість кроків: %.2f\n", avgSteps4)
	fmt.Printf("Час виконання: %v\n\n", elapsed3)

	// Порівняння
	fmt.Println("Порівняння")
	fmt.Printf("Sequential:   %v (базовий)\n", elapsed1)
	fmt.Printf("Interleaved:  %v (%.2f%% від базового)\n", elapsed2, float64(elapsed2)/float64(elapsed1)*100)
	fmt.Printf("Channels:     %v (%.2f%% від базового)\n\n", elapsed3, float64(elapsed3)/float64(elapsed1)*100)

	// Найшвидша версія
	fastest := elapsed1
	fastestName := "Sequential"
	if elapsed2 < fastest {
		fastest = elapsed2
		fastestName = "Interleaved"
	}
	if elapsed3 < fastest {
		fastest = elapsed3
		fastestName = "Channels"
	}
	fmt.Printf("Найшвидша версія: %s (%v)\n", fastestName, fastest)
}
