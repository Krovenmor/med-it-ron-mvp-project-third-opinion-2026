package jobs

import "go.uber.org/fx"

func Run(lc fx.Lifecycle, pool *Pool) {
	lc.Append(fx.StartStopHook(pool.Start, pool.Stop))
}
