package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/sacloud/iaas-api-go"
	sacloudServer "github.com/sacloud/iaas-service-go/server"
	sacloudLocalSwitch "github.com/sacloud/iaas-service-go/swytch"
)

var worker = workerpool.New(50)

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
	switchService := sacloudLocalSwitch.New(client)

	// FindServers
	serverFindRequest := sacloudServer.FindRequest{
		Zone: zone,
		Tags: args,
	}
	servers, err := serverService.Find(&serverFindRequest)
	if err != nil {
		panic(err)
	}
	startTime := time.Now()
	svc := len(servers)
	diskc := 0
	for _, v := range servers {
		serverDeleteWithDisk := sacloudServer.DeleteRequest{
			Zone:      zone,
			ID:        v.ID,
			WithDisks: true,
			Force:     true,
		}
		diskc += len(v.Disks)
		worker.Submit(func() {
			err := serverService.Delete(&serverDeleteWithDisk)
			if err != nil {
				log.Println(err)
			}
		})
	}

	findSwitchReq := sacloudLocalSwitch.FindRequest{
		Zone: zone,
	}
	sws, err := switchService.Find(&findSwitchReq)
	if err != nil {
		panic(err)
	}
	for _, v := range sws {
		// Switchに接続しているサーバが存在する場合
		if v.GetServerCount() != 0 {
			continue
		}
		deleteSwitchReq := sacloudLocalSwitch.DeleteRequest{
			Zone: zone,
			ID:   v.ID,
		}
		worker.Submit(func() {
			err := switchService.Delete(&deleteSwitchReq)
			if err != nil {
				log.Println(err)
			}
		})
	}

	worker.StopWait()

	str := fmt.Sprintf("Servers: %d/ Disks: %d/ 実行時間: %f 分", svc, diskc, time.Now().Sub(startTime).Minutes())
	log.Println(str)
}
