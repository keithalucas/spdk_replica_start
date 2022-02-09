package main

import (
	"fmt"
	"net"
	"os"

	"github.com/keithalucas/jsonrpc/pkg/jsonrpc"
	"github.com/keithalucas/jsonrpc/pkg/spdk"
)

func main() {

	conn, err := net.Dial("unix", "/var/tmp/spdk.sock")

	if err != nil {
		fmt.Printf("Error opening socket: %v", err)
	}

	client := jsonrpc.NewClient(conn)

	errChan := client.Init()

	tcpServer := spdk.NewTcpJsonServer(os.Args[1], 4421)
	client.SendMsg(tcpServer.GetMethod(), tcpServer)

	externalAddr := spdk.NewLonghornSetExternalAddress(os.Args[1])
	client.SendMsg(externalAddr.GetMethod(), externalAddr)

	aio := spdk.NewAioCreate("sata1", "/dev/sda", 4096)
	client.SendMsg(aio.GetMethod(), aio)

	lvs := spdk.NewBdevLvolCreateLvstore("sata1", "longhorn")
	client.SendMsg(lvs.GetMethod(), lvs)

	lrc := spdk.NewLonghornCreateReplica("demo", 4*1024*1024*1024, "longhorn", os.Args[1], 4420)
	client.SendMsg(lrc.GetMethod(), lrc)

	client.SendMsg("bdev_get_bdevs", nil)

	<-errChan
}
