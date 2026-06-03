// ==================== mDNS扫描相关类型 ====================

// ScanReq mDNS扫描请求
type ScanReq struct {
CIDR        string   `json:"cidr"`
PortRanges  []string `json:"portRanges"`
Timeout     int      `json:"timeout"`
Concurrency int      `json:"concurrency"`
}

// ScanResp mDNS扫描响应
type ScanResp struct {
Assets []AssetInfo `json:"assets"`
Stats  ScanStats   `json:"stats"`
}

// AssetInfo 资产信息
type AssetInfo struct {
IP       string            `json:"ip"`
Services []ServiceInfo     `json:"services"`
DeviceInfo *DeviceInfo     `json:"deviceInfo,omitempty"`
Answers  map[string][]string `json:"answers,omitempty"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
Port     int               `json:"port"`
Protocol string            `json:"protocol"`
Name     string            `json:"name"`
IPv4     string            `json:"ipv4"`
IPv6     string            `json:"ipv6"`
Hostname string            `json:"hostname"`
TTL      int               `json:"ttl"`
Details  map[string]string `json:"details"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
Name     string `json:"name"`
IPv4     string `json:"ipv4"`
IPv6     string `json:"ipv6"`
Hostname string `json:"hostname"`
TTL      int    `json:"ttl"`
Model    string `json:"model"`
}

// ScanStats 扫描统计
type ScanStats struct {
TotalTargets int    `json:"totalTargets"`
Scanned      int    `json:"scanned"`
Found        int    `json:"found"`
Duration     string `json:"duration"`
Error        string `json:"error,omitempty"`
}
