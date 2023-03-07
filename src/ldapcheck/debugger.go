package ldapcheck

import (
	"user-check/utils/logger"
)

type Debugger struct{}

func (d *Debugger) Debug(format string, v ...interface{}) {
	log := logger.SugaredLogger().With("package", "azure", "source", "azure_debugger")
	log.Debugf(format, v...)
}
