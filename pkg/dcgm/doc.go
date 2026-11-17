// Package dcgm provides Go bindings for NVIDIA Data Center GPU Manager.
//
// Initialize a session with Init and release it with the returned cleanup.
// InitEmbedded accepts typed logging and module denylist options for a new
// embedded session; it rejects an already initialized session. Init callers
// share the active session through reference counting.
//
// HostengineIsHealthy, HostengineEnvironmentVariable, ModuleGetStatuses and
// ModuleIDToName read administration data. HostengineSetLoggingSeverity and
// ModuleDenylist change the shared host engine and affect its other clients.
//
// GroupGetAllIDs and FieldGroupGetAll list existing objects without transferring
// ownership. Destroy only groups and field groups created by the caller.
// GetGroupTopology describes group affinity; SelectGPUsByTopology chooses nearby
// GPUs using DCGM's 64-bit GPU mask, so candidate IDs must be at most 63.
// GetGPUChipArchitecture reports each GPU's architecture.
package dcgm
