package main

import (
	"fmt"
	"os"
	"sort"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/dirsync/proto"
)

func (s *Server) runSyncEvent(ctx context.Context, ev *pb.Sync) error {
	//Fill out a blank sync matrix
	syncMatrix := make([][]int64, 0)

	indexMap := make(map[string]int)

	for i := 0; i < len(ev.GetServers()); i++ {
		syncMatrix = append(syncMatrix, make([]int64, 0))
		indexMap[ev.GetServers()[i]] = i
		for j := 0; j < len(ev.GetServers()); j++ {
			syncMatrix[i] = append(syncMatrix[i], int64(0))
		}
	}

	minSync := int64(0)
	for _, sync := range ev.GetEvents() {
		syncMatrix[indexMap[sync.GetFromServer()]][indexMap[sync.GetToServer()]] = sync.GetEndTime()
		if sync.GetEndTime() > minSync {
			minSync = sync.GetEndTime()
		}
	}

	bestI := 0
	bestJ := 1
	for i := 0; i < len(ev.GetServers()); i++ {
		for j := 0; j < len(ev.GetServers()); j++ {
			if i != j {
				if syncMatrix[i][j] < minSync {
					minSync = syncMatrix[i][j]
					bestI = i
					bestJ = j
				}
			}
		}
	}

	s.lastSync = fmt.Sprintf("%v -> %v", ev.GetServers()[bestI], ev.GetServers()[bestJ])
	return nil
}

func (s *Server) runSync(ctx context.Context) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	// Find an existing event and update it
	for _, sync := range config.GetSyncs() {
		for _, event := range sync.GetEvents() {
			if event.GetState() == pb.SyncEvent_IN_SYNC {
				//s.updateSync(sync, event, config)
				return s.save(ctx, config)
			}
		}
	}

	sort.SliceStable(config.GetSyncs(), func(i, j int) bool {
		return config.GetSyncs()[i].LastSyncTime < config.GetSyncs()[j].LastSyncTime
	})

	return s.runSyncEvent(ctx, config.GetSyncs()[0])
}

func (s *Server) hydrate(ctx context.Context) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	change := false
	for _, sync := range config.GetSyncs() {
		_, err := os.Stat(sync.GetDir())
		if err == nil {
			found := false
			for _, server := range sync.GetServers() {
				if server == s.Registry.Identifier {
					found = true
				}
			}

			if !found {
				sync.Servers = append(sync.Servers, s.Registry.Identifier)
				change = true
			}
		}
	}

	if change {
		return s.save(ctx, config)
	}
	return nil
}
