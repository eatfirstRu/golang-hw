package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	scanner *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		scanner: bufio.NewScanner(in),
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect to %s: %w", c.address, err)
	}
	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	if !c.scanner.Scan() {
		return io.EOF
	}
	_, err := c.conn.Write(append(c.scanner.Bytes(), '\n'))
	return err
}

func (c *telnetClient) Receive() error {
	buf := make([]byte, 4096)
	n, err := c.conn.Read(buf)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buf[:n])
	return err
}
