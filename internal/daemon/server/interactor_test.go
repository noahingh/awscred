package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/mock/gomock"
	"github.com/hanjunlee/awscred/core"
	"github.com/hanjunlee/awscred/internal/daemon/server/mock"
	"github.com/sirupsen/logrus"
)

type (
	SessionTokenGeneratorFunc func(*gomock.Controller) SessionTokenGenerator
	CredFileHandlerFunc       func(*gomock.Controller) CredFileHandler
	ConfigFileHandlerFunc     func(*gomock.Controller) ConfigFileHandler
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
								On: false,
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
								On: true,
								Cache: core.SessionToken{
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

func TestInteractor_On(t *testing.T) {
	type fields struct {
		ch                  chan fsnotify.Event
		stGenerator         SessionTokenGenerator
		watcher             FileWatcher
		origCredHandlerFunc CredFileHandlerFunc
		credHandler         CredFileHandler
		confHandlerFunc     ConfigFileHandlerFunc
		log                 *logrus.Entry
	}
	type args struct {
		profile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set the profile enabled",
			fields: fields{
				origCredHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{
							"profile_0": {},
							"profile_1": {},
						}, nil)

					return m
				},
				confHandlerFunc: func(ctrl *gomock.Controller) ConfigFileHandler {
					m := mock.NewMockConfigFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Config{
							"profile_0": {
								On:             false,
								SerialNumber:   "serial",
								DurationSecond: 3600,
							},
						}, nil).
						AnyTimes()

					m.EXPECT().
						Write(gomock.Eq(map[string]core.Config{
							"profile_0": {
								On:             true, // set true
								SerialNumber:   "serial",
								DurationSecond: 3600,
							},
						})).
						Return(nil)

					return m
				},
			},
			args: args{
				profile: "profile_0",
			},
			wantErr: false,
		},
		{
			name: "create a new config if it doesn't exist",
			fields: fields{
				origCredHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{
							"profile_0": {},
							"profile_1": {},
						}, nil)

					return m
				},
				confHandlerFunc: func(ctrl *gomock.Controller) ConfigFileHandler {
					m := mock.NewMockConfigFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Config{}, nil).
						AnyTimes()

					// expected.
					m.EXPECT().
						Write(gomock.Eq(map[string]core.Config{
							"profile_0": {
								On: true, // set true
							},
						})).
						Return(nil)

					return m
				},
			},
			args: args{
				profile: "profile_0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			i := &Interactor{
				ch:              tt.fields.ch,
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandlerFunc(ctrl),
				credHandler:     tt.fields.credHandler,
				confHandler:     tt.fields.confHandlerFunc(ctrl),
				log:             tt.fields.log,
			}
			if err := i.On(tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("Interactor.On() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInteractor_Off(t *testing.T) {
	type fields struct {
		ch                  chan fsnotify.Event
		stGenerator         SessionTokenGenerator
		watcher             FileWatcher
		origCredHandlerFunc CredFileHandlerFunc
		credHandler         CredFileHandler
		confHandlerFunc     ConfigFileHandlerFunc
		log                 *logrus.Entry
	}
	type args struct {
		profile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set the profile disabled",
			fields: fields{
				origCredHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{
							"profile_0": {},
							"profile_1": {},
						}, nil)

					return m
				},
				confHandlerFunc: func(ctrl *gomock.Controller) ConfigFileHandler {
					m := mock.NewMockConfigFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Config{
							"profile_0": {
								On:             true,
								SerialNumber:   "serial",
								DurationSecond: 3600,
							},
						}, nil).
						AnyTimes()

					m.EXPECT().
						Write(gomock.Eq(map[string]core.Config{
							"profile_0": {
								On:             false, // set false
								SerialNumber:   "serial",
								DurationSecond: 3600,
							},
						})).
						Return(nil)

					return m
				},
			},
			args: args{
				profile: "profile_0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			i := &Interactor{
				ch:              tt.fields.ch,
				stGenerator:     tt.fields.stGenerator,
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandlerFunc(ctrl),
				credHandler:     tt.fields.credHandler,
				confHandler:     tt.fields.confHandlerFunc(ctrl),
				log:             tt.fields.log,
			}
			if err := i.Off(tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Off() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInteractor_Gen(t *testing.T) {
	type fields struct {
		ch                  chan fsnotify.Event
		stGeneratorFunc     SessionTokenGeneratorFunc
		watcher             FileWatcher
		origCredHandlerFunc CredFileHandlerFunc
		credHandlerFunc     CredFileHandlerFunc
		confHandlerFunc     ConfigFileHandlerFunc
		log                 *logrus.Entry
	}
	type args struct {
		profile string
		token   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "generate a session token",
			fields: fields{
				stGeneratorFunc: func(ctrl *gomock.Controller) SessionTokenGenerator {
					m := mock.NewMockSessionTokenGenerator(ctrl)

					m.EXPECT().
						Generate(
							gomock.Any(),
							gomock.Eq(core.Config{
								On:           true,
								SerialNumber: "serial",
							}),
							gomock.Any(),
						).
						Return(core.SessionToken{
							AccessKeyID:     "sessionkey",
							SecretAccessKey: "sessionsecret",
							SessionToken:    "sessiontoken",
						}, nil)

					return m
				},
				origCredHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
						}, nil).
						AnyTimes()

					return m
				},
				credHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
						}, nil).
						AnyTimes()

						// after generate the session token.
					m.EXPECT().
						Write(gomock.Eq(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "sessionkey",
								SecretAccessKey: "sessionsecret",
								SessionToken:    "sessiontoken",
							},
						})).
						Return(nil)

					return m
				},
				confHandlerFunc: func(ctrl *gomock.Controller) ConfigFileHandler {
					m := mock.NewMockConfigFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Config{
							"profile_0": {
								On:           true,
								SerialNumber: "serial",
							},
						}, nil).
						AnyTimes()

						// after generate the session token
					m.EXPECT().
						Write(map[string]core.Config{
							"profile_0": {
								On:           true,
								SerialNumber: "serial",
								Cache: core.SessionToken{
									AccessKeyID:     "sessionkey",
									SecretAccessKey: "sessionsecret",
									SessionToken:    "sessiontoken",
								},
							},
						})

					return m
				},
			},
			args: args{
				profile: "profile_0",
				token:   "123456",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			i := &Interactor{
				ch:              tt.fields.ch,
				stGenerator:     tt.fields.stGeneratorFunc(ctrl),
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandlerFunc(ctrl),
				credHandler:     tt.fields.credHandlerFunc(ctrl),
				confHandler:     tt.fields.confHandlerFunc(ctrl),
				log:             tt.fields.log,
			}
			if err := i.Gen(tt.args.profile, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Gen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInteractor_GetOriginalCred(t *testing.T) {
	type fields struct {
		ch                  chan fsnotify.Event
		stGenerator         SessionTokenGenerator
		watcher             FileWatcher
		origCredHandlerFunc CredFileHandlerFunc
		credHandler         CredFileHandler
		confHandler         ConfigFileHandler
		log                 *logrus.Entry
	}
	type args struct {
		profile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    core.Cred
		want1   bool
		wantErr bool
	}{
		{
			name: "return false if couldn't find the profile",
			fields: fields{
				origCredHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

					m.EXPECT().
						Read().
						Return(map[string]core.Cred{}, nil)

					return m
				},
			},
			args: args{
				profile: "profile",
			},
			want:    core.Cred{},
			want1:   false,
			wantErr: false,
		},
		{
			name: "return the cred of profile",
			fields: fields{
				origCredHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

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
			},
			args: args{
				profile: "profile_0",
			},
			want: core.Cred{
				AccessKeyID:     "key",
				SecretAccessKey: "secret",
			},
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			i := &Interactor{
				ch:              tt.fields.ch,
				stGenerator:     tt.fields.stGenerator,
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandlerFunc(ctrl),
				credHandler:     tt.fields.credHandler,
				confHandler:     tt.fields.confHandler,
				log:             tt.fields.log,
			}
			got, got1, err := i.GetOriginalCred(tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interactor.GetOriginalCred() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interactor.GetOriginalCred() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Interactor.GetOriginalCred() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInteractor_SetCred(t *testing.T) {
	type fields struct {
		ch              chan fsnotify.Event
		stGenerator     SessionTokenGenerator
		watcher         FileWatcher
		origCredHandler CredFileHandler
		credHandlerFunc CredFileHandlerFunc
		confHandler     ConfigFileHandler
		log             *logrus.Entry
	}
	type args struct {
		profile string
		cred    core.Cred
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set a new credential",
			fields: fields{
				credHandlerFunc: func(ctrl *gomock.Controller) CredFileHandler {
					m := mock.NewMockCredFileHandler(ctrl)

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

					m.EXPECT().
						Write(gomock.Eq(map[string]core.Cred{
							"profile_0": {
								AccessKeyID:     "key",
								SecretAccessKey: "secret",
							},
							"profile_1": {
								AccessKeyID:     "tokenkey",
								SecretAccessKey: "tokensecret",
								SessionToken:    "token",
							},
						})).
						Return(nil)

					return m
				},
			},
			args: args{
				profile: "profile_1",
				cred: core.Cred{
					AccessKeyID:     "tokenkey",
					SecretAccessKey: "tokensecret",
					SessionToken:    "token",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			i := &Interactor{
				ch:              tt.fields.ch,
				stGenerator:     tt.fields.stGenerator,
				watcher:         tt.fields.watcher,
				origCredHandler: tt.fields.origCredHandler,
				credHandler:     tt.fields.credHandlerFunc(ctrl),
				confHandler:     tt.fields.confHandler,
				log:             tt.fields.log,
			}
			if err := i.SetCred(tt.args.profile, tt.args.cred); (err != nil) != tt.wantErr {
				t.Errorf("Interactor.SetCred() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
