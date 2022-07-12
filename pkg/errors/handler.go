package errors

func NoError() RsError {
	return RsError{}
}

func ConvertErrToRsErr(err error) *RsError {
	var rsErr *RsError
	switch err {
	//case ErrNotFound:
	//	rsErr = NewNotFoundErr()
	//	break
	case ErrRecordExist:
		rsErr = NewRecordExistErr()
		break
	default:
		rsErr = NewInternalServerErr(err)
	}

	return rsErr
}
