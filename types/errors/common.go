package errors

import (
	"fmt"
	"strings"
)

var (
	NotAuthenticated = &appError{
		Err:  "user not authenticated",
		Code: 401,
	}
	DecodeJSON = &appError{
		Err:  "could not decode incoming JSON",
		Code: 422,
	}
	IDMissing = &appError{
		Err:  "ID missing in URL",
		Code: 422,
	}
	NotNumber = &appError{
		Err:  "ID needs to be an integer",
		Code: 422,
	}
	ServiceNameMissing = &appError{
		Err:  "missing service name",
		Code: 422,
	}
	ServiceTypeMissing = &appError{
		Err:  "missing service type",
		Code: 422,
	}
	DomainNameMissing = &appError{
		Err:  "missing domain name",
		Code: 422,
	}
	CheckIntervalMissing = &appError{
		Err:  "missing check interval",
		Code: 422,
	}
	CommandConfigNotJson = &appError{
		Err:  "command config not JSON",
		Code: 422,
	}
	CommandConfigFieldCmdMissing = &appError{
		Err:  "command config field \"Cmd\" missing",
		Code: 422,
	}
)

func Missing(object interface{}, id interface{}) error {
	outErr := fmt.Errorf("%s with id %v was not found", splitVar(object), id)
	return &appError{
		Err:  outErr.Error(),
		Code: 404,
	}
}

func splitVar(val interface{}) string {
	s := strings.Split(fmt.Sprintf("%T", val), ".")
	return strings.ToLower(s[len(s)-1])
}
