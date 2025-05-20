package models

type SocialNetwork struct {
	ID                int    `json: "id"`
	UserID            string `json: "user_id"`
	SocialNetworkType string `json: "social_network_type"`
	LinkAccess        string `json: "link_access"`
	IsActive          bool   `json: "is_active"`
}
