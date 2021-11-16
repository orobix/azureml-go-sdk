package workspace

import "fmt"

type ResourceNotFoundError struct {
	resourceType       string
	resourceIdentifier string
}

func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", e.resourceType, e.resourceIdentifier)
}

type HttpResponseError struct {
	statusCode      int
	responseContent string
}

func (e HttpResponseError) Error() string {
	return fmt.Sprintf("HTTP Response is in error [status code %d]: %s", e.statusCode, e.responseContent)
}
