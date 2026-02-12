package mongodb

type DocUser struct {
	Name  string `bson:"name,omitempty"`
	Age   int    `bson:"age,omitempty"`
	Email string `bson:"email,omitempty"`
}
