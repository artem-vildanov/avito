package errors

//func HandleDbError[T any](err error) (*T, *AppError) {
//	if err == nil {
//		return nil, nil
//	}
//	if err.Error() == sql.ErrNoRows.Error() {
//		return nil, NotFoundErr
//	}
//	return nil, DatabaseInternalError(err.Error())
//}
//
//var NotFoundErr *AppError = &AppError{
//	Message: "not found",
//	Code:    http.StatusNotFound,
//}
//
//var CreateErr *AppError = &AppError{
//	Message: "failed to create",
//	Code:    http.StatusUnprocessableEntity,
//}
//
//var UpdateErr *AppError = &AppError{
//	Message: "failed to update",
//	Code:    http.StatusUnprocessableEntity,
//}
//
//func DatabaseInternalError(errMsg string) *AppError {
//	err := &AppError{
//		Message: fmt.Sprintf("Internal error: %s", errMsg),
//		Code:    http.StatusInternalServerError,
//	}
//	println(err)
//	return err
//}
