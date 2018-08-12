package auth

import (
	"log"
	"errors"
	"encoding/json"

	"github.com/patrickmn/go-cache"

	"github.com/chrootlogin/go-wiki/src/lib/repo"
	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
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

	// remove users from cache
	authCache.Delete("users")

	return nil
}

func GetUserList() (*UserList, error) {
	// check if users are in cache
	cachedUsers, found := authCache.Get("users")
	if found {
		userList := &UserList{
			Users: cachedUsers.(map[string]common.User),
		}

		return userList, nil
	}

	// otherwise read from disk
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

	// add to cache
	authCache.Set("users", users, cache.DefaultExpiration)

	var userList = &UserList{
		Users: users,
	}

	return userList, nil
}