package server

import (
	"context"
	"fmt"

	pb "github.com/hanjunlee/awscred/api"
	"github.com/sirupsen/logrus"
)

// Server interact with the client to manage and reflect the credentials file.
type Server struct {
	pb.UnimplementedAWSCredServer
	Inter *Interactor
	log   *logrus.Entry
}

// NewServer create a new server.
func NewServer(i *Interactor) *Server {
	return &Server{
		Inter: i,
		// TODO: write logs.
		log: logrus.NewEntry(logrus.New()),
	}
}

// SetOn set the profile enabled. It makes to reflect the session token on the credential file.
func (s *Server) SetOn(ctx context.Context, in *pb.SetOnRequest) (*pb.SetOnResponse, error) {
	if err := s.Inter.On(in.Profile); err != nil {
		return &pb.SetOnResponse{}, fmt.Errorf("failed to set enabled: %s", err)
	}

	if err := s.Inter.Reflect(); err != nil {
		return &pb.SetOnResponse{}, fmt.Errorf("failed to reflect: %s", err)
	}

	return &pb.SetOnResponse{}, nil
}

// SetOff set the profile disabled.
func (s *Server) SetOff(ctx context.Context, in *pb.SetOffRequest) (*pb.SetOffResponse, error) {
	if err := s.Inter.On(in.Profile); err != nil {
		return &pb.SetOffResponse{}, fmt.Errorf("failed to set enabled: %s", err)
	}

	if err := s.Inter.Reflect(); err != nil {
		return &pb.SetOffResponse{}, fmt.Errorf("failed to reflect: %s", err)
	}

	return &pb.SetOffResponse{}, nil
}

// SetConfig set the configuration of the profile.
func (s *Server) SetConfig(ctx context.Context, in *pb.SetConfigRequest) (*pb.SetConfigResponse, error) {
	config, ok, err := s.Inter.GetConfig(in.Profile)
	if err != nil {
		return &pb.SetConfigResponse{}, fmt.Errorf("failed to get the configuration: %s", err)
	}
	if !ok {
		return &pb.SetConfigResponse{}, fmt.Errorf("failed to set the configuration, set enabled the configuration first")
	}

	config.SerialNumber = in.Serial
	config.DurationSecond = in.Duration

	// set the configuration.
	err = s.Inter.SetConfig(in.Profile, config)
	if err != nil {
		return &pb.SetConfigResponse{}, fmt.Errorf("failed to set the configuraiton: %s", err)
	}

	return &pb.SetConfigResponse{}, nil
}

// SetGenerate generate the session token of the profile, and reflect it on the credential file.
func (s *Server) SetGenerate(ctx context.Context, in *pb.SetGenerateRequest) (*pb.SetGenerateResponse, error) {
	err := s.Inter.Gen(in.Profile, in.Token)
	if err != nil {
		return &pb.SetGenerateResponse{}, fmt.Errorf("failed to generate the session token: %s", err)
	}

	err = s.Inter.Reflect()
	if err != nil {
		return &pb.SetGenerateResponse{}, fmt.Errorf("failed to reflect: %s", err)
	}

	return &pb.SetGenerateResponse{}, nil
}

// GetProfileList return the information of profiles.
func (s *Server) GetProfileList(ctx context.Context, in *pb.GetProfileListRequest) (*pb.GetProfileListResponse, error) {
	// TODO: write the code.
	return &pb.GetProfileListResponse{Profiles: []*pb.Profile{}}, nil
}
