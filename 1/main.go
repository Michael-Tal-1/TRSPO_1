package main

import (
	"fmt"
	"time"
)

func fibonacci(n int, result chan int) {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
		time.Sleep(100 * time.Millisecond) // симулюємо роботу
	}
	fmt.Println("Fibonacci(", n, ") =", a)
	result <- a
}

func main() {
	results := make(chan int, 2)

	go fibonacci(10, results)
	go fibonacci(15, results)

	// Отримуємо результат (блокує).
	fib1 := <-results
	fib2 := <-results

	fmt.Println("Sum =", fib1+fib2)
}
