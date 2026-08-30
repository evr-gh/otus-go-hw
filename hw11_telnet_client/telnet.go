package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	ErrSendNoConnection    = errors.New("нет связи")
	ErrReceiveNoConnection = errors.New("нет связи")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	Address    string
	Timeout    time.Duration
	In         io.ReadCloser
	Out        io.Writer
	Connection net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		Address:    address,
		Timeout:    timeout,
		In:         in,
		Out:        out,
		Connection: nil,
	}
}

func (c *Client) Connect() error {
	var dialer net.Dialer
	cntx := context.Background()
	var cancel context.CancelFunc
	if c.Timeout > 0 {
		cntx, cancel = context.WithTimeout(cntx, c.Timeout)
		defer cancel()
	}
	connection, err := dialer.DialContext(cntx, "tcp", c.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось подключиться к серверу: %v\n", err)
		return err
	}
	c.Connection = connection
	return nil
}

func (c *Client) Close() error {
	defer func() {
		err := c.In.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка при закрытии входного потока: %v\n", err)
		}
	}()

	if c.Connection == nil {
		return nil
	}
	err := c.Connection.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось разорвать связь: %v\n", err)
		if c.Connection != nil {
			c.Connection = nil
			fmt.Fprintf(os.Stderr, "Принудительное отсоединение от %s\n", c.Address)
		}
		return err
	}

	return nil
}

func (c *Client) Send() error {
	if c.Connection == nil {
		err := ErrSendNoConnection
		fmt.Fprintf(os.Stderr, "Oшибка передачи данных: %v\n", err)
		return err
	}

	cntx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inChan := make(chan string)
	go func() {
		defer close(inChan)
		scanner := bufio.NewScanner(c.In)
		for scanner.Scan() {
			txt := scanner.Text()
			inChan <- txt
		}
		if scanner.Err() != nil {
			fmt.Fprintf(os.Stderr, "Ошибка при получении данных: %v\n", scanner.Err())
		}
	}()
OUTER:
	for {
		select {
		case <-cntx.Done():
			break OUTER
		case str, ok := <-inChan:
			if !ok {
				break OUTER
			}
			if _, err := fmt.Fprintf(c.Connection, "%s\n", str); err != nil {
				fmt.Fprintf(os.Stderr, "Oшибка передачи данных: %v\n", err)
				return err
			}
		}
	}

	return nil
}

func (c *Client) Receive() error {
	if c.Connection == nil {
		err := ErrReceiveNoConnection
		fmt.Fprintf(os.Stderr, "Oшибка приема данных: %v\n", err)
		return err
	}

	cntx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := bufio.NewScanner(c.Connection)
OUTER:
	for {
		select {
		case <-cntx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			text := scanner.Text()
			fmt.Fprintf(c.Out, "%s\n", text)
		}
	}

	return nil
}
