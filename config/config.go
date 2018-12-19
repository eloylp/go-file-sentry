package config

type Version struct {
	storagePath  string
	watchedFiles []string
}

func (v *Version) WatchedFiles() []string {
	return v.watchedFiles
}

func NewVersion(storagePath string, watchedFiles []string) *Version {
	return &Version{storagePath: storagePath, watchedFiles: watchedFiles}
}

func (v *Version) StoragePath() string {
	return v.storagePath
}

func (v *Version) SetStoragePath(storagePath string) {
	v.storagePath = storagePath
}
