package common

var DefaultFiles = map[string]interface{}{
	"prefs/_config.json": map[string]string{
		"title": "Go-Wiki",
		"registration": "0",
	},
	"prefs/_users.json": map[string]User{
		"admin": {
			PasswordHash: "$2a$14$4.xZgCG40zWYult3N1qrK.kid2RKqmKgjbOaEDqg8noS1FAdx6v5O",
			Email: "admin@example.com",
			IsEnabled: true,
			Permissions: []string{
				"ADMINISTRATOR",
			},
		},
	},
}
