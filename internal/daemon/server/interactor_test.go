package server

import (
	"context"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/mock/gomock"
	"github.com/hanjunlee/awscred/core"
	"github.com/hanjunlee/awscred/internal/daemon/server/mock"
	"github.com/sirupsen/logrus"
)

type (
	CredFileHandlerFunc   func(*gomock.Controller) CredFileHandler
	ConfigFileHandlerFunc func(*gomock.Controller) ConfigFileHandler
)

func TestInteractor_StartWatch(t *testing.T) {
	type fields struct {
		ch              chan fsnotify.Event
		watcher         FileWatcher
		origCredHandler CredFileHandler
		credHandler     CredFileHandler
		confHandler     ConfigFileHandler
		log             *logrus.Entry
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				ch:              tt.fields.ch,
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandler,
				credHandler:     tt.fields.credHandler,
				confHandler:     tt.fields.confHandler,
				log:             tt.fields.log,
			}
			i.StartWatch(tt.args.ctx)
		})
	}
}

func TestInteractor_runWorker(t *testing.T) {
	type fields struct {
		ch              chan fsnotify.Event
		watcher         FileWatcher
		origCredHandler CredFileHandler
		credHandler     CredFileHandler
		confHandler     ConfigFileHandler
		log             *logrus.Entry
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				ch:              tt.fields.ch,
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandler,
				credHandler:     tt.fields.credHandler,
				confHandler:     tt.fields.confHandler,
				log:             tt.fields.log,
			}
			i.runWorker(tt.args.ctx)
		})
	}
}

func TestInteractor_reflect(t *testing.T) {
	type fields struct {
		ch                  chan fsnotify.Event
		watcher             FileWatcher
		origCredHandlerFunc CredFileHandlerFunc
		credHandlerFunc     CredFileHandlerFunc
		confHandlerFunc     ConfigFileHandlerFunc
		log                 *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "reflect original creds",
			fields: fields{
				origCredHandlerFunc: func(c *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(c)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
							"profile_1": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
						}, nil)

					return m
				},
				credHandlerFunc: func(c *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(c)

					// expected
					m.EXPECT().
						Write(gomock.Eq(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
							"profile_1": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
						})).
						Return(nil)

					return m
				},
				confHandlerFunc: func(c *gomock.Controller) ConfigFileHandler {
					m := mock.NewMockConfigFileHandler(c)

					m.EXPECT().
						Read().
						Return(map[string]core.Config{
							"profile_0": {
								On: "false",
							},
						}, nil)

					return m
				},
			},
			wantErr: false,
		},
		{
			name: "reflect cache",
			fields: fields{
				origCredHandlerFunc: func(c *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(c)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
							"profile_1": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
						}, nil)

					return m
				},
				credHandlerFunc: func(c *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(c)

					m.EXPECT().
						Write(gomock.Eq(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "cachekey",
								SecretAccessKey: "cachesecret",
								SessionToken:    "cachetoken",
							},
							"profile_1": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
						})).
						Return(nil)

					return m
				},
				confHandlerFunc: func(c *gomock.Controller) ConfigFileHandler {
					m := mock.NewMockConfigFileHandler(c)

					m.EXPECT().
						Read().
						Return(map[string]core.Config{
							"profile_0": {
								On: "true",
								Cache: core.Cache{
									AccessKeyID:     "cachekey",
									SecretAccessKey: "cachesecret",
									SessionToken:    "cachetoken",
								},
							},
						}, nil)

					return m
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			i := &Interactor{
				ch:              tt.fields.ch,
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandlerFunc(c),
				credHandler:     tt.fields.credHandlerFunc(c),
				confHandler:     tt.fields.confHandlerFunc(c),
				log:             tt.fields.log,
			}
			if err := i.reflect(); (err != nil) != tt.wantErr {
				t.Errorf("Interactor.reflect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
