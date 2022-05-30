package csv

type User struct {
	Id, Name, LastName, Email, Phone string
	FriendIds                        []string
}
