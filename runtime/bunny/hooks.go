package bunny

// HookPoint is the point in time at which we hook
type HookPoint int

// the hook point constants
const (
	BeforeInsertHook HookPoint = iota + 1
	BeforeUpdateHook
	BeforeDeleteHook
	AfterInsertHook
	AfterSelectHook
	AfterUpdateHook
	AfterDeleteHook
)
