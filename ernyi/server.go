package ernyi

import (
   "net"
   "fmt"
   "os"
   "log"
   "bufio"
)

func StartServer(addr string) {
  	l, err := net.Listen("tcp", addr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    log.Printf(fmt.Sprintf("Listening on %s", addr))
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

func handleRequest(conn net.Conn) {
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Message Received:", string(message))
    conn.Close()
}