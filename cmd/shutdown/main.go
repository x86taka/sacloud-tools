package main

import (
	"flag"
	"log"
	"os"

	"github.com/sacloud/iaas-api-go"
	sacloudServer "github.com/sacloud/iaas-service-go/server"
)

func main() {

	flag.Parse()
	args := flag.Args()

	// APIキー
	token := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	secret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")
	zone := os.Getenv("SAKURACLOUD_ZONE")

	// クライアントの作成
	client := iaas.NewClient(token, secret)

	// Service

	serverService := sacloudServer.New(client)

	// FindServers
	serverFindRequest := sacloudServer.FindRequest{
		Zone: zone,
		Tags: args,
	}
	servers, err := serverService.Find(&serverFindRequest)
	if err != nil {
		panic(err)
	}

	for _, v := range servers {
		if v.InstanceStatus.IsDown() {
			continue
		}
		serverShutdownReq := sacloudServer.ShutdownRequest{
			Zone:          zone,
			ID:            v.ID,
			NoWait:        true,
			ForceShutdown: false,
		}
		err := serverService.Shutdown(&serverShutdownReq)
		if err != nil {
			log.Println(err)
		}
	}
}
