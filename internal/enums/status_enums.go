package enums

type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"
	TenderStatusPublished TenderStatus = "Published"
	TenderStatusClosed    TenderStatus = "Closed"
)

func GetTenderStatuses() []TenderStatus {
	return []TenderStatus{
		TenderStatusCreated,
		TenderStatusPublished,
		TenderStatusClosed,
	}
}

type BidStatus string

const (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusApproved  BidStatus = "Approved"
)

func GetBidStatuses() []BidStatus {
	return []BidStatus{
		BidStatusCreated,
		BidStatusPublished,
		BidStatusCanceled,
	}
}

type Decision string

const (
	DecisionApproved Decision = "Approved"
	DecisionRejected Decision = "Rejected"
)

func GetDecisions() []Decision {
	return []Decision{
		DecisionApproved,
		DecisionRejected,
	}
}
