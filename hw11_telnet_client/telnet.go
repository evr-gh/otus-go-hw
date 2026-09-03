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
	Connect(ctx context.Context) error
	io.Closer
	Send(ctx context.Context) error
	Receive(ctx context.Context) error
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

func (c *Client) Connect(ctx context.Context) error {
	var dialer net.Dialer
	cntx := ctx
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

func (c *Client) Send(ctx context.Context) error {
	if c.Connection == nil {
		err := ErrSendNoConnection
		fmt.Fprintf(os.Stderr, "Oшибка передачи данных: %v\n", err)
		return err
	}

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

	for {
		select {
		case <-ctx.Done():
			return nil
		case str, ok := <-inChan:
			if !ok {
				return nil
			}
			if _, err := fmt.Fprintf(c.Connection, "%s\n", str); err != nil {
				fmt.Fprintf(os.Stderr, "Oшибка передачи данных: %v\n", err)
				return err
			}
		}
	}
}

func (c *Client) Receive(ctx context.Context) error {
	if c.Connection == nil {
		err := ErrReceiveNoConnection
		fmt.Fprintf(os.Stderr, "Oшибка приема данных: %v\n", err)
		return err
	}

	scanner := bufio.NewScanner(c.Connection)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				err := scanner.Err()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Oшибка приема данных: %v\n", err)
				}
				return err
			}
			text := scanner.Text()
			fmt.Fprintf(c.Out, "%s\n", text)
		}
	}
}
