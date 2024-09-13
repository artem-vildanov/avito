package db

import "avito/internal/errors"

type BidApproveRepository struct {
	PostgresStorage
}

func NewBidApproveRepository(db *PostgresStorage) *BidApproveRepository {
	return &BidApproveRepository{*db}
}

func (self BidApproveRepository) AddApprove(bidId, employeeId string) *errors.AppError {
	_, err := self.Database.Exec(
		`insert into bid_approve (bid_id, employee_id) VALUES ($1, $2);`,
		bidId, employeeId,
	)
	if err != nil {
		return errors.DatabaseError
	}
	return nil
}

func (self BidApproveRepository) RemoveApprovesByBidId(bidId string) *errors.AppError {
	_, err := self.Database.Exec(
		`DELETE FROM bid_approve WHERE bid_id = $1;`,
		bidId,
	)
	if err != nil {
		return errors.BidNotFound(bidId)
	}
	return nil
}

func (self BidApproveRepository) CheckEmployeeApprovedBid(bidId, employeeId string) *errors.AppError {
	var exists bool
	err := self.Database.QueryRow(
		`SELECT EXISTS (
            SELECT 1 
            FROM bid_approve 
            WHERE employee_id = $1 and bid_id = $2
        );`,
		employeeId,
		bidId,
	).Scan(&exists)

	if err != nil {
		println(err.Error())
		return errors.DatabaseError
	}

	if exists {
		return errors.AlreadyApprovedBid(bidId, employeeId)
	}

	return nil
}

func (self BidApproveRepository) CountApprovementsByBidId(bidId string) (uint, *errors.AppError) {
	var count uint
	err := self.Database.QueryRow(
		`SELECT COUNT(*) 
        FROM bid_approve 
        WHERE bid_id = $1;`,
		bidId,
	).Scan(&count)

	if err != nil {
		return 0, errors.DatabaseError
	}

	return count, nil
}
