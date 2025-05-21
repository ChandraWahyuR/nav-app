package utils

type Response struct {
	Status  string
	Message string
	Data    interface{}
}

type MetadataResponse struct {
	Status   string
	Message  string
	Metadata interface{}
	Data     interface{}
}

func ResponseHandler(status, message string, data interface{}) Response {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return response
}

func MetadataFormatResponse(status string, message string, metadata interface{}, data interface{}) MetadataResponse {
	response := MetadataResponse{
		Status:   status,
		Message:  message,
		Metadata: metadata,
		Data:     data,
	}
	return response
}
