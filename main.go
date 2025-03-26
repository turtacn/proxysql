package main

import (
	"log"
	"net"
	"runtime"

	"github.com/go-mysql-org/go-mysql/server"
)

const (
	listenAddress = "127.0.0.1:3306"
	mysqlUser     = "root"
	mysqlPassword = "123" // Consider security implications in a real scenario
)

func main() {
	// Listen for connections on the specified address
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenAddress, err)
	}
	defer listener.Close()

	log.Printf("Listening on %s, connect with 'mysql -h %s -P %d -u %s'",
		listenAddress, net.ParseIP(listenAddress).String(), getPort(listenAddress), mysqlUser)

	// Accept connections in a loop to handle multiple clients concurrently
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue // Continue to accept other connections
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Accepted connection from %s", remoteAddr)

	// Create a connection with the specified user and password.
	// Replace EmptyHandler with your own handler for actual SQL command processing.
	mysqlConn, err := server.NewConn(conn, mysqlUser, mysqlPassword, server.EmptyHandler{})
	if err != nil {
		log.Printf("Failed to register connection from %s with server: %v", remoteAddr, err)
		return
	}

	log.Printf("Registered connection from %s with the server", remoteAddr)

	// Handle commands from the client until an error occurs (e.g., client disconnects)
	for {
		err := mysqlConn.HandleCommand()
		if err != nil {
			log.Printf("Error handling command from %s: %v", remoteAddr, err)
			return // Exit the goroutine when the client connection breaks
		}
	}
}

// Helper function to extract the port from the address string
func getPort(addr string) int {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0 // Or handle the error appropriately
	}
	port, err := net.LookupPort("tcp", portStr)
	if err != nil {
		return 0 // Or handle the error appropriately
	}
	return port
}

func init() {
	// Set the number of Goroutines that can run in parallel.
	// This is optional but can be beneficial for performance.
	runtime.GOMAXPROCS(runtime.NumCPU())
}