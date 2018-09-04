package store

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/patrickmn/go-cache"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

const USERS_CACHE = "users"

var (
	ErrUserNotExist = errors.New("user does not exist")
)

type userList struct {
	Users map[string]common.User
	Error error
}

func (ul *userList) Get(username string) (common.User, error) {
	// check for error
	if ul.Error != nil {
		return common.User{}, ul.Error
	}

	value, ok := ul.Users[username]
	if ok {
		value.Username = username
		return value, nil
	}

	return common.User{}, ErrUserNotExist
}

func (ul *userList) Add(user common.User) error {
	// check for error
	if ul.Error != nil {
		return ul.Error
	}

	ul.Users[user.Username] = user
	err := ul.save()
	if err != nil {
		return err
	}

	// write new users to cache
	storeCache.Set(USERS_CACHE, ul.Users, cache.DefaultExpiration)
	return nil
}

func (ul *userList) Delete(user common.User) {
	_, ok := ul.Users[user.Username]
	if ok {
		delete(ul.Users, user.Username)

		storeCache.Set(USERS_CACHE, ul.Users, cache.DefaultExpiration)
	}
}

func UserList() *userList {
	// check if users are in cache
	cachedUsers, found := storeCache.Get(USERS_CACHE)
	if found {
		ul := &userList{
			Users: cachedUsers.(map[string]common.User),
		}

		return ul
	}

	// otherwise read from disk
	usersRaw, err := filesystem.New(filesystem.WithChroot("prefs")).Get("_users.json")
	if err != nil {
		log.Fatal(err)
	}

	// Convert json to object
	var users map[string]common.User
	err = json.Unmarshal([]byte(usersRaw.Content), &users)
	if err != nil {
		log.Println("unmarshal: " + err.Error())

		return &userList{
			Error: err,
		}
	}

	// add to cache with no expiration
	storeCache.Set(USERS_CACHE, users, cache.DefaultExpiration)

	return &userList{
		Users: users,
	}
}

func (ul *userList) save() error {
	jsonData, err := json.Marshal(ul.Users)
	if err != nil {
		return err
	}

	err = filesystem.New(filesystem.WithChroot("prefs")).Save("_users.json", filesystem.File{Content: string(jsonData)})
	if err != nil {
		return err
	}

	return nil
}

func (ul *userList) refresh() error {
	// read from disk
	usersRaw, err := filesystem.New(filesystem.WithChroot("prefs")).Get("_users.json")
	if err != nil {
		log.Fatal(err)
	}

	// Convert json to object
	var users map[string]common.User
	err = json.Unmarshal([]byte(usersRaw.Content), &users)
	if err != nil {
		return err
	}

	ul.Users = users
	return nil
}