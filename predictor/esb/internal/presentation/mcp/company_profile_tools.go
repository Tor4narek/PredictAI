package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tor4narek/predictai/esb/internal/application/companycontext"
)

func registerCompanyProfileTools(server *mcp.Server, allowed map[string]struct{}, companyContext *companycontext.Service) {
	addTool(server, allowed, &mcp.Tool{Name: "get_company_profile", Description: "Читает профиль компании в формате Markdown из S3 через ESB."}, func(ctx context.Context, _ *mcp.CallToolRequest, input companyProfileInput) (*mcp.CallToolResult, companyProfileOutput, error) {
		markdown, key, err := companyContext.CompanyProfile(ctx, input.CompanyID)
		if err != nil {
			return nil, companyProfileOutput{}, err
		}
		return nil, companyProfileOutput{CompanyID: strings.TrimSpace(input.CompanyID), Key: key, Markdown: markdown}, nil
	})

	addTool(server, allowed, &mcp.Tool{Name: "save_company_profile_and_metrics", Description: "Идемпотентно сохраняет профиль компании в формате Markdown в S3 и обновляет её структурированные настройки метрик в PostgreSQL через ESB."}, func(ctx context.Context, _ *mcp.CallToolRequest, input saveCompanyContextInput) (*mcp.CallToolResult, companyContextOutput, error) {
		key, err := companyContext.SaveCompanyProfileAndMetrics(ctx, input.CompanyID, input.ProfileRevision, input.Markdown, input.Metrics)
		if err != nil {
			return nil, companyContextOutput{}, err
		}
		return nil, companyContextOutput{CompanyID: strings.TrimSpace(input.CompanyID), ProfileRevision: input.ProfileRevision, ProfileKey: key, MetricsSaved: len(input.Metrics)}, nil
	})

	addTool(server, allowed, &mcp.Tool{Name: "save_company_metrics", Description: "Идемпотентно обновляет структурированные настройки метрик компании в PostgreSQL через ESB."}, func(ctx context.Context, _ *mcp.CallToolRequest, input companyMetricsInput) (*mcp.CallToolResult, insertedOutput, error) {
		if err := companyContext.SaveCompanyMetrics(ctx, input.CompanyID, input.Metrics); err != nil {
			return nil, insertedOutput{}, err
		}
		return nil, insertedOutput{Inserted: len(input.Metrics)}, nil
	})

	addTool(server, allowed, &mcp.Tool{Name: "get_company_metrics", Description: "Читает структурированные настройки метрик компании из PostgreSQL через ESB. Параметр limit принимает целое число или числовую строку."}, func(ctx context.Context, _ *mcp.CallToolRequest, input companyLimitInput) (*mcp.CallToolResult, companyMetricsOutput, error) {
		limit, err := integerValue(input.Limit)
		if err != nil {
			return nil, companyMetricsOutput{}, fmt.Errorf("limit must be an integer: %w", err)
		}
		metrics, err := companyContext.CompanyMetrics(ctx, input.CompanyID, limit)
		if err != nil {
			return nil, companyMetricsOutput{}, err
		}
		return nil, companyMetricsOutput{Metrics: metrics}, nil
	})

}
