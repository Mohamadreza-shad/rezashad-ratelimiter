package config

type Window struct {
	Size int
}

func WindowSize() int {
	if cfg.Window.Size == 0 {
		return 3
	}
	return cfg.Window.Size
}
