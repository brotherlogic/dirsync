package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/dirsync/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(&utils.DiscoveryClientResolverBuilder{})
}

func main() {
	conn, err := grpc.Dial("discovery:///dirsync", grpc.WithInsecure(), grpc.WithBalancerName("my_pick_first"))
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewDirsyncServiceClient(conn)
	ctx, cancel := utils.BuildContext("dirsync-cli", "dirsync")
	defer cancel()

	switch os.Args[1] {
	case "add":
		addFlags := flag.NewFlagSet("AddRecords", flag.ExitOnError)
		var dir = addFlags.String("dir", "", "Dir to add")
		var subdir = addFlags.String("subdir", "", "Subdir")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			_, err := client.AddSync(ctx, &pb.AddSyncRequest{Sync: &pb.Sync{Dir: *dir, Subdir: *subdir}})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}

	}

}
