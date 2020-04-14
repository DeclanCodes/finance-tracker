package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Contribution is .
type Contribution struct {
	ContributionUUID uuid.UUID       `json:"contributionUuid,omitEmpty" db:"contribution_uuid"`
	Account          Account         `json:"account" db:"account"`
	Name             string          `json:"name" db:"name"`
	Description      string          `json:"description" db:"description"`
	Date             time.Time       `json:"date" db:"date_made"`
	Amount           decimal.Decimal `json:"amount" db:"amount"`
}
