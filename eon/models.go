package eon

import (
	"strings"
	"time"
)

// FlexibleTime handles JSON date fields that may be empty strings, null, or valid timestamps
type FlexibleTime struct {
	time.Time
}

func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`)

	// Handle empty string or null
	if s == "" || s == "null" {
		ft.Time = time.Time{}
		return nil
	}

	// Try multiple time formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // Without timezone
		time.RFC3339Nano,
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			ft.Time = t
			return nil
		}
	}

	return err
}

func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + ft.Time.Format(time.RFC3339) + `"`), nil
}

// Installations - based on swagger InstallationsWrapper and InstallationDto

type InstallationsWrapper struct {
	Installations []InstallationDto `json:"installations"`
}

type InstallationDto struct {
	ID                          string   `json:"id"`
	Active                      bool     `json:"active"`
	Address                     string   `json:"address"`
	Business                    string   `json:"business"`
	Category                    string   `json:"category"`
	City                        string   `json:"city"`
	EnergyClass                 string   `json:"energyClass"`
	GridArea                    string   `json:"gridArea"`
	Name                        string   `json:"name"`
	OrgNumber                   string   `json:"orgNumber"`
	PriceArea                   string   `json:"priceArea"`
	Resolution                  string   `json:"resolution"`
	SafetyLevel                 *float32 `json:"safetyLevel"`
	HasMeasurementsSubscription bool     `json:"hasMeasurementsSubscription"`
	HasCostsSubscription        bool     `json:"hasCostsSubscription"`
}

// Measurement Series - based on swagger InstallationsMeasurementsWrapper

type InstallationsMeasurementsWrapper struct {
	Installations []InstallationMeasurementsDto `json:"installations"`
}

type InstallationMeasurementsDto struct {
	ID                string                 `json:"id"`
	MeasurementSeries []MeasurementSeriesDto `json:"measurementSeries"`
}

type MeasurementSeriesDto struct {
	ID         int          `json:"id"`
	SeriesType string       `json:"seriesType"`
	Unit       string       `json:"unit"`
	LastUpdate FlexibleTime `json:"lastUpdate"`
}

// Measurements - based on swagger MeasurementsWrapper and MeasurementDto

type MeasurementsWrapper struct {
	ID           int              `json:"id"`
	Resolution   string           `json:"resolution"`
	Measurements []MeasurementDto `json:"measurements"`
}

type MeasurementDto struct {
	TimeStamp FlexibleTime `json:"timeStamp"`
	Value     *float64     `json:"value"` // nullable
}

// Costs - based on swagger CostsWrapper and various cost DTOs
// The API returns different schemas based on energy type (oneOf)

type CostsWrapper struct {
	EnergyClass  string `json:"energyClass"`
	Installation string `json:"installation"`
}

type CostsElectricityWrapper struct {
	CostsWrapper
	Costs []CostElectricityProductionDto `json:"costs"`
}

type CostsProductionWrapper struct {
	CostsWrapper
	Costs []CostElectricityProductionDto `json:"costs"`
}

type CostsHeatWrapper struct {
	CostsWrapper
	Costs []CostHeatColdDto `json:"costs"`
}

type CostsColdWrapper struct {
	CostsWrapper
	Costs []CostHeatColdDto `json:"costs"`
}

type CostsGasWrapper struct {
	CostsWrapper
	Costs []CostGasDto `json:"costs"`
}

type CostsBaseDto struct {
	Month time.Time `json:"month"`
}

type CostElectricityProductionDto struct {
	CostsBaseDto
	RetailCost      *float64            `json:"retailCost"`
	RetailCostVAT   *float64            `json:"retailCostVAT"`
	EnergyTax       *float64            `json:"energyTax"`
	EnergyTaxVAT    *float64            `json:"energyTaxVAT"`
	NetCost         *float64            `json:"netCost"`
	NetCostVAT      *float64            `json:"netCostVAT"`
	CostGridDetails *CostGridDetailsDto `json:"costGridDetails"`
}

