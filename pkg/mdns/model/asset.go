package model

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

// ScanTarget 表示单个扫描目标（IP:Port组合）
type ScanTarget struct {
	IP   net.IP // IP地址（支持IPv4和IPv6）
	Port int    // 端口号（1-65535）
}

// Validate 验证扫描目标的合法性
func (st *ScanTarget) Validate() error {
	if st.IP == nil {
		return fmt.Errorf("IP地址不能为空")
	}
	if st.Port < 1 || st.Port > 65535 {
		return fmt.Errorf("端口号必须在1-65535范围内，当前值: %d", st.Port)
	}
	return nil
}

// String 返回目标的可读字符串表示
func (st *ScanTarget) String() string {
	return fmt.Sprintf("%s:%d", st.IP.String(), st.Port)
}

// ServiceInfo 表示单个服务的详细信息
type ServiceInfo struct {
	Port     int               `json:"port" yaml:"port"`
	Protocol string            `json:"protocol" yaml:"protocol"`
	Name     string            `json:"name" yaml:"name"`
	IPv4     string            `json:"ipv4,omitempty" yaml:"ipv4,omitempty"`
	IPv6     string            `json:"ipv6,omitempty" yaml:"ipv6,omitempty"`
	Hostname string            `json:"hostname" yaml:"hostname"`
	TTL      int               `json:"ttl" yaml:"ttl"`
	Details  map[string]string `json:"details,omitempty" yaml:"details,omitempty"`
}

// NewServiceInfo 创建新的ServiceInfo实例
func NewServiceInfo() *ServiceInfo {
	return &ServiceInfo{
		Details: make(map[string]string),
	}
}

// AddDetail 添加服务详细信息（避免nil panic）
func (si *ServiceInfo) AddDetail(key, value string) {
	if si.Details == nil {
		si.Details = make(map[string]string)
	}
	si.Details[key] = value
}

// GetServiceLabel 获取服务的显示标签
func (si *ServiceInfo) GetServiceLabel() string {
	return fmt.Sprintf("%d/%s %s:", si.Port, si.Protocol, si.Name)
}

// DeviceInfo 表示设备的汇总信息
type DeviceInfo struct {
	Name     string `json:"name" yaml:"name"`
	IPv4     string `json:"ipv4" yaml:"ipv4"`
	IPv6     string `json:"ipv6" yaml:"ipv6"`
	Hostname string `json:"hostname" yaml:"hostname"`
	TTL      int    `json:"ttl" yaml:"ttl"`
	Model    string `json:"model" yaml:"model"`
}

// NewDeviceInfo 创建新的DeviceInfo实例
func NewDeviceInfo() *DeviceInfo {
	return &DeviceInfo{}
}

// Asset 表示一个完整的资产
type Asset struct {
	IP         string              `json:"ip" yaml:"ip"`
	Services   []ServiceInfo       `json:"services,omitempty" yaml:"services,omitempty"`
	DeviceInfo *DeviceInfo         `json:"device-info,omitempty" yaml:"device-info,omitempty"`
	Answers    map[string][]string `json:"answers,omitempty" yaml:"answers,omitempty"`
}

// NewAsset 创建新的Asset实例
func NewAsset(ip string) *Asset {
	return &Asset{
		IP:       ip,
		Services: make([]ServiceInfo, 0),
		Answers:  make(map[string][]string),
	}
}

// AddService 添加服务信息（自动去重和合并）
func (a *Asset) AddService(service ServiceInfo) {
	for i, existing := range a.Services {
		if existing.Port == service.Port && existing.Protocol == service.Protocol {
			if service.Details != nil {
				for k, v := range service.Details {
					a.Services[i].AddDetail(k, v)
				}
			}
			if a.Services[i].Hostname == "" && service.Hostname != "" {
				a.Services[i].Hostname = service.Hostname
			}
			if a.Services[i].IPv4 == "" && service.IPv4 != "" {
				a.Services[i].IPv4 = service.IPv4
			}
			if a.Services[i].IPv6 == "" && service.IPv6 != "" {
				a.Services[i].IPv6 = service.IPv6
			}
			return
		}
	}
	a.Services = append(a.Services, service)
}

// AddPTRRecord 添加PTR记录到answers分类
func (a *Asset) AddPTRRecord(recordType, value string) {
	if a.Answers == nil {
		a.Answers = make(map[string][]string)
	}

	values := a.Answers[recordType]
	for _, v := range values {
		if v == value {
			return
		}
	}
	a.Answers[recordType] = append(values, value)
}

