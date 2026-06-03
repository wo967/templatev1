package protocol

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/zeromicro/go-zero/core/logx"
	"templatev1/pkg/mdns/model"
)

// MDNSServiceTypes mDNS查询的服务类型列表
var MDNSServiceTypes = []string{
	"_workstation._tcp.local.",
	"_http._tcp.local.",
	"_https._tcp.local.",
	"_smb._tcp.local.",
	"_afpovertcp._tcp.local.",
	"_qdiscover._tcp.local.",
	"_device-info._tcp.local.",
	"_ssh._tcp.local.",
	"_ftp._tcp.local.",
	"_printer._tcp.local.",
}

// QueryMDNS 执行mDNS查询
func QueryMDNS(ctx context.Context, ip net.IP, port int) ([]*dns.Msg, error) {
	mdnsAddr := "224.0.0.251:5353"
	if isIPv6(ip) {
		mdnsAddr = "[ff02::fb]:5353"
	}

	var responses []*dns.Msg

	for _, serviceType := range MDNSServiceTypes {
		select {
		case <-ctx.Done():
			return responses, ctx.Err()
		default:
		}

		msg := new(dns.Msg)
		msg.SetQuestion(serviceType, dns.TypePTR)
		msg.RecursionDesired = false

		resp, err := sendQuery(ctx, msg, mdnsAddr)
		if err != nil {
			logx.Debugf("查询服务 %s 失败: %v", serviceType, err)
			continue
		}

		if resp != nil && len(resp.Answer) > 0 {
			responses = append(responses, resp)
		}
	}

	return responses, nil
}

func sendQuery(ctx context.Context, msg *dns.Msg, addr string) (*dns.Msg, error) {
	client := &dns.Client{
		Timeout: 2 * time.Second,
		Net:     "udp",
	}

	resp, _, err := client.ExchangeContext(ctx, msg, addr)
	if err != nil {
		return nil, fmt.Errorf("DNS交换失败: %w", err)
	}

	return resp, nil
}

// ParseResponse 解析DNS响应为ServiceInfo
func ParseResponse(msg *dns.Msg) *model.ServiceInfo {
	if msg == nil || len(msg.Answer) == 0 {
		return nil
	}

	service := model.NewServiceInfo()

	for _, answer := range msg.Answer {
		switch rr := answer.(type) {
		case *dns.PTR:
			service.Name = extractServiceName(rr.Ptr)
			service.TTL = int(rr.Hdr.Ttl)

		case *dns.SRV:
			service.Hostname = normalizeHostname(rr.Target)
			service.Port = int(rr.Port)
			service.Protocol = "tcp"
			service.TTL = int(rr.Hdr.Ttl)

		case *dns.TXT:
			parseTXTRecords(service, rr.Txt)

		case *dns.A:
			service.IPv4 = rr.A.String()
			service.TTL = int(rr.Hdr.Ttl)

		case *dns.AAAA:
			service.IPv6 = rr.AAAA.String()
			service.TTL = int(rr.Hdr.Ttl)
		}
	}

	if service.Name == "" {
		return nil
	}

	return service
}

func extractServiceName(ptr string) string {
	parts := strings.Split(ptr, ".")
	if len(parts) >= 3 {
		servicePart := parts[0]
		return strings.TrimPrefix(servicePart, "_")
	}
	return ptr
}

func normalizeHostname(hostname string) string {
	return strings.TrimSuffix(hostname, ".")
}

func parseTXTRecords(service *model.ServiceInfo, texts []string) {
	for _, txt := range texts {
		if key, value, ok := parseKeyValue(txt); ok {
			service.AddDetail(key, value)
		}
	}
}

func parseKeyValue(s string) (key, value string, ok bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}

func isIPv6(ip net.IP) bool {
	return ip.To4() == nil && ip.To16() != nil
}
