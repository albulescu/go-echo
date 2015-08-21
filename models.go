package main

type User struct {
	Id   string `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
}
