package main

import (
	"fmt"
	"log"
	"net"
)

type BackendRegister struct {
  backendAddresses []string
  lastUsedIndex int
}

var backendRegister = &BackendRegister{
    backendAddresses: []string{"127.0.0.1:8081", "127.0.0.1:8082"}, // TODO: add function to accept and add new backends here
    lastUsedIndex: -1,
}

func main() {
	portListener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("Failed to create listener on 127.0.0.1:8080")
	}

	for {
		clientConnection, err := portListener.Accept()
        if err != nil {
            log.Println("failed to accept connection from client", err)
            continue
        }

        backendConn, err := getAvailableBackendConnection()

        if err != nil {
          log.Println("No available backend.. closing client connection") // TODO: add reliability architecture later (writing somewhere else, holding queue, etc.)
          clientConnection.Close()
          continue
        }

        go handleClientConnection(clientConnection, backendConn)
	}
}

func getAvailableBackendConnection() (net.Conn, error) {
  if len(backendRegister.backendAddresses) == 0 {
        log.Println("Failed to fetch an available backed") // TODO: ideally this should be put in some sort of a holding queue.. like let the connection write to a random buffer.. however tradeoff is that the client assumes that the data is okay to go and may mislead the user (unless the connection status is separate from not being able to write through a passthrough.. so that means that the load balancer itself needs to have some sort of a ticking time bomb attached to each incoming client request in case backends fail)
        return nil, fmt.Errorf("no available backend")
    }

    backendRegister.lastUsedIndex = (backendRegister.lastUsedIndex + 1) % len(backendRegister.backendAddresses)
    availableBackendAddress := backendRegister.backendAddresses[backendRegister.lastUsedIndex]

    backendConn, err := net.Dial("tcp", availableBackendAddress)
    if err != nil {
        return nil, err
    }

    return backendConn, nil
}

func handleClientConnection(clientConn net.Conn, backendConn net.Conn) {
    defer clientConn.Close()

    buf := make([]byte, 4096)

    go handleBackendResponse(backendConn, clientConn)

    for {
        n, err := clientConn.Read(buf)
        if err != nil {
            log.Println("error while reading from client connection: ", err)
            return
        }

        if _, err := backendConn.Write(buf[:n]); err != nil {
			log.Printf("client→backend: write => %v", err)
			return
		}
    }
}

func handleBackendResponse(backendConn net.Conn, clientConn net.Conn) {
    defer backendConn.Close()
    defer clientConn.Close()

    buf := make([]byte, 4096)

    for {
      n, err := backendConn.Read(buf)
      if err != nil {
        log.Println("error while reading from backend connection: ", err)
        return
      }

      _, err = clientConn.Write(buf[:n])
      if err != nil {
        log.Printf("backend->client: write => %v", err)
        return
      }
    }
}
