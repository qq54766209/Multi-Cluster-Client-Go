package config

import (
	"testing"
)

func TestTokenConfigSource_GetConfigWithFile(t *testing.T) {
	// 测试从文件读取 Token 的情况
	source := &TokenConfigSource{
		APIServerURL: "https://10.19.225.29:6443",
		TokenPath:    "/path/to/nonexistent/token", // 使用一个不存在的路径
		Insecure:     true,
	}

	_, err := source.GetConfig()
	if err == nil {
		t.Fatalf("expected error due to invalid token path")
	}

	// 这里可以添加更多测试，例如使用一个实际存在的文件路径
}

func TestTokenConfigSource_GetConfigWithToken(t *testing.T) {
	// 测试直接传递 Token 的情况
	source := &TokenConfigSource{
		APIServerURL: "https://10.19.225.29:6443",
		Token:        "eyJhbGciOiJSUzI1NiIsImtpZCI6ImNhb1VlWW1HWlBWRjdFdUN0Z0N0R0x2djVzLWJ0dVNhNmhEY0lzOEFNb00ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6Im11bHRlci1zYS10b2tlbi13MjU4cCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJtdWx0ZXItc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJkYmZkNDVmYS0yNTBiLTQyNmQtYTdlNy02YjU0ZTk5ZTVlMTciLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDptdWx0ZXItc2EifQ.iHXc3Y-fR3unhRZ7JQ1WREZ4zlPNwIOk8AR4aDp-16f_6E18lNLu4e4DEGUntgzmNIL6ZeiL0lFf-a6cN2VWS2Vz2Zmc3t9rSRymHEpZQr0daEP57C-SwQofpPcDcHhSQBx7NzIkCsRYP7JA0SDDgMAIzOcQsDrv4DV-5hx9wzx6BrPSDJEOvwgRvBhfxShnz22nSDtWk6qtWiClA25LRIIwNM_KF-kxExxSHQYJVX4GU_YN5nakxfkMocS0bU-uNXATgE_-Qz_RV0baDrSFs7S-1L9HFbLIKMIHqI-7H3h7qR4KGA2taMvAhEjCfEegiCB9yyzgtp5gRHt89lIn-g",
		Insecure:     true,
	}

	_, err := source.GetConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewTokenConfigSource(t *testing.T) {
	// 测试通过正确的参数创建 TokenConfigSource
	apiServerURL := "https://10.19.225.29:6443"
	tokenPath := "token/cluster2-ca.crt"
	insecure := false

	source := NewTokenConfigSourceWithFile(apiServerURL, tokenPath, insecure, "")

	// 验证返回的 TokenConfigSource 是否正确配置
	if source.APIServerURL != apiServerURL {
		t.Errorf("expected API server URL to be %s, got %s", apiServerURL, source.APIServerURL)
	}

	if source.TokenPath != tokenPath {
		t.Errorf("expected token path to be %s, got %s", tokenPath, source.TokenPath)
	}

	if source.Insecure != insecure {
		t.Errorf("expected insecure to be %v, got %v", insecure, source.Insecure)
	}
}
