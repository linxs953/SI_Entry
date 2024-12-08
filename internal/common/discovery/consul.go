package discovery

import (
	"fmt"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	Host     string `json:",optional"`
	Port     int    `json:",optional"`
	Services map[string]string `json:",optional"` // service name -> consul service name mapping
}

type ServiceDiscovery struct {
	config     *ConsulConfig
	client     *consulapi.Client
	once       sync.Once
}

var (
	instance *ServiceDiscovery
	mu       sync.Mutex
)

// GetInstance returns the singleton instance of ServiceDiscovery
func GetInstance() *ServiceDiscovery {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = &ServiceDiscovery{}
		}
	}
	return instance
}

// Initialize initializes the service discovery with consul config
func (s *ServiceDiscovery) Initialize(config *ConsulConfig) error {
	var err error
	s.once.Do(func() {
		s.config = config
		if s.config.Host == "" {
			s.config.Host = "localhost"
		}
		if s.config.Port == 0 {
			s.config.Port = 8500
		}
		if s.config.Services == nil {
			s.config.Services = make(map[string]string)
		}

		consulConfig := consulapi.DefaultConfig()
		consulConfig.Address = fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

		s.client, err = consulapi.NewClient(consulConfig)
	})
	return err
}

// GetRPCAddress returns the RPC address for the given service
func (s *ServiceDiscovery) GetRPCAddress(serviceName string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("service discovery not initialized")
	}

	// Get the consul service name from mapping, or use the original name
	consulServiceName, exists := s.config.Services[serviceName]
	if !exists {
		consulServiceName = serviceName
	}

	services, _, err := s.client.Health().Service(consulServiceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get service address: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy service instances found for %s", serviceName)
	}

	service := services[0].Service
	address := service.Address
	if address == "" {
		address = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%d", address, service.Port), nil
}
