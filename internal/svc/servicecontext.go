package svc

import (
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"

	"entry/internal/config"
)

type ClientManager struct {
	clients sync.Map
	conf    *config.Config
}

func NewClientManager(c config.Config) *ClientManager {
	return &ClientManager{
		conf: &c,
	}
}

func (m *ClientManager) GetClient(serviceName string) (*grpc.ClientConn, bool) {
	if client, ok := m.clients.Load(serviceName); ok {
		return client.(*grpc.ClientConn), true
	}

	target := m.conf.Services[serviceName]
	logx.Infof("Service %s target: %s", serviceName, target)
	if target == "" {
		return nil, false
	}

	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: target,
		App:    serviceName,
	})

	conn := client.Conn()
	m.clients.Store(serviceName, conn)
	return conn, true
}

type ServiceContext struct {
	Config        config.Config
	ClientManager *ClientManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		ClientManager: NewClientManager(c),
	}
}
