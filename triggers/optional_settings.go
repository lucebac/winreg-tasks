// SPDX-License-Identifier: MIT

package triggers

import (
	"github.com/google/uuid"
	"github.com/lucebac/winreg-tasks/generated"
	"github.com/lucebac/winreg-tasks/utils"
)

type Privilege int64

const (
	SeCreateTokenPrivilege                    Privilege = 0x4
	SeAssignPrimaryTokenPrivilege             Privilege = 0x8
	SeLockMemoryPrivilege                     Privilege = 0x10
	SeIncreaseQuotaPrivilege                  Privilege = 0x20
	SeMachineAccountPrivilege                 Privilege = 0x40
	SeTcbPrivilege                            Privilege = 0x80
	SeSecurityPrivilege                       Privilege = 0x100
	SeTakeOwnershipPrivilege                  Privilege = 0x200
	SeLoadDriverPrivilege                     Privilege = 0x400
	SeSystemProfilePrivilege                  Privilege = 0x800
	SeSystemtimePrivilege                     Privilege = 0x1000
	SeProfileSingleProcessPrivilege           Privilege = 0x2000
	SeIncreaseBasePriorityPrivilege           Privilege = 0x4000
	SeCreatePagefilePrivilege                 Privilege = 0x8000
	SeCreatePermanentPrivilege                Privilege = 0x10000
	SeBackupPrivilege                         Privilege = 0x20000
	SeRestorePrivilege                        Privilege = 0x40000
	SeShutdownPrivilege                       Privilege = 0x80000
	SeDebugPrivilege                          Privilege = 0x100000
	SeAuditPrivilege                          Privilege = 0x200000
	SeSystemEnvironmentPrivilege              Privilege = 0x400000
	SeChangeNotifyPrivilege                   Privilege = 0x800000
	SeRemoteShutdownPrivilege                 Privilege = 0x1000000
	SeUndockPrivilege                         Privilege = 0x2000000
	SeSyncAgentPrivilege                      Privilege = 0x4000000
	SeEnableDelegationPrivilege               Privilege = 0x8000000
	SeManageVolumePrivilege                   Privilege = 0x10000000
	SeImpersonatePrivilege                    Privilege = 0x20000000
	SeCreateGlobalPrivilege                   Privilege = 0x40000000
	SeTrustedCredManAccessPrivilege           Privilege = 0x80000000
	SeRelabelPrivilege                        Privilege = 0x100000000
	SeIncreaseWorkingSetPrivilege             Privilege = 0x200000000
	SeTimeZonePrivilege                       Privilege = 0x400000000
	SeCreateSymbolicLinkPrivilege             Privilege = 0x800000000
	SeDelegateSessionUserImpersonatePrivilege Privilege = 0x1000000000
)

type OptionalSettings struct {
	length uint32

	IdleDuration            utils.Duration
	IdleWaitTimeout         utils.Duration
	ExecutionTimeLimit      utils.Duration
	DeleteExpiredTaskAfter  utils.Duration
	Priority                uint32
	RestartOnFailureDelay   utils.Duration
	RestartOnFailureRetries uint32
	NetworkId               uuid.UUID
	Privileges              Privilege
	Periodicity             utils.Duration
	Deadline                utils.Duration
	Exclusive               bool
}

func NewOptionalSettings(gen *generated.OptionalSettings) (*OptionalSettings, error) {
	optionalSettings := &OptionalSettings{length: gen.Len.Value}

	if gen.Len.Value == 0 {
		return optionalSettings, nil
	}

	optionalSettings.IdleDuration = utils.SecondsToDuration(gen.IdleDurationSeconds)
	optionalSettings.IdleWaitTimeout = utils.SecondsToDuration(gen.IdleWaitTimeoutSeconds)
	optionalSettings.ExecutionTimeLimit = utils.SecondsToDuration(gen.ExecutionTimeLimitSeconds)
	optionalSettings.DeleteExpiredTaskAfter = utils.SecondsToDuration(gen.DeleteExpiredTaskAfter)
	optionalSettings.Priority = gen.Priority
	optionalSettings.RestartOnFailureDelay = utils.SecondsToDuration(gen.RestartOnFailureDelay)
	optionalSettings.RestartOnFailureRetries = gen.RestartOnFailureRetries

	networkId, err := uuid.FromBytes(gen.NetworkId)
	if err != nil {
		return nil, err
	}
	optionalSettings.NetworkId = networkId

	if optionalSettings.length < 0x38 {
		return optionalSettings, nil
	}

	optionalSettings.Privileges = Privilege(gen.Privileges)

	if optionalSettings.length < 0x58 {
		return optionalSettings, nil
	}

	optionalSettings.Periodicity = utils.DurationFromTSTimePeriod(gen.Periodicity)
	optionalSettings.Deadline = utils.DurationFromTSTimePeriod(gen.Deadline)
	optionalSettings.Exclusive = gen.Exclusive != 0

	return optionalSettings, nil
}
