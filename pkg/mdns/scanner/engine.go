package scanner

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"templatev1/pkg/mdns/model"
	"templatev1/pkg/mdns/protocol"
)

// Engine 扫描引擎
type Engine struct {
	config    *model.ScanConfig
	semaphore chan struct{}
	mu        sync.RWMutex
	assetMap  map[string]*model.Asset
}

// NewEngine 创建新的扫描引擎
func NewEngine(config *model.ScanConfig) *Engine {
	return &Engine{
		config:    config,
		semaphore: make(chan struct{}, config.MaxGoroutines),
		assetMap:  make(map[string]*model.Asset),
	}
}

// Run 执行扫描
func (e *Engine) Run(ctx context.Context, targets []*model.ScanTarget, onFound func(*model.Asset)) error {
	var wg sync.WaitGroup
	var scannedCount int64

	logx.Infof("启动扫描引擎，并发数: %d", e.config.MaxGoroutines)

	for _, target := range targets {
		select {
		case <-ctx.Done():
			logx.Warnf("扫描被取消")
			wg.Wait()
			return ctx.Err()
		default:
		}

		e.semaphore <- struct{}{}

		wg.Add(1)
		go func(t *model.ScanTarget) {
			defer wg.Done()
			defer func() { <-e.semaphore }()

			asset, err := e.scanTarget(ctx, t)
			if err != nil {
				logx.Debugf("扫描目标 %s 失败: %v", t.String(), err)
				return
			}

			atomic.AddInt64(&scannedCount, 1)

			if asset != nil && asset.HasServices() {
				e.mu.Lock()
				e.assetMap[asset.IP] = asset
				e.mu.Unlock()

				onFound(asset)
			}

			if e.config.Verbose && atomic.LoadInt64(&scannedCount)%100 == 0 {
				logx.Infof("扫描进度: %d/%d", scannedCount, len(targets))
			}
		}(target)
	}

	wg.Wait()

	logx.Infof("扫描完成，共扫描 %d 个目标", scannedCount)
	return nil
}

func (e *Engine) scanTarget(ctx context.Context, target *model.ScanTarget) (*model.Asset, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(e.config.Timeout)*time.Second)
	defer cancel()

	responses, err := protocol.QueryMDNS(ctx, target.IP, target.Port)
	if err != nil {
		return nil, err
	}

	if len(responses) == 0 {
		return nil, nil
	}

	asset := model.NewAsset(target.IP.String())
	for _, resp := range responses {
		service := protocol.ParseResponse(resp)
		if service != nil {
			asset.AddService(*service)
		}
	}

	return asset, nil
}

// GetAssets 获取所有发现的资产
func (e *Engine) GetAssets() []*model.Asset {
	e.mu.RLock()
	defer e.mu.RUnlock()

	assets := make([]*model.Asset, 0, len(e.assetMap))
	for _, asset := range e.assetMap {
		assets = append(assets, asset)
	}

	return assets
}

// GenerateTargets 根据配置生成所有扫描目标
func GenerateTargets(config *model.ScanConfig) ([]*model.ScanTarget, error) {
	ips, err := parseCIDR(config.CIDR)
	if err != nil {
		return nil, fmt.Errorf("解析CIDR失败: %w", err)
	}

	ports, err := config.ParsePortRanges()
	if err != nil {
		return nil, fmt.Errorf("解析端口范围失败: %w", err)
	}

	targets := make([]*model.ScanTarget, 0, len(ips)*len(ports))
	for _, ip := range ips {
		for _, port := range ports {
			targets = append(targets, &model.ScanTarget{
				IP:   ip,
				Port: port,
			})
		}
	}

	return targets, nil
}

func parseCIDR(cidr string) ([]net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		ips = append(ips, net.ParseIP(ip.String()))
	}

	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}

	return ips, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
