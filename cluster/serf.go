package cluster

import (
	"strings"
	"net"
	"strconv"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
	"github.com/hashicorp/raft"
)

type SerfNode struct {
	//discovery *discovery.Discover
	cluster   *serf.Serf
	eventChan chan serf.Event
	stopChan  chan struct{}
	log       *logrus.Logger
}

func NewNode(advertiseAddr, clusterAddr, bindAddr, name string) (*SerfNode, error) {

	node := &SerfNode{
		log:       logrus.New(),
		stopChan:  make(chan struct{}),
		eventChan: make(chan serf.Event, 100),
	}

	node.initSerf(advertiseAddr, clusterAddr, bindAddr, name)

	return node, nil
}


func parseAddr(addr string) (string, int, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return host, portInt, nil
}



func (n *SerfNode) initSerf(advertiseAddr, clusterAddr, bindAddr, name string) (error) {

	conf := serf.DefaultConfig()
	conf.Init()
	conf.EventCh = n.eventChan

	if len(name) > 0 {
		conf.NodeName = name
	}
	advertiseHost, advertisePort, err := parseAddr(advertiseAddr)
	if err != nil {
		return err
	}
	conf.Tags = map[string]string{
		"name": name,
	}
	conf.MemberlistConfig.AdvertisePort = advertisePort
	conf.MemberlistConfig.AdvertiseAddr = advertiseHost
	bindHost, bindPort, err := parseAddr(bindAddr)
	if err != nil {
		return err
	}
	conf.MemberlistConfig.BindPort = bindPort
	conf.MemberlistConfig.BindAddr = bindHost
	cluster, err := serf.Create(conf)
	if err != nil {
		return err
	}
	clusterAddresses := strings.Split(clusterAddr, ",")
	_, err = cluster.Join(clusterAddresses, true)
	if err != nil {
		logrus.Warning("Couldn't join cluster, starting own: %v", err)
	}

	n.cluster = cluster
	return nil
}


func (n *SerfNode) ProcessEvents(raft  *raft.Raft ){

	go n.processClusterEvents()
}



func (n *SerfNode) processClusterEvents() {

	for {
		select {

		case event := <-n.eventChan:

			if serf.EventMemberJoin == event.EventType(){

				memberEvent := event.(serf.MemberEvent)
				n.log.Infof("Event receive events %s for members %s",  memberEvent.Type, memberEvent.Members)
			}

		case <-n.stopChan:
			n.log.Info("Stopping cluster event process")
			return
		}
	}
}
