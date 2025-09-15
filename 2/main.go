package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const totalPoints = 1000000

func calculatePiParallel(points, goroutines int) float64 {
	pointsPerGoroutine := points / goroutines
	results := make(chan int, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			inside := 0
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for j := 0; j < pointsPerGoroutine; j++ {
				x := r.Float64()*2 - 1
				y := r.Float64()*2 - 1
				if x*x+y*y <= 1 {
					inside++
				}
			}
			results <- inside
		}()
	}

	totalInside := 0
	for i := 0; i < goroutines; i++ {
		totalInside += <-results
	}

	return 4.0 * float64(totalInside) / float64(points)
}

func benchmarkCalculation(points, goroutines int) (float64, time.Duration) {
	start := time.Now()
	var pi float64

	pi = calculatePiParallel(points, goroutines)

	duration := time.Since(start)
	return pi, duration
}

func main() {
	goroutineCounts := []int{1, 2, 4, 8, 16, 32, 64, 256, 1024} // Додав 256 задля кращої демонстрації.

	file, err := os.Create("monte_carlo_pi_report.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("Monte Carlo Pi Calculation Performance Report\n")
	file.WriteString("-------------------------------------------\n\n")
	file.WriteString(fmt.Sprintf("Total points: %d\n\n", totalPoints))
	file.WriteString(fmt.Sprintf("%10s %12s %10s %8s\n", "Goroutines", "Pi Value", "Time (ms)", "Speedup"))
	file.WriteString("-------------------------------------------\n")

	var baseTime time.Duration

	for i, count := range goroutineCounts {
		pi, duration := benchmarkCalculation(totalPoints, count)

		if i == 0 {
			baseTime = duration
		}

		speedup := float64(baseTime) / float64(duration)

		result := fmt.Sprintf("%10d %12.6f %10.4f %8.2fx\n",
			count, pi, float64(duration.Nanoseconds())/1e6, speedup)

		fmt.Print(result)
		file.WriteString(result)
	}

	fmt.Println("\nReport saved to monte_carlo_pi_report.txt")
}
