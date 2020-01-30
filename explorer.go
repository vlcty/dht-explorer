package main

import (
	"flag"
	"fmt"
	"net"
	"time"
	"sort"

	_ "github.com/anacrolix/envpprof"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/dht/krpc"
)

var (
	serveAddr = flag.String("serveAddr", ":0", "local UDP address")
	infoHash  = flag.String("infoHash", "", "torrent infohash")

	s       *dht.Server
	max     = 200
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
		fmt.Println("require 20 byte infohash")
	}
	conn, err := net.ListenPacket("udp", *serveAddr)
	if err != nil {
		fmt.Println(err)
	}
	sc := dht.ServerConfig{
		Conn:          conn,
		StartingNodes: dht.GlobalBootstrapAddrs}

	s, err = dht.NewServer(&sc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("dht server on %s, ID is %x\n", s.Addr(), s.ID())

	seen := make(map[string]krpc.NodeAddr)

	var ih [20]byte
	copy(ih[:], *infoHash)
	ps, err := s.Announce(ih, 0, false)
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

		if counter % 10 == 0 {
			fmt.Print(".")
		}

		if len(v.Peers) == 0 {
			continue
		}

		fmt.Printf("Received %d peers from %x\n", len(v.Peers), v.NodeInfo.ID)

		for _, p := range v.Peers {
			if _, ok := seen[p.String()]; ok {
				continue
			} else {
				seen[p.String()] = p
			}
		}

		if counter >= max {
			printTable(seen)
			return
		}
	}
}
