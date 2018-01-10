package models

import (
	"time"
)

type Stock struct {
	Name              string  `json:"name"`
	Code              string  `json:"code"`
	Description       string  `json:"description"`
	IssueTs           int64   `json:"issue_ts"`
	TotalSupply       float64 `json:"total_supploy"`
	CirculatingSupply float64 `json:"circulating_supply"`
	Reference         string  `json:"reference"`
}

func NewStock(name, code, desc string, total, circul float64, ref string) *Stock {
	return &Stock{
		Name:              name,
		Code:              code,
		Description:       desc,
		IssueTs:           time.Now().Unix(),
		TotalSupply:       total,
		CirculatingSupply: circul,
		Reference:         ref,
	}
}

func (stock *Stock) RefLink() *Link {
	return NewLink(stock.Reference)
}
