package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func progUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [Флаги] <IP-адрес> <порт>\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, "Флаги:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = progUsage
	timeoutStr := ""
	flag.StringVar(&timeoutStr, "timeout", "", "Задержка при подключении к серверу (последовательность десятичных чисел, "+
		"где у каждого может быть необязательная дробная часть и суффикс единицы измерения: "+
		"ns — наносекунды; us или µs — микросекунды; ms — миллисекунды; s — секунды; m — минуты; h — часы)")

	// Парсим флаги из командной строки
	flag.Parse()

	// Получаем оставшиеся (после обработки флагов) аргументы — это позиционные
	args := flag.Args()

	// Проверяем, что ровно два позиционных аргумента были переданы
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Ошибка: должно быть ровно два позиционных аргумента.")
		flag.Usage()
		os.Exit(1)
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil || timeout <= 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: неверное значение флага timeout")
		flag.Usage()
		os.Exit(1)
	}

	port, err := strconv.Atoi(args[1])
	if err != nil || port <= 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: неверное значение  для номера порта")
		flag.Usage()
		os.Exit(1)
	}

	client := NewTelnetClient(args[0]+":"+strconv.Itoa(port), timeout, os.Stdin, os.Stdout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	err = client.Connect(ctx)
	if err != nil {
		client.Close()
		return
	}

	defer client.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	wg.Go(func() {
		err := client.Receive(ctx)
		if err != nil {
			stop()
		}
	})
	wg.Go(func() {
		err := client.Send(ctx)
		if err != nil {
			stop()
		}
	})

	wg.Wait()
}
