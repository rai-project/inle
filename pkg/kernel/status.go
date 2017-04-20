package kernel

// KernelStatus holds a kernel state, for status broadcast messages.
type KernelStatus struct {
	ExecutionState string `json:"execution_state"`
}
