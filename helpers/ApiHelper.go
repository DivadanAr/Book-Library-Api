package helpers

import "main.go/models/response"

func GetResponse(status int, data interface{}, err error) response.Response {
	var response response.Response

	switch status {
	case 200:
		response.Message = "Success"
	default:
		response.Message = err.Error()
	}

	response.Status = status
	response.Data = data

	return response
}
