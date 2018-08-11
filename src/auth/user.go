package auth

import (
	"log"
	"errors"
	"encoding/json"

	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
	"github.com/chrootlogin/go-wiki/src/filesystem"
)

type UserList struct {
	Users map[string]common.User
}

func (u *UserList) GetAll() map[string]common.User {
	return u.Users
}

func (u *UserList) Get(username string) (common.User, error) {
	value, ok := u.Users[username]
	if ok {
		value.Username = username
		return value, nil
	} else {
		return common.User{}, errors.New("User not found: " + username)
	}
}

func (u *UserList) Add(user common.User) error {
	u.Users[user.Username] = user

	jsonData, err := json.Marshal(u.Users)
	if err != nil {
		log.Println("marshal json: " + err.Error())
		return err
	}

	err = filesystem.New().Commit("prefs/_users.json", filesystem.File{Content:string(jsonData)}, repo.Commit{})
	if err != nil {
		log.Println("save file: " + err.Error())
		return err
	}

	return nil
}

func GetUserList() (*UserList, error) {
	usersRaw, err := filesystem.New().Get("prefs/_users.json")
	if err != nil {
		log.Fatal(err)
	}

	// Convert json to object
	var users map[string]common.User
	err = json.Unmarshal([]byte(usersRaw.Content), &users)
	if err != nil {
		log.Println("unmarshal: " + err.Error())
		return nil, err
	}

	var userList = &UserList{
		Users: users,
	}

	return userList, nil
}

func GetUserById(userId string) (common.User, error) {
	ul, err := GetUserList()
	if err != nil {
		return common.User{}, err
	}
	u, err := ul.Get(userId)
	if err != nil {
		return common.User{}, err
	}

	return u, err
}