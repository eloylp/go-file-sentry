package config

type Config struct {
	storagePath  string
	watchedFiles []string
}

func (v *Config) WatchedFiles() []string {
	return v.watchedFiles
}

func NewVersion(storagePath string, watchedFiles []string) *Config {
	return &Config{storagePath: storagePath, watchedFiles: watchedFiles}
}

func (v *Config) StoragePath() string {
	return v.storagePath
}

func (v *Config) SetStoragePath(storagePath string) {
	v.storagePath = storagePath
}
