package data

import (
	"log/slog"

	"github.com/sonal3323/form-poc/types"
)

type checkoutPrice struct {
	Amount           float32                `json:"amount"`
	BillingPeriod    types.BillingMode      `json:"billing_period"`
	ResponsePerMonth types.ResponseLimitNum `json:"response_per_month"`
	Discount         float32                `json:"discount"`
	Total            float32                `json:"total"`
	TaxExcluded      bool                   `json:"tax_excluded"`
}

func (s *PostgresStore) PriceCalculation(v *types.PriceCalculationVariables) (*checkoutPrice, error) {

	planPrice, err := s.getPricing(v.PlanId)
	if err != nil {
		slog.Error("Unable to get price for price calculation", "details", err.Error())
		return nil, err
	}

	responseLimit, err := s.getResponseLimit(v.ResponseLimitID)
	if err != nil {
		slog.Error("Unable to get respone limit for price calculation", "details", err.Error())
		return nil, err
	}

	var total float32 = (planPrice.Price - planPrice.Features.Discount) + (float32(responseLimit.ResponseLimit) * responseLimit.PPR)

	return &checkoutPrice{
		Amount:           planPrice.Price,
		BillingPeriod:    planPrice.BillingPeriod,
		ResponsePerMonth: responseLimit.ResponseLimit,
		Discount:         planPrice.Discount,
		Total:            total,
		TaxExcluded:      true,
	}, nil

}
