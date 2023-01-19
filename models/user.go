package models

type User struct {
	ApiKey                string   `firestore:"API-KEY"`
	AllowedObjectPatterns []string `firestore:"allowed_object_patterns"`
}
