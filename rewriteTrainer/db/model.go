package db

import (
	"gatoor/orca/base"
	"time"
)

type UtilisationStatistic struct {
	Timestamp  time.Time
	Cpu        base.Usage
	Mbytes     base.Usage
	Network    base.Usage

	AppName    base.AppName
	AppVersion base.Version
	Host       base.HostId
}

type ApplicationUtilisationStatistic struct {
	Cpu        base.Usage
	Mbytes     base.Usage
	Network    base.Usage
	AppName    base.AppName
	Timestamp  time.Time
}

type ApplicationCountStatistic struct {
	Timestamp  time.Time
	AppName    base.AppName
	Running    base.DeploymentCount
	Desired    base.DeploymentCount
}

type AuditEvent struct {
	Timestamp time.Time
	Details map[string]string
}
