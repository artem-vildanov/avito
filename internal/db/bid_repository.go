package db

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/models"
)

type BidRepository struct {
	PostgresStorage
}

func NewBidRepository(storage *PostgresStorage) *BidRepository {
	return &BidRepository{
		PostgresStorage: *storage,
	}
}

func (self BidRepository) GetBidById(id string) (*models.BidDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`SELECT * FROM bid 
		WHERE id = $1
		limit 1;`,
		id,
	)
	return models.NewBidDbModel(row, errors.BidNotFound(id))
}

func (self BidRepository) GetBidsListByTenderId(tenderId string, limit, offset uint) ([]*models.BidDbModel, *errors.AppError) {
	rows, err := self.Database.Query(
		`SELECT * FROM bid 
		WHERE tender_id = $1
		ORDER BY name ASC
		limit $2
		offset $3`,
		tenderId,
		limit,
		offset,
	)
	if err != nil {
		return nil, errors.DatabaseError
	}
	defer rows.Close()
	return models.NewBidDbModelsList(rows)
}

func (self BidRepository) GetBidsListByUsername(username string, limit, offset uint) ([]*models.BidDbModel, *errors.AppError) {
	rows, err := self.Database.Query(
		`SELECT bid.id, 
		bid.name, 
		bid.description, 
		bid.status,
		bid.tender_id,
		bid.author_id,
		bid.author_type,
		bid.version,
		bid.created_at
		FROM bid 
        join employee on bid.author_id = employee.id
		WHERE employee.username = $1
		ORDER BY bid.name ASC
		limit $2
		offset $3;`,
		username,
		limit,
		offset,
	)
	if err != nil {
		println(err.Error())
		return nil, errors.DatabaseError
	}
	defer rows.Close()
	return models.NewBidDbModelsList(rows)
}

func (self BidRepository) CreateBid(createModel *models.BidCreateModel) (*models.BidDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`insert into bid 
		(name, description, tender_id, author_id, author_type) 
		values ($1, $2, $3, $4, $5)
		returning *;`,
		createModel.Name,
		createModel.Description,
		createModel.TenderId,
		createModel.AuthorId,
		string(createModel.AuthorType),
	)
	return models.NewBidDbModel(row, errors.FailedToCreateBid)
}

func (self BidRepository) UpdateBidParams(
	id string,
	updateModel *models.BidUpdateModel,
) (*models.BidDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`update bid
		set name = $1, description = $2, version = version + 1
		where id = $3
		returning *`,
		updateModel.Name,
		updateModel.Description,
		id,
	)
	return models.NewBidDbModel(row, errors.FailedToUpdateBid(id))
}

func (self BidRepository) UpdateBidStatus(
	id string, status enums.BidStatus,
) (*models.BidDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`update bid
		set status = $1, version = version + 1
		where id = $2
		returning *`,
		string(status),
		id,
	)
	return models.NewBidDbModel(row, errors.FailedToUpdateBid(id))
}

func (self BidRepository) CancelBidsByTenderId(tenderId string) ([]*models.BidDbModel, *errors.AppError) {
	rows, err := self.Database.Query(
		`UPDATE bid
		SET status = $1, version = version + 1  
		WHERE tender_id = $2 AND status != 'Approved'
		returning *;`,
		enums.BidStatusCanceled,
		tenderId,
	)
	if err != nil {
		return nil, errors.DatabaseError
	}
	return models.NewBidDbModelsList(rows)
}

func (self BidRepository) RollbackBid(rollbackModel *models.BidRollbackDbModel) (*models.BidDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`update bid set
		name = $1,
		description = $2,
		status = $3, 
		tender_id = $4,
		author_id = $5,
		author_type = $6,
		version = version + 1
		where id = $7
		returning *;`,
		rollbackModel.Name,
		rollbackModel.Description,
		rollbackModel.Status,
		rollbackModel.TenderId,
		rollbackModel.AuthorId,
		rollbackModel.AuthorType,
		rollbackModel.BidId,
	)
	return models.NewBidDbModel(row, errors.FailedToUpdateBid(rollbackModel.BidId))
}
