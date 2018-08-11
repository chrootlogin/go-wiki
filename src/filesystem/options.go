package filesystem

type Option func(*filesystem) error

func WithChroot(chrootDir string) Option {
	return func(fs *filesystem) error {
		newFs, err := fs.fs.Chroot(chrootDir)
		if err != nil {
			return err
		}

		fs.fs = newFs
		return nil
	}
}

/*type Options struct {
	UID         int
	GID         int
	Flags       int
	Contents    string
	Permissions os.FileMode
}



func UID(userID int) Option {
	return func(args *Options) {
		args.UID = userID
	}
}

func GID(groupID int) Option {
	return func(args *Options) {
		args.GID = groupID
	}
}

func Contents(c string) Option {
	return func(args *Options) {
		args.Contents = c
	}
}

func Permissions(perms os.FileMode) Option {
	return func(args *Options) {
		args.Permissions = perms
	}
}*/