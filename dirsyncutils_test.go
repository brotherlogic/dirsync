package main

import (
	"testing"
	"time"

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

func TestLoadFailSync(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true
	err := s.runSync(context.Background())
	if err == nil {
		t.Errorf("No error on bad hydrate")
	}
}

func TestBasicSync(t *testing.T) {
	s := InitTestServer()
	s.AddSync(context.Background(), &pb.AddSyncRequest{Sync: &pb.Sync{Dir: "blah", Subdir: "blah2", Servers: []string{"server1", "server2"}, Events: []*pb.SyncEvent{&pb.SyncEvent{EndTime: time.Now().Unix(), FromServer: "server1", ToServer: "server2"}}}})
	s.AddSync(context.Background(), &pb.AddSyncRequest{Sync: &pb.Sync{LastSyncTime: time.Now().Unix(), Dir: "blah2", Subdir: "blah3", Servers: []string{"server1", "server2"}, Events: []*pb.SyncEvent{&pb.SyncEvent{EndTime: time.Now().Unix(), FromServer: "server1", ToServer: "server2"}}}})

	err := s.runSync(context.Background())
	if err != nil {
		t.Errorf("Error in sync: %v", err)
	}
	err = s.runSync(context.Background())
	if err != nil {
		t.Errorf("Error in sync: %v", err)
	}

	if s.lastSync != "server2 -> server1" {
		t.Errorf("Bad sync run %v", s.lastSync)
	}
}

func TestRunSync(t *testing.T) {
	s := InitTestServer()
	s.AddSync(context.Background(), &pb.AddSyncRequest{Sync: &pb.Sync{Dir: "blah", Subdir: "blah2", Servers: []string{"server1", "server2"}, Events: []*pb.SyncEvent{&pb.SyncEvent{EndTime: time.Now().Unix(), FromServer: "server1", ToServer: "server2", State: pb.SyncEvent_IN_SYNC}}}})

	err := s.runSync(context.Background())
	if err != nil {
		t.Errorf("Bad sync run: %v", err)
	}

	if s.lastSync != "" {
		t.Errorf("Ran sync when we shouldn't have: %v", s.lastSync)
	}
}
