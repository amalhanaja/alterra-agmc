package app

type App interface {
	Run() error
	OnInit() error
	OnDestroy()
}
