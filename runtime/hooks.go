package runtime

import "context"

type HookFn func(ctx context.Context) error

type Hooks struct {
	startHooks []HookFn
	stopHooks  []HookFn
}

func NewHooks() *Hooks {
	return &Hooks{
		startHooks: []HookFn{},
		stopHooks:  []HookFn{},
	}
}

func (h *Hooks) OnStart(f HookFn) {
	h.startHooks = append(h.startHooks, f)
}

func (h *Hooks) OnStop(f HookFn) {
	h.stopHooks = append(h.stopHooks, f)
}

func (h *Hooks) start(ctx context.Context) error {
	for _, hook := range h.startHooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hooks) stop(ctx context.Context) error {
	for _, hook := range h.stopHooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}
