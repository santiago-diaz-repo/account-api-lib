package models

import "time"

type ResponseObject struct {
	Data  *ResponseData `json:"data,omitempty"`
	Links *Link         `json:"links,omitempty"`
}

type ResponseData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	CreateOn       time.Time          `json:"created_on,omitempty"`
	ID             string             `json:"id,omitempty"`
	ModifiedOn     time.Time          `json:"modified_on,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type Link struct {
	Self string `json:"self,omitempty"`
}
