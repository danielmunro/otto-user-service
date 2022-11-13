/*
 * Otto user service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// Session struct for Session
type Session struct {
	User  *User  `json:"user"`
	Token string `json:"token,omitempty"`
}

func CreateSession(user *User, token string) *Session {
	return &Session{
		User:  user,
		Token: token,
	}
}
