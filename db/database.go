package database

// RepositoryReady defines methods needed to this program
type RepositoryReady interface {
	IsAlive() error
}
