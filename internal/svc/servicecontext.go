package svc

import (
	"net"
	"sync"

	"entry/internal/config"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type ClientManager struct {
	clients sync.Map
	conf    *config.Config
	mu      sync.Mutex
}

func NewClientManager(c config.Config) *ClientManager {
	return &ClientManager{
		conf: &c,
	}
}

func (m *ClientManager) GetClient(serviceName string) (*grpc.ClientConn, error) {
	// 检查现有连接
	if conn, ok := m.clients.Load(serviceName); ok {
		return conn.(*grpc.ClientConn), nil
	}

	// 获取目标地址
	target, ok := m.conf.Services[serviceName]
	if !ok || target == "" {
		return nil, fmt.Errorf("service %s not configured", serviceName)
	}

	// 加锁以确保并发安全
	m.mu.Lock()
	defer m.mu.Unlock()

	// 再次检查，避免在等待锁时已被其他goroutine创建
	if conn, ok := m.clients.Load(serviceName); ok {
		return conn.(*grpc.ClientConn), nil
	}

	// 检查服务是否可用，如果可用再创建连接
	if !m.serviceAvailable(serviceName) {
		return nil, fmt.Errorf("service %s is unavailable", serviceName)
	}

	// 创建新的客户端连接
	client, err := zrpc.NewClient(zrpc.RpcClientConf{
		Target:  target,
		App:     serviceName,
		Timeout: 500,
	})

	if err != nil {
		logx.Errorf("Failed to connect to service %s: %v", serviceName, err)
		return nil, fmt.Errorf("service unavailable: %s", serviceName)
	}

	conn := client.Conn()
	m.clients.Store(serviceName, conn)
	return conn, nil
}

func (m *ClientManager) serviceAvailable(serviceName string) bool {
	target, ok := m.conf.Services[serviceName]
	if !ok || target == "" {
		return false
	}
	_, err := net.Dial("tcp", target)
	return err == nil
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
