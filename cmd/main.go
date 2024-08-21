package main

import (
	"context"
	"fmt"
	"log"
	"multi-cluster-clientgo/internal/config"
	"multi-cluster-clientgo/internal/manager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// 假设我们有多个集群的配置信息
	clusters := []struct {
		Name         string
		APIServerURL string
		Token        string
		TokenPath    string
		Insecure     bool
		CACertPath   string // CA 证书路径
	}{
		{
			Name:         "cluster1",
			APIServerURL: "https://10.19.225.29:6443",
			TokenPath:    "token/cluster1-token.txt",
			Insecure:     true,
		},
		{
			Name:         "cluster2",
			APIServerURL: "https://10.19.225.30:6443",
			Token:        "eyJhbGciOiJSUzI1NiIsImtpZCI6InQ4dnVIUFhaVkxEWFhzTmQ4RWRMQkp6MW5uc3F5dWlFSEd2a19vRThXTU0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6Im11bHRlci1zYS10b2tlbi1ua3NkdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJtdWx0ZXItc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI0ZGEwMTkxZS1kOTdmLTRhNDgtYmNlNi1hNDQxZmU3ZmE5MzciLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDptdWx0ZXItc2EifQ.fYmFkdVgsrLXSZozt3McQRM3qNTFZWiwdqfQQvILAtd8pRbpOH8wI6BuVRLaia2zFw1rI5q4w87oZypgCIgpMqWtrCftXmqpM0q0U6X3OzxqQQZi-cvdlN6PodtUF0988RrsSMe1DkQ8dPpzNkzUFAFVig4HZH6071MWRxh-YR8yq1q1czcxWC0AKfDNzjA2mZ--vI3QNN_5n8cRnA6JMi0E6TKGw1gp2Oy-EFdfMkZZR1God6yGlRWrCAyDtyJtZGm-pWXrK_RLpZQX2jte42E43PrFVvjwP_82KKvbGZvi1J_yFpqai6P84PLKA8tECVwb7m860RRF7Mfz0l1aQg", // 直接传递 Token
			Insecure:     true,
			CACertPath:   "token/cluster2-ca.crt",
		},
	}

	// 创建一个 K8sClientManager 实例
	clientManager := manager.NewK8sClientManager()

	// 为每个集群创建一个 TokenConfigSource 并添加到 K8sClientManager 中
	for _, cluster := range clusters {
		var tokenSource *config.TokenConfigSource
		if cluster.Token != "" {
			tokenSource = config.NewTokenConfigSourceWithToken(cluster.APIServerURL, cluster.Token, cluster.Insecure, cluster.CACertPath)
		} else if cluster.TokenPath != "" {
			tokenSource = config.NewTokenConfigSourceWithFile(cluster.APIServerURL, cluster.TokenPath, cluster.Insecure, cluster.CACertPath)
		} else {
			log.Fatalf("No token or tokenPath provided for cluster %s", cluster.Name)
		}

		if err := clientManager.AddClient(cluster.Name, tokenSource); err != nil {
			log.Fatalf("Failed to add %s: %v", cluster.Name, err)
		}
	}

	// 示例操作：获取并列出每个集群的 namespaces
	for _, cluster := range clusters {
		client, err := clientManager.GetClient(cluster.Name)
		if err != nil {
			log.Fatalf("Failed to get client for %s: %v", cluster.Name, err)
		}

		// 获取所有命名空间中的所有 Pod 名称
		pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Failed to list pods: %v", err)
		}

		for _, pod := range pods.Items {
			fmt.Printf("Pod Name: %s, Namespace: %s\n", pod.Name, pod.Namespace)
		}
	}

	// 优雅关闭所有客户端
	clientManager.Shutdown()
}
