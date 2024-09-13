package api

import (
	"avito/internal/db"
	"avito/internal/handlers"
)

func (self Router) TenderRoutes(storage *db.PostgresStorage) {
	h := handlers.NewTenderHandler(storage)
	r := "/tenders"
	self.createRoute(r, h.GetTendersList).Methods("GET")
	self.createRoute(r+"/my", h.GetTendersListByUsername).Methods("GET")
	self.createRoute(r+"/new", h.CreateTender).Methods("POST")
	self.createRoute(r+"/{tenderId}/status", h.GetTenderStatus).Methods("GET")
	self.createRoute(r+"/{tenderId}/status", h.UpdateTenderStatus).Methods("PUT")
	self.createRoute(r+"/{tenderId}/edit", h.UpdateTenderParams).Methods("PATCH")
	self.createRoute(r+"/{tenderId}/rollback/{version}", h.RollbackTender).Methods("PUT")
}
