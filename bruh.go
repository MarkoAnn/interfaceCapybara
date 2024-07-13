package main

import (
	"fmt"
	"time"
)

// Определяем интерфейс Printer
type Printer interface {
	Print(data string)
}


type ConsolePrinter struct{}

func (cp ConsolePrinter) Print(data string) {
	fmt.Println(data)
}


type TimestampPrinter struct{}

func (tp TimestampPrinter) Print(data string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s: %s\n", currentTime, data)
}

func sayMyName(p Printer) {
	p.Print("My name is Capybara")
}

func main() {
	// Создаем переменные для каждой из реализаций интерфейса Printer
	var consolePrinter Printer = ConsolePrinter{}
	var timestampPrinter Printer = TimestampPrinter{}

	// Используем метод Print для каждой из реализаций
	consolePrinter.Print("Hello, Golang!")
	timestampPrinter.Print("Hello, Golang with timestamp!")
	
	// Вызываем функцию sayMyName для каждой реализации
	sayMyName(consolePrinter)
	sayMyName(timestampPrinter)
}
