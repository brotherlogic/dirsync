package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"

	pb "github.com/brotherlogic/dirsync/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

func init() {
	resolver.Register(&utils.DiscoveryServerResolverBuilder{})
}

const (
	// CONFIG - Where we store syncs
	CONFIG = "/github.com/brotherlogic/dirsync/config"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	syncs []*pb.Sync
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
	}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterDirsyncServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

//Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

func (s *Server) load(ctx context.Context) (*pb.Config, error) {
	data, _, err := s.KSclient.Read(ctx, CONFIG, &pb.Config{})
	if err != nil {
		return nil, err
	}
	return data.(*pb.Config), nil
}

func (s *Server) save(ctx context.Context, config *pb.Config) error {
	return s.KSclient.Save(ctx, CONFIG, config)
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "syncs", Text: fmt.Sprintf("%v", s.syncs)},
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("dirsync", false, true)
	if err != nil {
		return
	}

	fmt.Printf("%v", server.Serve())
}
