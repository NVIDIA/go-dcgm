package dcgm

// administrationCallsWithoutSession lists public administration operations that must reject an absent session.
func administrationCallsWithoutSession() map[string]func() error {
	return map[string]func() error{
		"health":       func() error { _, err := HostengineIsHealthy(); return err },
		"logging":      func() error { return HostengineSetLoggingSeverity(0, LoggingNone) },
		"environment":  func() error { _, err := HostengineEnvironmentVariable("DCGM_DIAG_REMOVED"); return err },
		"modules":      func() error { _, err := ModuleGetStatuses(); return err },
		"module name":  func() error { _, err := ModuleIDToName(ModuleHealth); return err },
		"denylist":     func() error { return ModuleDenylist(ModuleHealth) },
		"groups":       func() error { _, err := GroupGetAllIDs(); return err },
		"remove GPU":   func() error { return RemoveFromGroup(GroupHandle{}, 0) },
		"field group":  func() error { _, err := FieldGroupGetInfo(FieldHandle{}); return err },
		"field groups": func() error { _, err := FieldGroupGetAll(); return err },
		"topology":     func() error { _, err := GetGroupTopology(GroupHandle{}); return err },
		"architecture": func() error { _, err := GetGPUChipArchitecture(0); return err },
		"selection":    func() error { _, err := SelectGPUsByTopology(nil, 1, 0); return err },
	}
}
