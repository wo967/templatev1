package scanner

import (
	"context"
	"fmt"

	"api/internal/svc"
	"api/internal/types"
	"templatev1/pkg/mdns/model"
	"templatev1/pkg/mdns/scanner"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScanLogic {
	return &ScanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Scan 执行mDNS扫描
func (l *ScanLogic) Scan(req *types.ScanReq) (resp *types.ScanResp, err error) {
	config := &model.ScanConfig{
		CIDR:          req.CIDR,
		PortRanges:    req.PortRanges,
		Timeout:       req.Timeout,
		MaxGoroutines: req.Concurrency,
		OutputFormat:  "json",
		Verbose:       false,
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	l.Infof("开始mDNS扫描，CIDR: %s, 端口: %v", config.CIDR, config.PortRanges)

	result := scanner.ScanSync(config)

	if result.Stats.Error != nil {
		l.Errorf("扫描失败: %v", result.Stats.Error)
		return nil, fmt.Errorf("扫描失败: %w", result.Stats.Error)
	}

	resp = &types.ScanResp{
		Assets: convertAssets(result.Assets),
		Stats: types.ScanStats{
			TotalTargets: result.Stats.TotalTargets,
			Scanned:      result.Stats.Scanned,
			Found:        result.Stats.Found,
			Duration:     result.Stats.Duration,
		},
	}

	l.Infof("扫描完成，发现 %d 个资产", resp.Stats.Found)
	return resp, nil
}

func convertAssets(assets []*model.Asset) []types.AssetInfo {
	result := make([]types.AssetInfo, 0, len(assets))

	for _, asset := range assets {
		info := types.AssetInfo{
			IP:       asset.IP,
			Services: make([]types.ServiceInfo, 0, len(asset.Services)),
			Answers:  asset.Answers,
		}

		for _, svc := range asset.Services {
			info.Services = append(info.Services, types.ServiceInfo{
				Port:     svc.Port,
				Protocol: svc.Protocol,
				Name:     svc.Name,
				IPv4:     svc.IPv4,
				IPv6:     svc.IPv6,
				Hostname: svc.Hostname,
				TTL:      svc.TTL,
				Details:  svc.Details,
			})
		}

		if asset.DeviceInfo != nil {
			info.DeviceInfo = &types.DeviceInfo{
				Name:     asset.DeviceInfo.Name,
				IPv4:     asset.DeviceInfo.IPv4,
				IPv6:     asset.DeviceInfo.IPv6,
				Hostname: asset.DeviceInfo.Hostname,
				TTL:      asset.DeviceInfo.TTL,
				Model:    asset.DeviceInfo.Model,
			}
		}

		result = append(result, info)
	}

	return result
}
