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

    de := twp.NewDecoder(conn)
    en := twp.NewEncoder(conn)
    
    if err := en.InitWithProtocol(2); err != nil {
        fmt.Println(err)
    }
    if _, err := en.Write([]byte("\x04")); err != nil {
        panic(err)
    }
    if err := en.EncodeString("Hello, World!"); err != nil {
        panic(err)
    }
    if _, err := en.Write([]byte{0}); err != nil {
        panic(err)
    }
    // wr.Flush()

    fmt.Println("Sent")

    var msg *twp.Message
    if msg, err = de.DecodeMessage(); err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    fmt.Printf("%v", msg)
}
