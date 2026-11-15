package main

import (
	"fmt"
	"log"
	"net"
)

type BackendRegister struct {
  backend_addresses []string
  last_used_index int
}

var backend_register = &BackendRegister{
    backend_addresses: []string{"127.0.0.1:8081", "127.0.0.1:8082"}, // TODO: add function to accept and add new backends here
    last_used_index: -1,
}

func main() {
	port_listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("Failed to create listener on 127.0.0.1:8080")
	}

	for {
		client_connection, err := port_listener.Accept()
        if err != nil {
            log.Println("failed to accept connection from client", err)
            continue
        }

        backend_conn, err := get_available_backend_connection()

        if err != nil {
          log.Println("No available backend.. closing client connection") // TODO: add reliability architecture later (writing somewhere else, holding queue, etc.)
          client_connection.Close()
          continue
        }

        go handle_client_connection(client_connection, backend_conn)
	}
}

func get_available_backend_connection() (net.Conn, error) {
  if len(backend_register.backend_addresses) == 0 {
        log.Println("Failed to fetch an available backed") // TODO: ideally this should be put in some sort of a holding queue.. like let the connection write to a random buffer.. however tradeoff is that the client assumes that the data is okay to go and may mislead the user (unless the connection status is separate from not being able to write through a passthrough.. so that means that the load balancer itself needs to have some sort of a ticking time bomb attached to each incoming client request in case backends fail)
        return nil, fmt.Errorf("no available backend")
    }

    backend_register.last_used_index = (backend_register.last_used_index + 1) % len(backend_register.backend_addresses)
    available_backend_address := backend_register.backend_addresses[backend_register.last_used_index]

    backend_conn, err := net.Dial("tcp", available_backend_address)
    if err != nil {
        return nil, err
    }

    return backend_conn, nil
}

func handle_client_connection(client_conn net.Conn, backend_conn net.Conn) {
    defer client_conn.Close()

    buf := make([]byte, 4096)

    go handle_backend_response(backend_conn, client_conn)

    for {
        n, err := client_conn.Read(buf)
        if err != nil {
            log.Println("error while reading from client connection: ", err)
            return
        }

        if _, err := backend_conn.Write(buf[:n]); err != nil {
			log.Printf("client→backend: write => %v", err)
			return
		}
    }
}

func handle_backend_response(backend_conn net.Conn, client_conn net.Conn) {
    defer backend_conn.Close()
    defer client_conn.Close()

    buf := make([]byte, 4096)

    for {
      n, err := backend_conn.Read(buf)
      if err != nil {
        log.Println("error while reading from backend connection: ", err)
        return
      }

      _, err = client_conn.Write(buf[:n])
      if err != nil {
        log.Printf("backend->client: write => %v", err)
        return
      }
    }
}
