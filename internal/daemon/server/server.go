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
		s.log.Errorf("failed to set the profile enabled: %s", err)
		return &pb.SetOnResponse{}, fmt.Errorf("failed to set enabled: %s", err)
	}

	if err := s.Inter.Reflect(); err != nil {
		return &pb.SetOnResponse{}, fmt.Errorf("failed to reflect: %s", err)
	}

	s.log.Infof("set the profile enabled: %s", in.Profile)
	return &pb.SetOnResponse{}, nil
}

// SetOff set the profile disabled.
func (s *Server) SetOff(ctx context.Context, in *pb.SetOffRequest) (*pb.SetOffResponse, error) {
	if err := s.Inter.Off(in.Profile); err != nil {
		s.log.Errorf("failed to set the profile disabled: %s", err)
		return &pb.SetOffResponse{}, fmt.Errorf("failed to set enabled: %s", err)
	}

	if err := s.Inter.Reflect(); err != nil {
		return &pb.SetOffResponse{}, fmt.Errorf("failed to reflect: %s", err)
	}

	s.log.Infof("set the profile disabled: %s", in.Profile)
	return &pb.SetOffResponse{}, nil
}

// SetConfig set the configuration of the profile.
func (s *Server) SetConfig(ctx context.Context, in *pb.SetConfigRequest) (*pb.SetConfigResponse, error) {
	config, ok, err := s.Inter.GetConfig(in.Profile)
	if err != nil {
		s.log.Errorf("failed to load the config file: %s", err)
		return &pb.SetConfigResponse{}, fmt.Errorf("failed to load the config file: %s", err)
	}
	if !ok {
		s.log.Warnf("failed to get the config, there is no such a profile: %s", in.Profile)
		return &pb.SetConfigResponse{}, fmt.Errorf("failed to get the config, there is no such a profile: %s", in.Profile)
	}

	config.SerialNumber = in.Serial
	config.DurationSecond = in.Duration

	// set the configuration.
	err = s.Inter.SetConfig(in.Profile, config)
	if err != nil {
		s.log.Errorf("failed to set the configuraiton: %s", err)
		return &pb.SetConfigResponse{}, fmt.Errorf("failed to set the configuraiton: %s", err)
	}

	s.log.Infof("set the config of \"%s\" [serial: \"%s\", duration: \"%d\"]", in.Profile, in.Serial, in.Duration)
	return &pb.SetConfigResponse{}, nil
}

// SetGenerate generate the session token of the profile, and reflect it on the credential file.
func (s *Server) SetGenerate(ctx context.Context, in *pb.SetGenerateRequest) (*pb.SetGenerateResponse, error) {
	err := s.Inter.Gen(in.Profile, in.Token)
	if err != nil {
		s.log.Errorf("failed to generate the session token: %s", err)
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
	profiles := make([]*pb.Profile, 0)
	pl, err := s.Inter.GetProfileList()
	if err != nil {
		return &pb.GetProfileListResponse{Profiles: profiles}, err
	}

	for _, p := range pl {
		profiles = append(profiles, &pb.Profile{
			Name: p.Name,
			On: p.On,
			Serial: p.Serial,
			Duration: p.Duration,
			Expired: p.Expired,
		})
	}
	s.log.Info("list profiles")
	return &pb.GetProfileListResponse{Profiles: profiles}, nil
}
