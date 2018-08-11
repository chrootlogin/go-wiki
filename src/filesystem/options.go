package filesystem

type Option func(*filesystem) error

func WithChroot(chrootDir string) Option {
	return func(fs *filesystem) error {
		newFs, err := fs.Filesystem.Chroot(chrootDir)
		if err != nil {
			return err
		}

		fs.Chroot = chrootDir
		fs.Filesystem = newFs

		return nil
	}
}