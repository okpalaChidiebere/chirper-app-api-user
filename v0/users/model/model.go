package usermodel

type User struct {
	Id string  `json:"id" dynamodbav:"id"`
  	Name string  `json:"name" dynamodbav:"name"`
  	AvatarURL string  `json:"avatarURL,omitempty" dynamodbav:"avatarURL,omitempty"`
  	Tweets []string  `json:"tweets,omitempty" dynamodbav:"tweets,omitempty,omitemptyelem,stringset"`
}