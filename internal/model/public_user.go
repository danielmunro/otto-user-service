package model

import (
"time"
)

type PublicUser struct {

	Uuid string `json:"uuid"`

	Name string `json:"name,omitempty"`

	Username string `json:"username"`

	ProfilePic string `json:"profile_pic,omitempty"`

	BioMessage string `json:"bio_message,omitempty"`

	Birthday string `json:"birthday,omitempty"`

	AddressStreet string `json:"address_street,omitempty"`

	AddressCity string `json:"address_city,omitempty"`

	AddressZip string `json:"address_zip,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

