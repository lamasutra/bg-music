package model

type LocalFile interface {
	func GetPath() string
}