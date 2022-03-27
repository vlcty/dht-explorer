package main

import (
    // "fmt"
    "os"
    "time"
    "net"

    "github.com/anacrolix/dht"
)

func main() {
    conn, _ := net.ListenPacket("udp", "[2001:16e0:300:b002:6e:b668:77fb:68d6]:1234")
	server, err := dht.NewServer(&dht.ServerConfig{
		Conn:          conn,
		StartingNodes: dht.GlobalBootstrapAddrs})

    if err != nil {
        panic(err)
    }

    var infohash [20]byte
	copy(infohash[:], "af5e8f524e1ccfe33b210e285f3009bae9faae41")

    server.Announce(infohash, 0, false)

    for {
        server.WriteStatus(os.Stdout)
        time.Sleep(time.Second * 10)
    }
}
