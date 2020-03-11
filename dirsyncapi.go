package main

import "golang.org/x/net/context"
import pb "github.com/brotherlogic/dirsync/proto"

//AddSync adds a sync dir
func (s *Server) AddSync(ctx context.Context, req *pb.AddSyncRequest) (*pb.AddSyncResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	for _, sync := range config.GetSyncs() {
		if sync.GetDir() == req.GetSync().GetDir() &&
			sync.GetSubdir() == req.GetSync().GetSubdir() {
			return &pb.AddSyncResponse{}, nil
		}
	}

	config.Syncs = append(config.GetSyncs(), req.GetSync())

	return &pb.AddSyncResponse{}, s.save(ctx, config)
}
