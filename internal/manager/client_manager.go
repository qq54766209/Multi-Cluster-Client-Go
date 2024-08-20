package manager

import (
	"errors"
	"log"
	"sync"

	"k8s.io/client-go/kubernetes"
	"multi-cluster-clientgo/internal/config"
)

// K8sClientManager 管理多个 Kubernetes 集群的客户端
type K8sClientManager struct {
	mu      sync.RWMutex
	clients map[string]*kubernetes.Clientset
}

// NewK8sClientManager 创建一个新的 K8sClientManager
func NewK8sClientManager() *K8sClientManager {
	return &K8sClientManager{
		clients: make(map[string]*kubernetes.Clientset),
	}
}

// AddClient 使用指定的 K8sConfigSource 添加一个新的集群客户端
func (m *K8sClientManager) AddClient(name string, source config.K8sConfigSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[name]; exists {
		log.Printf("Client for cluster %s already exists, skipping creation", name)
		return errors.New("client for cluster " + name + " already exists")
	}

	config, err := source.GetConfig()
	if err != nil {
		log.Printf("Failed to get config for cluster %s: %v", name, err)
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create clientset for cluster %s: %v", name, err)
		return err
	}

	m.clients[name] = clientset
	log.Printf("Successfully added client for cluster %s", name)
	return nil
}

// UpdateClient 更新现有集群的客户端配置
func (m *K8sClientManager) UpdateClient(name string, source config.K8sConfigSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	config, err := source.GetConfig()
	if err != nil {
		log.Printf("Failed to update config for cluster %s: %v", name, err)
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to update clientset for cluster %s: %v", name, err)
		return err
	}

	m.clients[name] = clientset
	log.Printf("Successfully updated client for cluster %s", name)
	return nil
}

// GetClient 根据集群名返回对应的 Kubernetes Clientset
func (m *K8sClientManager) GetClient(name string) (*kubernetes.Clientset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, exists := m.clients[name]
	if !exists {
		log.Printf("Client not found for cluster %s", name)
		return nil, errors.New("client not found for cluster: " + name)
	}
	return client, nil
}

// Shutdown 优雅地关闭所有客户端连接
func (m *K8sClientManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name := range m.clients {
		delete(m.clients, name)
		log.Printf("Client for cluster %s has been shutdown", name)
	}
}
