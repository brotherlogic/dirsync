package main

import (
	"context"
	"testing"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/dirsync/proto"
)

//InitTestServer gets a test version of the server
func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), CONFIG, &pb.Config{})
	return s
}

func TestBasicAdd(t *testing.T) {
	s := InitTestServer()

	s.AddSync(context.Background(), &pb.AddSyncRequest{Sync: &pb.Sync{Dir: "blah", Subdir: "blah2"}})
	s.AddSync(context.Background(), &pb.AddSyncRequest{Sync: &pb.Sync{Dir: "blah", Subdir: "blah2"}})

	config, err := s.load(context.Background())
	if err != nil {
		t.Errorf("Bad load: %v", err)
	}

	if len(config.GetSyncs()) != 1 {
		t.Errorf("Bad config: %v", config)
	}
}

func TestAddFail(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	_, err := s.AddSync(context.Background(), &pb.AddSyncRequest{Sync: &pb.Sync{Dir: "blah", Subdir: "blah2"}})
	if err == nil {
		t.Errorf("Should have failed")
	}
}
