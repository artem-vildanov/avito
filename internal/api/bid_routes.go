package api

import (
	"avito/internal/db"
	"avito/internal/handlers"
)

func (self Router) BidRoutes(storage *db.PostgresStorage) {
	bh := handlers.NewBidHandler(storage)
	rh := handlers.NewReviewHandler(storage)
	r := "/bids"
	self.createRoute(r+"/new", bh.CreateBid).Methods("POST")
	self.createRoute(r+"/my", bh.GetBidsListByUsername).Methods("GET")
	self.createRoute(r+"/{tenderId}/list", bh.GetBidsListByTenderId).Methods("GET")
	self.createRoute(r+"/{bidId}/status", bh.GetBidStatus).Methods("GET")
	self.createRoute(r+"/{bidId}/status", bh.UpdateBidStatus).Methods("PUT")
	self.createRoute(r+"/{bidId}/edit", bh.UpdateBidParams).Methods("PATCH")
	self.createRoute(r+"/{bidId}/submit_decision", bh.SubmitDecision).Methods("PUT")
	self.createRoute(r+"/{bidId}/rollback/{version}", bh.RollbackBid).Methods("PUT")
	self.createRoute(r+"/{bidId}/feedback", rh.LeaveFeedback).Methods("PUT")
	self.createRoute(r+"/{tenderId}/reviews", rh.GetReviews).Methods("GET")
}
