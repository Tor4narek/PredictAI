package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tor4narek/predictai/esb/internal/application/analytics"
	"github.com/tor4narek/predictai/esb/internal/domain"
)

func registerReportTools(server *mcp.Server, allowed map[string]struct{}, service *analytics.Service) {
	addTool(server, allowed, &mcp.Tool{Name: "save_company_analysis", Description: "Сохраняет анализ компании с оценками вероятности. Вероятность передаётся числом JSON или числовой строкой; ESB преобразует строку в число."}, func(ctx context.Context, _ *mcp.CallToolRequest, input companyAnalysisInput) (*mcp.CallToolResult, insertedOutput, error) {
		probability, err := numericValue(input.Probability)
		if err != nil {
			return nil, insertedOutput{}, fmt.Errorf("analysis probability must be a number: %w", err)
		}
		analysis := domain.CompanyAnalysis{CompanyID: input.CompanyID, PeriodFrom: input.PeriodFrom, PeriodTo: input.PeriodTo, Probability: probability, MetricIDs: input.MetricIDs, Text: input.Text}
		if err := service.SaveCompanyAnalysis(ctx, analysis); err != nil {
			return nil, insertedOutput{}, err
		}
		return nil, insertedOutput{Inserted: 1}, nil
	})

	addTool(server, allowed, &mcp.Tool{Name: "get_company_report", Description: "Читает отчёт компании со связанными фактами, настроенными метриками, итогами расчётов и анализами."}, func(ctx context.Context, _ *mcp.CallToolRequest, input periodInput) (*mcp.CallToolResult, domain.CompanyReport, error) {
		report, err := service.CompanyReport(ctx, input.CompanyID, input.PeriodFrom, input.PeriodTo)
		if err != nil {
			return nil, domain.CompanyReport{}, err
		}
		return nil, report, nil
	})

	addTool(server, allowed, &mcp.Tool{Name: "get_consultant_context", Description: "Читает связанные факты и наблюдения из ClickHouse, а через ESB пересчитывает итоги настроенных метрик компании. Ответ содержит ID источников; итоги используют время публикации новостей за всю сохранённую историю. Параметр limit принимает целое число или числовую строку."}, func(ctx context.Context, _ *mcp.CallToolRequest, input companyLimitInput) (*mcp.CallToolResult, domain.ConsultantContext, error) {
		limit, err := integerValue(input.Limit)
		if err != nil {
			return nil, domain.ConsultantContext{}, fmt.Errorf("limit must be an integer: %w", err)
		}
		result, err := service.ConsultantContext(ctx, input.CompanyID, limit)
		if err != nil {
			return nil, domain.ConsultantContext{}, err
		}
		return nil, result, nil
	})
}
