package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/practice/futures"
)

const (
	ClientErrorClose = "client error closed"
	ClientErrorType  = "client error expecting client type"
	ClientErrorTimeout  = "accept tcp [::]:8080: i/o timeout"
)

type Server struct {
	Listener net.Listener
	Clients  []*Client
	Future   futures.Future
}

type Client struct {
	Connection net.Conn
	Future     futures.Future
	Buffer     []byte
	I, N       int
}

func client_async_read_once(client *Client) futures.Future {
	return futures.Poll(func(state any) (futures.Result, error) {
		c, ok := state.(*Client)
		if !ok {
			return futures.Result{}, errors.New(ClientErrorType)
		}
		// Read method will block until a connection is made
		// We need to set a deadline to avoid blocking forever
		// Then reset the deadline after the accept call
		c.Connection.SetDeadline(time.Now().Add(10 * time.Millisecond))
		len, err := c.Connection.Read(c.Buffer)
		c.Connection.SetDeadline(time.Time{})
		if err != nil {
			return futures.Result{}, err
		}

		if len == 0 {
			return futures.Result{}, errors.New(ClientErrorClose)
		}
		c.N = len
		return futures.Finished(c), nil
	}, client)
}

func client_async_write_everything(client *Client) futures.Future {
	return futures.Poll(func(state any) (futures.Result, error) {
		c, ok := state.(*Client)
		if !ok {
			return futures.Result{}, errors.New(ClientErrorType)
		}
		len, err := c.Connection.Write(c.Buffer[c.I:c.N])
		if err != nil {
			return futures.Result{}, err
		}

		if len == 0 {
			return futures.Result{}, errors.New(ClientErrorClose)
		}

		c.I += len
		c.N -= len

		fmt.Println("Write", c.I, c.N)
		if c.N == 0 {
			return futures.Finished(client), nil
		}
		return futures.Pending(), nil
	}, client)
}

func client_async_catch_error(err error) futures.Future {
	if err.Error() == ClientErrorClose {
		return futures.Resolve(nil)
	}

	return futures.Reject(err)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	clients := []*Client{}
	clients_next := []*Client{}

	server := Server{
		Listener: listener,
		Clients:  clients,
	}

	server.Future = futures.Poll(func(state any) (futures.Result, error) {
		s, ok := state.(*Server)
		if !ok {
			return futures.Result{}, errors.New(ClientErrorType)
		}
		// Accept method will block until a connection is made
		// We need to set a deadline to avoid blocking forever
		// Then reset the deadline after the accept call
		s.Listener.(*net.TCPListener).SetDeadline(time.Now().Add(10 * time.Millisecond))
		conn, err := s.Listener.Accept()
		s.Listener.(*net.TCPListener).SetDeadline(time.Time{})
		if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return futures.Pending(), nil
		}
			return futures.Result{}, err
		}

		client := &Client{
			Connection: conn,
			Buffer:     make([]byte, 1024),
		}
		client.Future = futures.Forever(func(state any) futures.Future {
			c, ok := state.(*Client)
			if !ok {
				return futures.Reject(errors.New(ClientErrorType))
			}
			c.I = 0
			c.N = 0
			return client_async_read_once(c).Then(func(reslut any) futures.Future {
				c, ok := reslut.(*Client)
				if !ok {
					return futures.Reject(errors.New(ClientErrorType))
				}
				return client_async_write_everything(c)
			}).Catch(client_async_catch_error)
		}, client)
		server.Clients = append(server.Clients, client)
		clients = server.Clients
		fmt.Println("Client connected")
		return futures.Pending(), nil
	}, &server)

	for {
		_, err := server.Future.Poll()
		if err != nil {
			panic(err)
		}

		for _, c := range server.Clients {
			res, _ := c.Future.Poll()
			if !res.Finished {
				clients_next = append(clients_next, c)
			} else {
				fmt.Println("Client disconnected")
				c.Connection.Close()
			}
		}

		// Swap clients and clients_next
		clients = clients[:0]
		clients = append(clients, clients_next...)
		clients_next = clients_next[:0]
	}
}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	buf := make([]byte, 1024)
// 	fmt.Println("New connection from", conn.RemoteAddr())

// 	for {
// 		len, err := conn.Read(buf)
// 		if err != nil {
// 			fmt.Printf("Error reading: %#v\n", err)
// 			return
// 		}

// 		if len == 0 {
// 			break
// 		}

// 		fmt.Printf("Message received: %s\n", string(buf[:len]))

// 		conn.Write(buf[:len])
// 	}
// }
