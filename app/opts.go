package app

type AppOpt func(*App)

func WithPort(port int) AppOpt {
	return func(a *App) {
		a.mainPort = port
	}
}

func WithMetricsPort(port int) AppOpt {
	return func(a *App) {
		a.metricsPort = port
	}
}

func WithPprofPort(port int) AppOpt {
	return func(a *App) {
		a.pprofPort = port
	}
}

func WithPprofEnabled(enabled bool) AppOpt {
	return func(a *App) {
		a.pprofEnabled = enabled
	}
}
