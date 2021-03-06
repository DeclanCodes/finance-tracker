package repositories

import (
	"database/sql"

	"github.com/DeclanCodes/finance-tracker/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PortfolioRepository is the means for interacting with Portfolio storage.
type PortfolioRepository struct{}

// CreatePortfolios creates Portfolio entities in db.
func (r *PortfolioRepository) CreatePortfolios(db *sqlx.DB, ps []*models.Portfolio) ([]uuid.UUID, error) {
	query := `
	INSERT INTO portfolio (
		portfolio_uuid,
		name,
		description
	)
	VALUES (
		:portfolio_uuid,
		:name,
		:description
	)
	RETURNING portfolio_uuid;`

	IDs, err := createAndGetIDs(db, query, ps)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}

// CreatePortfolioHoldingMappings creates PortfolioHoldingMapping entities in db.
func (r *PortfolioRepository) CreatePortfolioHoldingMappings(db *sqlx.DB, phms []*models.PortfolioHoldingMapping) ([]uuid.UUID, error) {
	query := `
	INSERT INTO portfolio_holding_mapping (
		portfolio_holding_mapping_uuid,
		portfolio_uuid,
		holding_uuid
	)
	VALUES (
		:portfolio_holding_mapping_uuid,
		:portfolio.portfolio_uuid,
		:holding.holding_uuid
	)
	RETURNING portfolio_holding_mapping_uuid;`

	IDs, err := createAndGetIDs(db, query, phms)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}

// CreatePortfolioAssetCategoryMappings creates PortfolioAssetCategoryMapping entities in db.
func (r *PortfolioRepository) CreatePortfolioAssetCategoryMappings(db *sqlx.DB, pacms []*models.PortfolioAssetCategoryMapping) ([]uuid.UUID, error) {
	query := `
	INSERT INTO portfolio_asset_category_mapping (
		portfolio_asset_category_mapping_uuid,
		portfolio_uuid,
		asset_category_uuid,
		percentage
	)
	VALUES (
		:portfolio_asset_category_mapping_uuid,
		:portfolio.portfolio_uuid,
		:asset_category.asset_category_uuid,
		:percentage
	)
	RETURNING portfolio_asset_category_mapping_uuid;`

	IDs, err := createAndGetIDs(db, query, pacms)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}

// GetPortfolio retrieves Portfolio with pID from db.
func (r *PortfolioRepository) GetPortfolio(db *sqlx.DB, pID uuid.UUID) (*models.Portfolio, error) {
	mValues := map[string]interface{}{
		"portfolio": pID.String(),
	}

	ps, err := r.GetPortfolios(db, mValues)
	if err != nil {
		return nil, err
	}
	return ps[0], nil
}

// GetPortfolios gets Portfolios from db.
func (r *PortfolioRepository) GetPortfolios(db *sqlx.DB, mValues map[string]interface{}) ([]*models.Portfolio, error) {
	query := `
	SELECT
		portfolio.portfolio_uuid,
		portfolio.name,
		portfolio.description
	FROM portfolio`

	mFilters := map[string]string{
		"portfolio": "portfolio.portfolio_uuid = ",
	}

	q, args, err := getGetQueryAndValues(query, mValues, mFilters)
	if err != nil {
		return nil, err
	}

	var ps []*models.Portfolio
	err = db.Select(&ps, q, args...)
	if err != nil {
		return nil, err
	}
	if len(ps) == 0 {
		return nil, sql.ErrNoRows
	}
	return ps, nil
}

// GetPortfolioHoldingMapping retrieves PortfolioHoldingMapping with phmID from db.
func (r *PortfolioRepository) GetPortfolioHoldingMapping(db *sqlx.DB, phmID uuid.UUID) (*models.PortfolioHoldingMapping, error) {
	mValues := map[string]interface{}{
		"mapping": phmID.String(),
	}

	phms, err := r.GetPortfolioHoldingMappings(db, mValues)
	if err != nil {
		return nil, err
	}
	return phms[0], nil
}

// GetPortfolioHoldingMappings gets PortfolioHoldingMappings from db.
func (r *PortfolioRepository) GetPortfolioHoldingMappings(db *sqlx.DB, mValues map[string]interface{}) ([]*models.PortfolioHoldingMapping, error) {
	query := `
	SELECT
		portfolio_holding_mapping.portfolio_holding_mapping_uuid,
		portfolio.portfolio_uuid AS "portfolio.portfolio_uuid",
		portfolio.name AS "portfolio.name",
		portfolio.description AS "portfolio.description",
		holding.holding_uuid AS "holding.holding_uuid",
		account.account_uuid AS "account.account_uuid",
		account_category.account_category_uuid AS "account_category.account_category_uuid",
		account_category.name AS "account_category.name",
		account_category.description AS "account_category.description",
		account.name AS "account.name",
		account.description AS "account.description",
		account.amount AS "account.amount",
		fund.fund_uuid AS "fund.fund_uuid",
		asset_category.asset_category_uuid AS "asset_category.asset_category_uuid",
		asset_category.name AS "asset_category.name",
		asset_category.description AS "asset_category.description",
		fund.name AS "fund.name",
		fund.ticker_symbol AS "fund.ticker_symbol",
		fund.share_price AS "fund.share_price",
		fund.expense_ratio AS "fund.expense_ratio",
		holding.shares
	FROM portfolio_holding_mapping
	INNER JOIN portfolio
		ON portfolio_holding_mapping.portfolio_uuid = portfolio.portfolio_uuid
	INNER JOIN holding
		ON portfolio_holding_mapping.holding_uuid = holding.holding_uuid
	INNER JOIN account
		ON holding.account_uuid = account.account_uuid
	INNER JOIN account_category
		ON account.account_category_uuid = account_category.account_category_uuid
	INNER JOIN fund
		ON holding.fund_uuid = fund.fund_uuid
	INNER JOIN asset_category
		ON fund.asset_category_uuid = asset_category.asset_category_uuid`

	mFilters := map[string]string{
		"mapping":    "portfolio_holding_mapping.portfolio_holding_mapping_uuid = ",
		"portfolios": "portfolio.portfolio_uuid IN ",
	}

	q, args, err := getGetQueryAndValues(query, mValues, mFilters)
	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phms []*models.PortfolioHoldingMapping
	for rows.Next() {
		var phm models.PortfolioHoldingMapping

		err = rows.Scan(&phm.ID,
			&phm.Portfolio.ID, &phm.Portfolio.Name, &phm.Portfolio.Description,
			&phm.Holding.ID,
			&phm.Holding.Account.ID,
			&phm.Holding.Account.Category.ID, &phm.Holding.Account.Category.Name, &phm.Holding.Account.Category.Description,
			&phm.Holding.Account.Name, &phm.Holding.Account.Description, &phm.Holding.Account.Amount,
			&phm.Holding.Fund.ID,
			&phm.Holding.Fund.Category.ID, &phm.Holding.Fund.Category.Name, &phm.Holding.Fund.Category.Description,
			&phm.Holding.Fund.Name, &phm.Holding.Fund.TickerSymbol, &phm.Holding.Fund.SharePrice, &phm.Holding.Fund.ExpenseRatio,
			&phm.Holding.Shares)

		phms = append(phms, &phm)
	}
	return phms, nil
}

// GetPortfolioAssetCategoryMapping retrieves PortfolioAssetCategoryMapping with pacmID from db.
func (r *PortfolioRepository) GetPortfolioAssetCategoryMapping(db *sqlx.DB, pacmID uuid.UUID) (*models.PortfolioAssetCategoryMapping, error) {
	mValues := map[string]interface{}{
		"mapping": pacmID.String(),
	}

	pacms, err := r.GetPortfolioAssetCategoryMappings(db, mValues)
	if err != nil {
		return nil, err
	}
	return pacms[0], nil
}

// GetPortfolioAssetCategoryMappings gets PortfolioAssetCategoryMappings from db.
func (r *PortfolioRepository) GetPortfolioAssetCategoryMappings(db *sqlx.DB, mValues map[string]interface{}) ([]*models.PortfolioAssetCategoryMapping, error) {
	query := `
	SELECT
		portfolio_asset_category_mapping.portfolio_asset_category_mapping_uuid,
		portfolio.portfolio_uuid AS "portfolio.portfolio_uuid",
		portfolio.name AS "portfolio.name",
		portfolio.description AS "portfolio.description",
		asset_category.asset_category_uuid AS "asset_category.asset_category_uuid",
		asset_category.name AS "asset_category.name",
		asset_category.description AS "asset_category.description",
		portfolio_asset_category_mapping.percentage
	FROM portfolio_asset_category_mapping
	INNER JOIN portfolio
		ON portfolio_asset_category_mapping.portfolio_uuid = portfolio.portfolio_uuid
	INNER JOIN asset_category
		ON portfolio_asset_category_mapping.asset_category_uuid = asset_category.asset_category_uuid`

	mFilters := map[string]string{
		"mapping":    "portfolio_asset_category_mapping.portfolio_asset_category_mapping_uuid = ",
		"portfolios": "portfolio.portfolio_uuid IN ",
	}

	q, args, err := getGetQueryAndValues(query, mValues, mFilters)
	if err != nil {
		return nil, err
	}

	var pacms []*models.PortfolioAssetCategoryMapping
	err = db.Select(&pacms, q, args...)
	if err != nil {
		return nil, err
	}
	return pacms, nil
}

// UpdatePortfolio updates a Portfolio in db.
func (r *PortfolioRepository) UpdatePortfolio(db *sqlx.DB, p *models.Portfolio) error {
	query := `
	UPDATE portfolio
	SET
		name = :name,
		description = :description
	WHERE
		portfolio_uuid = :portfolio_uuid;`

	return updateEntity(db, query, p)
}

// UpdatePortfolioHoldingMapping updates a PortfolioHoldingMapping in db.
func (r *PortfolioRepository) UpdatePortfolioHoldingMapping(db *sqlx.DB, phm *models.PortfolioHoldingMapping) error {
	query := `
	UPDATE portfolio_holding_mapping
	SET
		portfolio_uuid = :portfolio_uuid,
		holding_uuid = :holding_uuid
	WHERE
		portfolio_holding_mapping_uuid = :portfolio_holding_mapping_uuid;`

	return updateEntity(db, query, phm)
}

// UpdatePortfolioAssetCategoryMapping updates a PortfolioAssetCategoryMapping in db.
func (r *PortfolioRepository) UpdatePortfolioAssetCategoryMapping(db *sqlx.DB, pacm *models.PortfolioAssetCategoryMapping) error {
	query := `
	UPDATE portfolio_asset_category_mapping
	SET
		portfolio_uuid = :portfolio_uuid,
		asset_category_uuid = :asset_category_uuid,
		percentage = :percentage
	WHERE
		portfolio_asset_category_mapping_uuid = :portfolio_asset_category_mapping_uuid;`

	return updateEntity(db, query, pacm)
}

// DeletePortfolio deletes a Portfolio from db.
func (r *PortfolioRepository) DeletePortfolio(db *sqlx.DB, pID uuid.UUID) error {
	query := `
	DELETE FROM portfolio
	WHERE
		portfolio_uuid = $1;`

	return deleteEntity(db, query, pID)
}

// DeletePortfolioHoldingMapping deletes a PortfolioHoldingMapping from db.
func (r *PortfolioRepository) DeletePortfolioHoldingMapping(db *sqlx.DB, phmID uuid.UUID) error {
	query := `
	DELETE FROM portfolio_holding_mapping
	WHERE
		portfolio_holding_mapping_uuid = $1;`

	return deleteEntity(db, query, phmID)
}

// DeletePortfolioAssetCategoryMapping deletes a PortfolioAssetCategoryMapping from db.
func (r *PortfolioRepository) DeletePortfolioAssetCategoryMapping(db *sqlx.DB, pacmID uuid.UUID) error {
	query := `
	DELETE FROM portfolio_asset_category_mapping
	WHERE
		portfolio_asset_category_mapping_uuid = $1;`

	return deleteEntity(db, query, pacmID)
}
