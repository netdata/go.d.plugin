// SPDX-License-Identifier: GPL-3.0-or-later

package job

type Job interface {
	Name() string
	ModuleName() string
	FullName() string
	AutoDetection() bool
	AutoDetectionEvery() int
	RetryAutoDetection() bool
	Tick(clock int)
	Start()
	Stop()
	Cleanup()
}