// HasServices 判断资产是否有服务信息
func (a *Asset) HasServices() bool {
	return len(a.Services) > 0
}

// GetPrimaryHostname 获取主要主机名
func (a *Asset) GetPrimaryHostname() string {
	if a.DeviceInfo != nil && a.DeviceInfo.Hostname != "" {
		return a.DeviceInfo.Hostname
	}
	for _, svc := range a.Services {
		if svc.Hostname != "" {
			return svc.Hostname
		}
	}
	return ""
}

// SortServices 对服务列表按端口排序
func (a *Asset) SortServices() {
	sort.Slice(a.Services, func(i, j int) bool {
		return a.Services[i].Port < a.Services[j].Port
	})
}

// ScanConfig 扫描配置参数
type ScanConfig struct {
	CIDR          string
	PortRanges    []string
	Timeout       int
	MaxGoroutines int
	OutputFormat  string
	OutputFile    string
	Verbose       bool
}

// Validate 验证扫描配置的合法性
func (sc *ScanConfig) Validate() error {
	if sc.CIDR == "" {
		return fmt.Errorf("CIDR网段不能为空")
	}

	_, _, err := net.ParseCIDR(sc.CIDR)
	if err != nil {
		return fmt.Errorf("无效的CIDR格式 '%s': %w", sc.CIDR, err)
	}

	if len(sc.PortRanges) == 0 {
		return fmt.Errorf("端口范围不能为空")
	}

	for _, portRange := range sc.PortRanges {
		if err := validatePortRange(portRange); err != nil {
			return fmt.Errorf("无效的端口范围 '%s': %w", portRange, err)
		}
	}

	if sc.Timeout <= 0 {
		sc.Timeout = 2
	}
	if sc.MaxGoroutines <= 0 {
		sc.MaxGoroutines = 100
	}
	if sc.OutputFormat == "" {
		sc.OutputFormat = "json"
	}

	if sc.OutputFormat != "yaml" && sc.OutputFormat != "json" {
		return fmt.Errorf("不支持的输出格式 '%s'，仅支持yaml或json", sc.OutputFormat)
	}

	return nil
}

func validatePortRange(portRange string) error {
	parts := strings.Split(portRange, "-")

	switch len(parts) {
	case 1:
		port := parts[0]
		if !isValidPortString(port) {
			return fmt.Errorf("无效端口号 '%s'", port)
		}
	case 2:
		start, end := parts[0], parts[1]
		if !isValidPortString(start) || !isValidPortString(end) {
			return fmt.Errorf("无效端口范围 '%s'", portRange)
		}

		startPort := parsePortString(start)
		endPort := parsePortString(end)

		if startPort > endPort {
			return fmt.Errorf("起始端口(%d)不能大于结束端口(%d)", startPort, endPort)
		}
	default:
		return fmt.Errorf("端口范围格式错误，应为 'port' 或 'start-end'")
	}

	return nil
}

func isValidPortString(s string) bool {
	port := parsePortString(s)
	return port >= 1 && port <= 65535
}

func parsePortString(s string) int {
	var port int
	if _, err := fmt.Sscanf(s, "%d", &port); err != nil {
		return 0
	}
	return port
}

// ParsePortRanges 解析所有端口范围为整数切片
func (sc *ScanConfig) ParsePortRanges() ([]int, error) {
	portSet := make(map[int]bool)

	for _, portRange := range sc.PortRanges {
		ports, err := parseSinglePortRange(portRange)
		if err != nil {
			return nil, err
		}
		for _, port := range ports {
			portSet[port] = true
		}
	}

	ports := make([]int, 0, len(portSet))
	for port := range portSet {
		ports = append(ports, port)
	}

	for i := 0; i < len(ports); i++ {
		for j := i + 1; j < len(ports); j++ {
			if ports[i] > ports[j] {
				ports[i], ports[j] = ports[j], ports[i]
			}
		}
	}

	return ports, nil
}

func parseSinglePortRange(portRange string) ([]int, error) {
	parts := strings.Split(portRange, "-")

	if len(parts) == 1 {
		port := parsePortString(parts[0])
		return []int{port}, nil
	}

	start := parsePortString(parts[0])
	end := parsePortString(parts[1])

	ports := make([]int, 0, end-start+1)
	for port := start; port <= end; port++ {
		ports = append(ports, port)
	}

	return ports, nil
}
