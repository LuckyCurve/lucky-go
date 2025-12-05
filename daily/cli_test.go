package daily

import (
	"strings"
	"testing"

	"lucky-go/forex"
)

func TestFormatDailyMessage(t *testing.T) {
	report := &DailyReport{
		Treasury: 4.5,
		AAA:      5.0,
		BAA:      5.5,
		CAPE:     30.0,
		FairPE:   22.22,
		Premium:  35.0,
		ForexResult: &forex.ExchangeResult{
			From:       "USD",
			To:         "CNY",
			Rate:       7.2345,
			Amount:     1,
			Converted:  7.2345,
			UpdateDate: "2025-01-01",
		},
	}

	message := formatDailyMessage(report)

	// 验证消息不为空
	if message == "" {
		t.Error("消息不应为空")
	}

	// 验证包含标题
	if !strings.Contains(message, "每日金融综合报告") {
		t.Error("消息应包含标题 '每日金融综合报告'")
	}

	// 验证包含 PE 数据
	if !strings.Contains(message, "PE 估值") {
		t.Error("消息应包含 'PE 估值'")
	}

	// 验证包含 CAPE 数据
	if !strings.Contains(message, "CAPE 估值") {
		t.Error("消息应包含 'CAPE 估值'")
	}

	// 验证包含汇率数据
	if !strings.Contains(message, "汇率信息") {
		t.Error("消息应包含 '汇率信息'")
	}

	// 验证包含货币对
	if !strings.Contains(message, "USD → CNY") {
		t.Error("消息应包含货币对 'USD → CNY'")
	}
}

func TestFormatDailyMessage_HighPremium(t *testing.T) {
	report := &DailyReport{
		Treasury: 4.5,
		AAA:      5.0,
		BAA:      5.5,
		CAPE:     45.0,
		FairPE:   22.22,
		Premium:  102.5, // 高溢价
		ForexResult: &forex.ExchangeResult{
			From:       "USD",
			To:         "CNY",
			Rate:       7.2345,
			Amount:     1,
			Converted:  7.2345,
			UpdateDate: "2025-01-01",
		},
	}

	message := formatDailyMessage(report)

	// 高溢价应该显示警告
	if !strings.Contains(message, "估值偏高") {
		t.Error("高溢价时应显示 '估值偏高' 评级")
	}
}

func TestFormatDailyMessage_LowPremium(t *testing.T) {
	report := &DailyReport{
		Treasury: 4.5,
		AAA:      5.0,
		BAA:      5.5,
		CAPE:     18.0,
		FairPE:   22.22,
		Premium:  -19.0, // 低估
		ForexResult: &forex.ExchangeResult{
			From:       "EUR",
			To:         "USD",
			Rate:       1.08,
			Amount:     100,
			Converted:  108.0,
			UpdateDate: "2025-01-01",
		},
	}

	message := formatDailyMessage(report)

	// 低溢价应该显示估值偏低
	if !strings.Contains(message, "估值偏低") {
		t.Error("低溢价时应显示 '估值偏低' 评级")
	}
}

func TestDailyReport_Fields(t *testing.T) {
	report := &DailyReport{
		Treasury: 4.5,
		AAA:      5.0,
		BAA:      5.5,
		CAPE:     30.0,
		FairPE:   22.22,
		Premium:  35.0,
		ForexResult: &forex.ExchangeResult{
			From:       "USD",
			To:         "CNY",
			Rate:       7.2345,
			Amount:     1,
			Converted:  7.2345,
			UpdateDate: "2025-01-01",
		},
	}

	// 验证字段值
	if report.Treasury != 4.5 {
		t.Errorf("Treasury = %v, want 4.5", report.Treasury)
	}
	if report.AAA != 5.0 {
		t.Errorf("AAA = %v, want 5.0", report.AAA)
	}
	if report.BAA != 5.5 {
		t.Errorf("BAA = %v, want 5.5", report.BAA)
	}
	if report.CAPE != 30.0 {
		t.Errorf("CAPE = %v, want 30.0", report.CAPE)
	}
	if report.ForexResult.Rate != 7.2345 {
		t.Errorf("ForexResult.Rate = %v, want 7.2345", report.ForexResult.Rate)
	}
}
