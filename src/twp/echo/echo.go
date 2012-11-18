package main

import (
    "fmt"
    "net"
    "twp"
)

var (
    host string = "www.dcl.hpi.uni-potsdam.de"
    port string = "80"
    // host string = "localhost"
    // port string = "8001"
)

func main() {
    target := fmt.Sprintf("%s:%s", host, port)    
    conn, err := net.Dial("tcp", target)
    if err != nil {
        fmt.Println("Error: ", err)
    }
    defer conn.Close()
    fmt.Println("Conn: ", conn)

    rd := twp.NewReader(conn)
    wr := twp.NewWriter(conn)
    
    if err := wr.InitWithProtocol(2); err != nil {
        fmt.Println(err)
    }
    if _, err := wr.Write([]byte("\x04\x1eHello, World!\x00")); err != nil {
        panic(err)
    }
    wr.Flush()

    fmt.Println("Sent")

    msg, err := rd.ReadValue(); 
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
    fmt.Println(msg)
}
