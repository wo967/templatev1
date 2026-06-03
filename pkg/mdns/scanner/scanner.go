package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"templatev1/pkg/mdns/model"
)

// ScanResult 扫描结果
type ScanResult struct {
	Assets []*model.Asset
	Stats  ScanStats
}

// ScanStats 扫描统计信息
type ScanStats struct {
	TotalTargets int    `json:"total_targets" yaml:"total_targets"`
	Scanned      int    `json:"scanned" yaml:"scanned"`
	Found        int    `json:"found" yaml:"found"`
	Duration     string `json:"duration" yaml:"duration"`
	Error        error  `json:"error,omitempty" yaml:"error,omitempty"`
}

// Scan 执行mDNS扫描（主入口函数）
func Scan(ctx context.Context, config *model.ScanConfig) (*ScanResult, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	logx.Infof("开始执行mDNS扫描...")

	targets, err := GenerateTargets(config)
	if err != nil {
		return nil, fmt.Errorf("生成扫描目标失败: %w", err)
	}

	logx.Infof("生成 %d 个扫描目标", len(targets))

	engine := NewEngine(config)

	result := &ScanResult{
		Assets: make([]*model.Asset, 0),
		Stats: ScanStats{
			TotalTargets: len(targets),
		},
	}

	startTime := time.Now()

	err = engine.Run(ctx, targets, func(asset *model.Asset) {
		if asset != nil && asset.HasServices() {
			result.Assets = append(result.Assets, asset)
			result.Stats.Found++
		}
		result.Stats.Scanned++
	})

	duration := time.Since(startTime)
	result.Stats.Duration = duration.String()

	if err != nil {
		result.Stats.Error = err
		return result, fmt.Errorf("扫描执行失败: %w", err)
	}

	logx.Infof("扫描完成，发现 %d 个资产", result.Stats.Found)
	return result, nil
}

// ScanSync 同步扫描（简化版）
func ScanSync(config *model.ScanConfig) *ScanResult {
	ctx := context.Background()
	result, err := Scan(ctx, config)
	if err != nil {
		return &ScanResult{
			Stats: ScanStats{
				Error: err,
			},
		}
	}
	return result
}
