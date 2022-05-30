package csv

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	FILE_NAME      = "csv/students.csv"
	MAX_GOROUTINES = 10
)

func ProcessFile() {
	f, err := os.Open(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}

	users := scanFile(f)

	// sequential processing
	// sequentialProcessing(users)

	// concurrent processing
	concurrentProcessing(users)
}

func concurrentProcessing(users []*User) {
	usersCh := make(chan []*User)
	unvisitedUsers := make(chan *User)
	go func() { usersCh <- users }()
	initializeWorkers(unvisitedUsers, usersCh, users)
	processUsers(unvisitedUsers, usersCh, len(users))
}

func initializeWorkers(unvisitedUsers <-chan *User, usersCh chan []*User, users []*User) {
	for i := 0; i < MAX_GOROUTINES; i++ {
		go func() {
			for user := range unvisitedUsers {
				sendSmsNotification(user)
				go func(user *User) {
					friendIds := user.FriendIds
					friends := []*User{}
					for _, friendId := range friendIds {
						friend, err := findUserById(friendId, users)
						if err != nil {
							fmt.Printf("Error %v\n", err)
							continue
						}
						friends = append(friends, friend)
					}

					_, ok := <-usersCh
					if ok {
						usersCh <- friends
					}
				}(user)
			}
		}()
	}
}

func processUsers(unvisitedUsers chan<- *User, usersCh chan []*User, size int) {
	visitedUsers := make(map[string]bool)
	count := 0
	for users := range usersCh {
		for _, user := range users {
			if !visitedUsers[user.Id] {
				visitedUsers[user.Id] = true
				count++
				if count >= size {
					close(usersCh)
				}
				unvisitedUsers <- user
			}
		}
	}
}

func sequentialProcessing(users []*User) {
	visited := make(map[string]bool)
	for _, user := range users {
		if !visited[user.Id] {
			visited[user.Id] = true
			sendSmsNotification(user)
			for _, friendId := range user.FriendIds {
				friend, err := findUserById(friendId, users)
				if err != nil {
					fmt.Printf("Error %v\n", err)
					continue
				}

				if !visited[friend.Id] {
					visited[friend.Id] = true
					sendSmsNotification(friend)
				}
			}
		}
	}
}

func sendSmsNotification(user *User) {
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("Sending sms notification to %v\n", user.Phone)
}

func findUserById(userId string, users []*User) (*User, error) {
	for _, user := range users {
		if user.Id == userId {
			return user, nil
		}
	}

	return nil, fmt.Errorf("User not found with id %v", userId)
}

func scanFile(f *os.File) []*User {
	s := bufio.NewScanner(f)
	users := []*User{}
	for s.Scan() {
		line := strings.Trim(s.Text(), " ")
		lineArray := strings.Split(line, ",")
		ids := strings.Split(lineArray[5], " ")
		ids = ids[1 : len(ids)-1]
		user := &User{
			Id:        lineArray[0],
			Name:      lineArray[1],
			LastName:  lineArray[2],
			Email:     lineArray[3],
			Phone:     lineArray[4],
			FriendIds: ids,
		}
		users = append(users, user)
	}
	return users
}
