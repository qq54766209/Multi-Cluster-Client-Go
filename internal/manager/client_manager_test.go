package manager

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"sync"
	"testing"

	"multi-cluster-clientgo/internal/config"
)

func TestK8sClientManager(t *testing.T) {
	manager := NewK8sClientManager()

	// 模拟 TokenConfigSource
	source := &config.TokenConfigSource{
		APIServerURL: "https://10.19.225.29:6443",
		Token:        "eyJhbGciOiJSUzI1NiIsImtpZCI6ImNhb1VlWW1HWlBWRjdFdUN0Z0N0R0x2djVzLWJ0dVNhNmhEY0lzOEFNb00ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6Im11bHRlci1zYS10b2tlbi13MjU4cCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJtdWx0ZXItc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJkYmZkNDVmYS0yNTBiLTQyNmQtYTdlNy02YjU0ZTk5ZTVlMTciLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDptdWx0ZXItc2EifQ.iHXc3Y-fR3unhRZ7JQ1WREZ4zlPNwIOk8AR4aDp-16f_6E18lNLu4e4DEGUntgzmNIL6ZeiL0lFf-a6cN2VWS2Vz2Zmc3t9rSRymHEpZQr0daEP57C-SwQofpPcDcHhSQBx7NzIkCsRYP7JA0SDDgMAIzOcQsDrv4DV-5hx9wzx6BrPSDJEOvwgRvBhfxShnz22nSDtWk6qtWiClA25LRIIwNM_KF-kxExxSHQYJVX4GU_YN5nakxfkMocS0bU-uNXATgE_-Qz_RV0baDrSFs7S-1L9HFbLIKMIHqI-7H3h7qR4KGA2taMvAhEjCfEegiCB9yyzgtp5gRHt89lIn-g",
		Insecure:     true,
	}

	err := manager.AddClient("test-cluster", source)
	if err != nil {
		t.Fatalf("failed to add client: %v", err)
	}

	_, err = manager.GetClient("test-cluster")
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	client, err := manager.GetClient("test-cluster")
	if err != nil {
		log.Fatalf("Failed to get client for %s: %v", "cluster.Name", err)
	}

	pod, err := client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	for i := range pod.Items {
		fmt.Println(pod.Items[i].Name)
	}

	// 并发测试
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			manager.GetClient("test-cluster")
		}(i)
	}
	wg.Wait()
}
