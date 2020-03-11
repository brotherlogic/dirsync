package main

import (
	"os"

	"golang.org/x/net/context"
)

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
