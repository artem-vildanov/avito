package db

import (
	"avito/internal/errors"
	"avito/internal/models"
)

type BidRollbackRepository struct {
	PostgresStorage
}

func NewBidRollbackRepository(db *PostgresStorage) *BidRollbackRepository {
	return &BidRollbackRepository{*db}
}

func (self BidRollbackRepository) SaveBidRollback(rollbackModel *models.BidDbModel) *errors.AppError {
	_, err := self.Database.Exec(
		`insert into bid_rollback (
			bid_id,
			name,
			description,
			status,
			tender_id,
			author_id,
			author_type,
			version,
			created_at
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		rollbackModel.Id,
		rollbackModel.Name,
		rollbackModel.Description,
		rollbackModel.Status,
		rollbackModel.TenderId,
		rollbackModel.AuthorId,
		rollbackModel.AuthorType,
		rollbackModel.Version,
		rollbackModel.CreatedAt,
	)
	if err != nil {
		return errors.FailedToSaveBidRollback
	}
	return nil
}

func (self BidRollbackRepository) SaveBidRollbacksList(rollbackModels []*models.BidDbModel) *errors.AppError {
	// Begin a new transaction
	tx, err := self.Database.Begin()
	if err != nil {
		return errors.DatabaseError
	}

	// Prepare the insert statement
	stmt, err := tx.Prepare(`
		INSERT INTO bid_rollback (
			bid_id,
			name,
			description,
			status,
			tender_id,
			author_id,
			author_type,
			version,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`)
	if err != nil {
		tx.Rollback()
		return errors.DatabaseError
	}
	defer stmt.Close()

	// Execute the prepared statement for each rollback model
	for _, rollbackModel := range rollbackModels {
		_, err := stmt.Exec(
			rollbackModel.Id,
			rollbackModel.Name,
			rollbackModel.Description,
			rollbackModel.Status,
			rollbackModel.TenderId,
			rollbackModel.AuthorId,
			rollbackModel.AuthorType,
			rollbackModel.Version,
			rollbackModel.CreatedAt,
		)
		if err != nil {
			tx.Rollback()
			return errors.FailedToSaveBidRollback
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.DatabaseError
	}

	return nil
}

func (self BidRollbackRepository) GetBidRollback(id string, version uint) (*models.BidRollbackDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`select * from bid_rollback 
		where bid_id = $1 and version = $2
		limit 1;`,
		id, version,
	)
	return models.NewBidRollbackDbModel(row, errors.BidRollbackNotFound(id, version))
}
