package main

import (
	"testing"

	pb "github.com/brotherlogic/dirsync/proto"
	"golang.org/x/net/context"
)

func TestHydrate(t *testing.T) {
	s := InitTestServer()
	s.AddSync(context.Background(), &pb.AddSyncRequest{Sync: &pb.Sync{Dir: "/etc/"}})

	err := s.hydrate(context.Background())
	if err != nil {
		t.Errorf("Bad hydrate: %v", err)
	}

	err = s.hydrate(context.Background())
	if err != nil {
		t.Errorf("Bad hydrate: %v", err)
	}

	config, _ := s.load(context.Background())

	if len(config.GetSyncs()[0].Servers) != 1 {
		t.Errorf("Sync was not hydrated")
	}
}

func TestLoadFailHydrate(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true
	err := s.hydrate(context.Background())
	if err == nil {
		t.Errorf("No error on bad hydrate")
	}
}
