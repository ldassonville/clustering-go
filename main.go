package main

import (
	"flag"

	"os/signal"
	"syscall"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ldassonville/clustering-go/cluster"
)

var (
	// Node name
	NodeName      = flag.String("name", "Node1", "Node name")
	// Node advertising address (ex : 127.0.0.1:7947)
	AdvertiseAddr = flag.String("aa", "127.0.0.1:7946", "Advertise address")
	// Node binding address. Could be same as AdvertiseAddr (ex : 127.0.0.1:7947)
	BindAddr      = flag.String("ba", "127.0.0.1:7946", "Bind address")
	// Coma separated cluster nodes address ex : (127.0.0.1:7947,127.0.0.1:7946)
	ClusterAddr   = flag.String("ca", "127.0.0.1:7946", "Cluster address")

)



func main(){

	// Get cluster config
	flag.Parse()

	// Initialise serf node
	_, err := cluster.NewNode(*AdvertiseAddr, *ClusterAddr, *BindAddr, *NodeName)
	if err != nil {
		logrus.Fatal(err)
	}


	logrus.Print("Node created " )

	runDeamon()


}


func runDeamon(){

	// Wait sign term for close agent
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done

	fmt.Println("exiting")
}