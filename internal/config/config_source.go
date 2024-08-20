package config

import "k8s.io/client-go/rest"

// K8sConfigSource 是一个接口，定义了从不同来源获取 Kubernetes 连接配置的方法
type K8sConfigSource interface {
	GetConfig() (*rest.Config, error)
}