type CostGasDto struct {
	CostsBaseDto
	RetailCost        *float64              `json:"retailCost"`
	RetailCostVAT     *float64              `json:"retailCostVAT"`
	EnergyTax         *float64              `json:"energyTax"`
	EnergyTaxVAT      *float64              `json:"energyTaxVAT"`
	CostBioGasDetails *CostBioGasDetailsDto `json:"costBioGasDetails"`
}

type CostHeatColdDto struct {
	CostsBaseDto
	RetailCost    *float64 `json:"retailCost"`
	RetailCostVAT *float64 `json:"retailCostVAT"`
	EffectCost    *float64 `json:"effectCost"`
	EffectCostVAT *float64 `json:"effectCostVAT"`
	EnergyCost    *float64 `json:"energyCost"`
	EnergyCostVAT *float64 `json:"energyCostVAT"`
	FlowCost      *float64 `json:"flowCost"`
	FlowCostVAT   *float64 `json:"flowCostVAT"`
}

type CostGridDetailsDto struct {
	GridSubscription                   *float64 `json:"gridSubscription"`
	GridSubscriptionVAT                *float64 `json:"gridSubscriptionVAT"`
	GridSubscribedEffectReactiveIn     *float64 `json:"gridSubscribedEffectReactiveIn"`
	GridSubscribedEffectReactiveInVAT  *float64 `json:"gridSubscribedEffectReactiveInVAT"`
	GridSubscribedEffectReactiveOut    *float64 `json:"gridSubscribedEffectReactiveOut"`
	GridSubscribedEffectReactiveOutVAT *float64 `json:"gridSubscribedEffectReactiveOutVAT"`
	GridSubscribedEffectWinter         *float64 `json:"gridSubscribedEffectWinter"`
	GridSubscribedEffectWinterVAT      *float64 `json:"gridSubscribedEffectWinterVAT"`
	GridSubscribedEffect               *float64 `json:"gridSubscribedEffect"`
	GridSubscribedEffectVAT            *float64 `json:"gridSubscribedEffectVAT"`
	GridEffectCompensation             *float64 `json:"gridEffectCompensation"`
	GridEffectCompensationVAT          *float64 `json:"gridEffectCompensationVAT"`
	GridEffect                         *float64 `json:"gridEffect"`
	GridEffectVAT                      *float64 `json:"gridEffectVAT"`
	GridCompensationEnergy             *float64 `json:"gridCompensationEnergy"`
	GridCompensationEnergyVAT          *float64 `json:"gridCompensationEnergyVAT"`
	GridExceededReactiveEffectOut      *float64 `json:"gridExceededReactiveEffectOut"`
	GridExceededReactiveEffectOutVAT   *float64 `json:"gridExceededReactiveEffectOutVAT"`
	GridCompensationLoss               *float64 `json:"gridCompensationLoss"`
	GridCompensationLossVAT            *float64 `json:"gridCompensationLossVAT"`
	GridOther                          *float64 `json:"gridOther"`
	GridOtherVAT                       *float64 `json:"gridOtherVAT"`
	GridExceededActiveEffect           *float64 `json:"gridExceededActiveEffect"`
	GridExceededActiveEffectVAT        *float64 `json:"gridExceededActiveEffectVAT"`
	GridFixed                          *float64 `json:"gridFixed"`
	GridFixedVAT                       *float64 `json:"gridFixedVAT"`
	GridExceededReactiveEffect         *float64 `json:"gridExceededReactiveEffect"`
	GridExceededReactiveEffectVAT      *float64 `json:"gridExceededReactiveEffectVAT"`
	GridTransfer                       *float64 `json:"gridTransfer"`
	GridTransferVAT                    *float64 `json:"gridTransferVAT"`
}

type CostBioGasDetailsDto struct {
	BioGasCarbonDioxideTax    *float64 `json:"bioGasCarbonDioxideTax"`
	BioGasCarbonDioxideTaxVAT *float64 `json:"bioGasCarbonDioxideTaxVAT"`
	BiogasEnergyTax           *float64 `json:"biogasEnergyTax"`
	BiogasEnergyTaxVAT        *float64 `json:"biogasEnergyTaxVAT"`
	BiogasAccumulatedTax      *float64 `json:"biogasAccumulatedTax"`
	BiogasAccumulatedTaxVAT   *float64 `json:"biogasAccumulatedTaxVAT"`
}
