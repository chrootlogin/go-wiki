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
	} else {
		return common.User{}, errors.New("User not found: " + username)
	}
}

func (ul *userList) Add(user common.User) error {
	// check for error
	if ul.Error != nil {
		return ul.Error
	}

	ul.Users[user.Username] = user

	jsonData, err := json.Marshal(ul.Users)
	if err != nil {
		return err
	}

	err = filesystem.New(filesystem.WithChroot("prefs")).Save("_users.json", filesystem.File{Content: string(jsonData)})
	if err != nil {
		return err
	}

	// remove users from cache
	storeCache.Set(USERS_CACHE, ul.Users, cache.NoExpiration)
	return nil
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
	storeCache.Set(USERS_CACHE, users, cache.NoExpiration)

	return &userList{
		Users: users,
	}
}
