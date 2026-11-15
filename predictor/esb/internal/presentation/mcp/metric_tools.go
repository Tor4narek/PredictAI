package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tor4narek/predictai/esb/internal/application/analytics"
	"github.com/tor4narek/predictai/esb/internal/domain"
)

func registerMetricTools(server *mcp.Server, allowed map[string]struct{}, service *analytics.Service) {
	addTool(server, allowed, &mcp.Tool{Name: "save_metric_observations", Description: "Сохраняет наблюдения метрик, связанные с новостями и фактами. Значения и вероятности передаются числами JSON или числовыми строками; ESB преобразует строки в числа. Для неизвестных значений используйте null. Границы периода у наблюдений не хранятся. ESB формирует стабильный ID из компании, метрики, новости и исходных фактов; если разные прогнозы используют одни и те же ID фактов, передайте разные observation_id."}, func(ctx context.Context, _ *mcp.CallToolRequest, input metricObservationsInput) (*mcp.CallToolResult, insertedOutput, error) {
		observations := make([]domain.MetricObservation, 0, len(input.Observations))
		for i, observation := range input.Observations {
			if observation.CompanyID != "" && !strings.EqualFold(strings.TrimSpace(observation.CompanyID), strings.TrimSpace(input.CompanyID)) {
				return nil, insertedOutput{}, fmt.Errorf("observation %d company_id must match the request company_id", i)
			}
			var probability *float64
			if observation.Probability != nil {
				parsed, err := numericValue(observation.Probability)
				if err != nil {
					return nil, insertedOutput{}, fmt.Errorf("observation %d probability must be a number or null: %w", i, err)
				}
				probability = &parsed
			}
			var value *float64
			if observation.Value != nil {
				parsed, err := numericValue(observation.Value)
				if err != nil {
					return nil, insertedOutput{}, fmt.Errorf("observation %d value must be a number or null: %w", i, err)
				}
				value = &parsed
			}
			observations = append(observations, domain.MetricObservation{ID: observation.ID, CompanyID: input.CompanyID, MetricID: observation.MetricID, NewsID: observation.NewsID,
				SourceID: observation.SourceID, SourceURL: observation.SourceURL, MaterialKind: observation.MaterialKind, Value: value,
				Probability: probability, SourceFactIDs: observation.SourceFactIDs})
		}
		if err := service.SaveMetricObservations(ctx, input.CompanyID, observations); err != nil {
			return nil, insertedOutput{}, err
		}
		return nil, insertedOutput{Inserted: len(observations)}, nil
	})

	addTool(server, allowed, &mcp.Tool{Name: "calculate_metric_total", Description: "Пересчитывает метрику по наблюдениям ClickHouse для включительного диапазона времени публикации новостей и сохраняет снимок результата в ClickHouse. Принадлежность метрики компании и единица измерения берутся из настроек PostgreSQL."}, func(ctx context.Context, _ *mcp.CallToolRequest, input metricTotalInput) (*mcp.CallToolResult, domain.MetricTotal, error) {
		total, err := service.CalculateMetricTotal(ctx, input.MetricID, input.PeriodFrom, input.PeriodTo)
		if err != nil {
			return nil, domain.MetricTotal{}, err
		}
		return nil, total, nil
	})

	addTool(server, allowed, &mcp.Tool{Name: "calculate_all_metric_totals", Description: "Читает настройки и единицы измерения метрик компании из PostgreSQL, пересчитывает каждую метрику по наблюдениям ClickHouse, сохраняет аналитические снимки в ClickHouse и возвращает результаты в порядке каталога."}, func(ctx context.Context, _ *mcp.CallToolRequest, input periodInput) (*mcp.CallToolResult, domain.CompanyMetricTotals, error) {
		totals, err := service.CalculateAllMetricTotals(ctx, input.CompanyID, input.PeriodFrom, input.PeriodTo)
		if err != nil {
			return nil, domain.CompanyMetricTotals{}, err
		}
		return nil, totals, nil
	})

}
