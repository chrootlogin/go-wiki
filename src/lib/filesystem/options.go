package filesystem

type Option func(*Filesystem) error

func WithChroot(chrootDir string) Option {
	return func(fs *Filesystem) error {
		newFs, err := fs.Filesystem.Chroot(chrootDir)
		if err != nil {
			return err
		}

		fs.ChrootDirectory = chrootDir
		fs.Filesystem = newFs

		return nil
	}
}