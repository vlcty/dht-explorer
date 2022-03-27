package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"time"

	_ "github.com/anacrolix/envpprof"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/dht/krpc"
)

var (
	serveAddr = flag.String("serveAddr", ":0", "local UDP address")
	infoHash  = flag.String("infoHash", "", "torrent infohash")
	maxNodes  = flag.Int("maxNodes", 200, "max amount of nodes to query")

	server       *dht.Server
	counter = 0
)

func printTable(seen map[string]krpc.NodeAddr) {
	fmt.Printf("\n\nPeers associated with that torrent:\n")

	peers := make([]string, 0, len(seen))
	for _, peer := range seen {
		peers = append(peers, peer.String())
	}
	sort.Strings(peers)

	for _, peer := range peers {
		fmt.Printf("%s\n", peer)
	}

	fmt.Printf("\nTotal: %d\n", len(peers))
}

func main() {
	flag.Parse()
	switch len(*infoHash) {
	case 20:
	case 40:
		_, err := fmt.Sscanf(*infoHash, "%x", infoHash)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("require 20 or 40 byte infohash")
	}
	conn, err := net.ListenPacket("udp", *serveAddr)
	if err != nil {
		fmt.Println(err)
	}
	sc := dht.ServerConfig{
		Conn:          conn,
		StartingNodes: dht.GlobalBootstrapAddrs}

	server, err = dht.NewServer(&sc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("dht server on %s, ID is %x\n", server.Addr(), server.ID())

	seen := make(map[string]krpc.NodeAddr)

	var infohash [20]byte
	copy(infohash[:], *infoHash)
	ps, err := server.Announce(infohash, 0, false)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second * 5)

	for {
		v, ok := <-ps.Peers

		counter++

		if !ok {
			fmt.Println("Failed to read peers")
			return
		}

		if counter%20 == 0 {
			fmt.Print("\r")
			fmt.Printf("Progress: %d / %d nodes discovered. Found: %d", counter, *maxNodes, len(seen))
		}

		for _, p := range v.Peers {
			if _, ok := seen[p.IP.String()]; ok {
				continue
			} else {
				seen[p.IP.String()] = p
			}
		}

		if counter >= *maxNodes {
			printTable(seen)
			return
		}
	}
}
