package enums

type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"
	TenderStatusPublished TenderStatus = "Published"
	TenderStatusClosed    TenderStatus = "Closed"
)

var TenderStatusesList = []string{
	string(TenderStatusCreated),
	string(TenderStatusPublished),
	string(TenderStatusClosed),
}

type BidStatus string

const (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusApproved  BidStatus = "Approved"
)

var BidStatusesList = []string{
	string(BidStatusCreated),
	string(BidStatusPublished),
	string(BidStatusCanceled),
}

type Decision string

const (
	DecisionApproved Decision = "Approved"
	DecisionRejected Decision = "Rejected"
)

var DecisionsList = []string{
	string(DecisionApproved),
	string(DecisionRejected),
}
