package model

import (
	"time"
)

type PublicUser struct {
	Uuid string `json:"uuid"`

	Name string `json:"name"`

	Username string `json:"username"`

	ProfilePic string `json:"profile_pic"`

	BioMessage string `json:"bio_message"`

	Birthday string `json:"birthday"`

	Role Role `json:"role,omitempty"`

	AddressStreet string `json:"address_street,omitempty"`

	AddressCity string `json:"address_city,omitempty"`

	AddressZip string `json:"address_zip,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}
