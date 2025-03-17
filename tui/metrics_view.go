package tui

import (
	"fmt"
	"strings"
)

// MetricsView component
type MetricsView struct {
	*CustomTextView
	theme            *Theme
	validatorMetrics ValidatorMetrics
	financialMetrics FinancialMetrics
}

func NewMetricsView(theme *Theme) *MetricsView {
	return &MetricsView{
		CustomTextView: NewCustomTextView(theme),
		theme:          theme,
	}
}

// Add this method to the MetricsView struct
func (m *MetricsView) UpdateMetrics(validatorMetrics ValidatorMetrics, financialMetrics FinancialMetrics) {
	var content strings.Builder

	// Validator Metrics
	content.WriteString(fmt.Sprintf("[#%06x::b]Validator Metrics[-:-:-]\n", m.theme.Colors.Primary.Hex()))
	content.WriteString(fmt.Sprintf("%s Session Time: %s\n", string(m.theme.Symbols.Bullet), validatorMetrics.SessionTime))
	content.WriteString(fmt.Sprintf("%s Solutions: [#%06x]%d/%d (%.1f%%)[-:-:-]\n",
		string(m.theme.Symbols.Bullet),
		m.theme.Colors.Success.Hex(),
		validatorMetrics.SolutionsLastMinute.Success,
		validatorMetrics.SolutionsLastMinute.Total,
		validatorMetrics.SolutionsLastMinute.Rate*100))
	content.WriteString(fmt.Sprintf("%s Avg Solutions/min: %.2f\n\n", string(m.theme.Symbols.Bullet), validatorMetrics.AverageSolutionsPerMin))

	// Financial Metrics
	content.WriteString(fmt.Sprintf("[#%06x::b]Financial Metrics[-:-:-]\n", m.theme.Colors.Secondary.Hex()))
	content.WriteString(fmt.Sprintf("%s Token/min: %.4f\n", string(m.theme.Symbols.Bullet), financialMetrics.TokenIncomePerMin))
	content.WriteString(fmt.Sprintf("%s USD/hour: $%.2f\n", string(m.theme.Symbols.Bullet), financialMetrics.IncomePerHour))
	content.WriteString(fmt.Sprintf("%s Profit/day: [#%06x]$%.2f[-:-:-]",
		string(m.theme.Symbols.Bullet),
		m.theme.Colors.Success.Hex(),
		financialMetrics.ProfitPerDay))

	m.SetText(content.String())
}
