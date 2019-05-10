package scaleio

import (
	"github.com/netdata/go.d.plugin/modules/scaleio/client"
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

type selectedStatistics struct {
	System struct {
		NumOfDevices                             float64
		NumOfFaultSets                           float64
		NumOfMappedToAllVolumes                  float64
		NumOfProtectionDomains                   float64
		NumOfRfcacheDevices                      float64
		NumOfScsiInitiators                      float64
		NumOfSdc                                 float64
		NumOfSds                                 float64
		NumOfSnapshots                           float64
		NumOfStoragePools                        float64
		NumOfThickBaseVolumes                    float64
		NumOfThinBaseVolumes                     float64
		NumOfUnmappedVolumes                     float64
		NumOfVolumes                             float64
		NumOfVolumesInDeletion                   float64
		NumOfVtrees                              float64
		ActiveBckRebuildCapacityInKb             float64
		ActiveFwdRebuildCapacityInKb             float64
		ActiveMovingCapacityInKb                 float64
		ActiveNormRebuildCapacityInKb            float64
		ActiveRebalanceCapacityInKb              float64
		AtRestCapacityInKb                       float64
		BckRebuildCapacityInKb                   float64
		CapacityAvailableForVolumeAllocationInKb float64
		CapacityInUseInKb                        float64
		CapacityLimitInKb                        float64
		DegradedFailedCapacityInKb               float64
		DegradedHealthyCapacityInKb              float64
		FailedCapacityInKb                       float64
		FwdRebuildCapacityInKb                   float64
		InMaintenanceCapacityInKb                float64
		MaxCapacityInKb                          float64
		MovingCapacityInKb                       float64
		NormRebuildCapacityInKb                  float64
		PendingBckRebuildCapacityInKb            float64
		PendingFwdRebuildCapacityInKb            float64
		PendingMovingCapacityInKb                float64
		PendingNormRebuildCapacityInKb           float64
		PendingRebalanceCapacityInKb             float64
		ProtectedCapacityInKb                    float64
		RebalanceCapacityInKb                    float64
		SemiProtectedCapacityInKb                float64
		SnapCapacityInUseInKb                    float64
		SnapCapacityInUseOccupiedInKb            float64
		SpareCapacityInKb                        float64
		ThickCapacityInUseInKb                   float64
		ThinCapacityAllocatedInKb                float64
		ThinCapacityInUseInKb                    float64
		UnreachableUnusedCapacityInKb            float64
		UnusedCapacityInKb                       float64
		NormRebuildReadBwc                       client.Bwc // TODO: ???
		NormRebuildWriteBwc                      client.Bwc // TODO: ???
		BckRebuildReadBwc                        client.Bwc // failed node/disk is back alive
		BckRebuildWriteBwc                       client.Bwc // failed node/disk is back alive
		FwdRebuildReadBwc                        client.Bwc // node/disk fails
		FwdRebuildWriteBwc                       client.Bwc // node/disk fails
		RebalanceReadBwc                         client.Bwc
		RebalanceWriteBwc                        client.Bwc
		PrimaryReadBwc                           client.Bwc // Backend (SDSs + Devices)
		PrimaryWriteBwc                          client.Bwc // Backend (SDSs + Devices)
		SecondaryReadBwc                         client.Bwc // Backend (SDSs + Devices, 2nd)
		SecondaryWriteBwc                        client.Bwc // Backend (SDSs + Devices, 2nd)
		UserDataReadBwc                          client.Bwc // Frontend (Volumes + SDCs)
		UserDataWriteBwc                         client.Bwc // Frontend (Volumes + SDCs)
		PrimaryReadFromDevBwc                    client.Bwc // TODO: ???
		PrimaryReadFromRmcacheBwc                client.Bwc // TODO: ???
		SecondaryReadFromDevBwc                  client.Bwc // TODO: ???
		SecondaryReadFromRmcacheBwc              client.Bwc // TODO: ???
		TotalReadBwc                             client.Bwc // *ReadBwc
		TotalWriteBwc                            client.Bwc // *WriteBwc
	}
}

func (s *ScaleIO) collect() (map[string]int64, error) {
	var mx metrics

	err := s.collectSystemOverview(&mx)
	if err != nil {
		return nil, err
	}

	return stm.ToMap(mx), nil
}

func (s *ScaleIO) collectSystemOverview(mx *metrics) error {
	var stats selectedStatistics

	err := s.apiClient.GetSelectedStatistics(&stats, selectedStatisticsQuery)
	if err != nil {
		return err
	}

	collectSystemCapacity(mx, stats)
	collectSystemComponents(mx, stats)
	collectSystemIOWorkload(mx, stats)
	collectSystemRebalance(mx, stats)
	collectSystemRebuild(mx, stats)

	return nil
}

func collectSystemCapacity(mx *metrics, stats selectedStatistics) {
	var (
		m = &mx.SystemOverview.Capacity
		s = stats.System
	)

	// Physical Capacity Calculation (as in the ScaleIO GUI)
	{
		m.AvailableForVolumeAllocation.Set(s.CapacityAvailableForVolumeAllocationInKb)
		m.MaxCapacity.Set(s.MaxCapacityInKb)

		// TODO: investigate decreased calculation, current implementation maybe wrong under some conditions
		m.Decreased.Set(sum(s.MaxCapacityInKb, -s.CapacityLimitInKb))
		m.Degraded.Set(sum(s.DegradedFailedCapacityInKb, s.DegradedHealthyCapacityInKb))
		m.Failed.Set(s.FailedCapacityInKb)
		m.InMaintenance.Set(s.InMaintenanceCapacityInKb)
		// TODO: ProtectedCapacityInKb + SemiProtectedCapacityInKb?
		m.Protected.Set(s.ProtectedCapacityInKb)
		m.Spare.Set(s.SpareCapacityInKb)
		m.UnreachableUnused.Set(s.UnreachableUnusedCapacityInKb)

		// Note: can't use 'UnusedCapacityInKb' directly, dashboard shows calculated value
		// TODO: some other values?
		used := sum(
			s.ProtectedCapacityInKb,
			s.InMaintenanceCapacityInKb,
			m.Decreased.Value(),
			m.Degraded.Value(),
			s.FailedCapacityInKb,
			s.SpareCapacityInKb,
			s.UnreachableUnusedCapacityInKb,
		)
		m.Unused.Set(sum(s.MaxCapacityInKb, -used))
	}

	// TODO: do we need this?
	{
		m.Limit.Set(s.CapacityLimitInKb)
		m.InUse.Set(s.CapacityInUseInKb)
		m.SemiProtected.Set(s.SemiProtectedCapacityInKb)
		m.SnapInUse.Set(s.SnapCapacityInUseInKb)
		m.SnapInUseOccupied.Set(s.SnapCapacityInUseOccupiedInKb)
	}

	m.ThickInUse.Set(s.ThickCapacityInUseInKb)
	m.ThinAllocated.Set(s.ThinCapacityAllocatedInKb)
	m.ThinInUse.Set(s.ThinCapacityInUseInKb)
	m.ThinFree.Set(sum(s.ThinCapacityAllocatedInKb, -s.ThinCapacityInUseInKb))
}

func collectSystemComponents(mx *metrics, stats selectedStatistics) {
	var (
		m = &mx.SystemOverview.Components
		s = stats.System
	)

	m.Devices.Set(s.NumOfDevices)
	m.FaultSets.Set(s.NumOfFaultSets)
	m.MappedToAllVolumes.Set(s.NumOfMappedToAllVolumes)
	m.ProtectionDomains.Set(s.NumOfProtectionDomains)
	m.RfcacheDevices.Set(s.NumOfRfcacheDevices)
	m.ScsiInitiators.Set(s.NumOfScsiInitiators)
	m.Sdc.Set(s.NumOfSdc)
	m.Sds.Set(s.NumOfSds)
	m.Snapshots.Set(s.NumOfSnapshots)
	m.StoragePools.Set(s.NumOfStoragePools)
	m.Vtrees.Set(s.NumOfVtrees)
	m.Volumes.Set(s.NumOfVolumes)

	m.ThickBaseVolumes.Set(s.NumOfThickBaseVolumes)
	m.ThinBaseVolumes.Set(s.NumOfThinBaseVolumes)
	m.UnmappedVolumes.Set(s.NumOfUnmappedVolumes)
	m.VolumesInDeletion.Set(s.NumOfVolumesInDeletion)
	m.MappedVolumes.Set(sum(s.NumOfVolumes, -s.NumOfUnmappedVolumes))
}

func collectSystemIOWorkload(mx *metrics, stats selectedStatistics) {
	var (
		m = &mx.SystemOverview.IOWorkload
		s = stats.System
	)

	// --Bandwidth--

	m.Total.BW.set(
		calcBW(s.TotalReadBwc),
		calcBW(s.TotalWriteBwc),
	)
	m.Frontend.BW.set(
		calcBW(s.UserDataReadBwc),
		calcBW(s.UserDataWriteBwc),
	)
	m.Backend.Primary.BW.set(
		calcBW(s.PrimaryReadBwc),
		calcBW(s.PrimaryWriteBwc),
	)
	m.Backend.Secondary.BW.set(
		calcBW(s.SecondaryReadBwc),
		calcBW(s.SecondaryWriteBwc),
	)
	m.Backend.Total.BW.set(
		sumGauge(m.Backend.Primary.BW.Read, m.Backend.Secondary.BW.Read),
		sumGauge(m.Backend.Primary.BW.Write, m.Backend.Secondary.BW.Write),
	)

	// --IOPS--

	m.Total.IOPS.set(
		calcIOPS(s.TotalReadBwc),
		calcIOPS(s.TotalWriteBwc),
	)
	m.Frontend.IOPS.set(
		calcIOPS(s.UserDataReadBwc),
		calcIOPS(s.UserDataWriteBwc),
	)
	m.Backend.Primary.IOPS.set(
		calcIOPS(s.PrimaryReadBwc),
		calcIOPS(s.PrimaryWriteBwc),
	)
	m.Backend.Secondary.IOPS.set(
		calcIOPS(s.SecondaryReadBwc),
		calcIOPS(s.SecondaryWriteBwc),
	)
	m.Backend.Total.IOPS.set(
		sumGauge(m.Backend.Primary.IOPS.Read, m.Backend.Secondary.IOPS.Read),
		sumGauge(m.Backend.Primary.IOPS.Write, m.Backend.Secondary.IOPS.Write),
	)

	// --I/O SIZE--

	m.Total.IOSize.set(
		calcIOSize(s.TotalReadBwc),
		calcIOSize(s.TotalWriteBwc),
	)
	m.Frontend.IOSize.set(
		calcIOSize(s.UserDataReadBwc),
		calcIOSize(s.UserDataWriteBwc),
	)
	m.Backend.Primary.IOSize.set(
		calcIOSize(s.PrimaryReadBwc),
		calcIOSize(s.PrimaryWriteBwc),
	)
	m.Backend.Secondary.IOSize.set(
		calcIOSize(s.SecondaryReadBwc),
		calcIOSize(s.SecondaryWriteBwc),
	)
	m.Backend.Total.IOSize.set(
		sumGauge(m.Backend.Primary.IOSize.Read, m.Backend.Secondary.IOSize.Read),
		sumGauge(m.Backend.Primary.IOSize.Write, m.Backend.Secondary.IOSize.Write),
	)
}

func collectSystemRebuild(mx *metrics, stats selectedStatistics) {
	var (
		m = &mx.SystemOverview.Rebuild
		s = stats.System
	)

	// --Bandwidth--

	m.Forward.BW.set(
		calcBW(s.FwdRebuildReadBwc),
		calcBW(s.FwdRebuildWriteBwc),
	)
	m.Backward.BW.set(
		calcBW(s.BckRebuildReadBwc),
		calcBW(s.BckRebuildWriteBwc),
	)
	m.Normal.BW.set(
		calcBW(s.NormRebuildReadBwc),
		calcBW(s.NormRebuildWriteBwc),
	)
	m.Total.BW.set(
		sumGauge(m.Forward.BW.Read, m.Backward.BW.Read, m.Normal.BW.Read),
		sumGauge(m.Forward.BW.Write, m.Backward.BW.Write, m.Normal.BW.Write),
	)

	// --IOPS--

	m.Forward.IOPS.set(
		calcIOPS(s.FwdRebuildReadBwc),
		calcIOPS(s.FwdRebuildWriteBwc),
	)
	m.Backward.IOPS.set(
		calcIOPS(s.BckRebuildReadBwc),
		calcIOPS(s.BckRebuildWriteBwc),
	)
	m.Normal.IOPS.set(
		calcIOPS(s.NormRebuildReadBwc),
		calcIOPS(s.NormRebuildWriteBwc),
	)
	m.Total.IOPS.set(
		sumGauge(m.Forward.IOPS.Read, m.Backward.IOPS.Read, m.Normal.IOPS.Read),
		sumGauge(m.Forward.IOPS.Write, m.Backward.IOPS.Write, m.Normal.IOPS.Write),
	)

	// --IO SIZE--

	m.Forward.IOSize.set(
		calcIOSize(s.FwdRebuildReadBwc),
		calcIOSize(s.FwdRebuildWriteBwc),
	)
	m.Backward.IOSize.set(
		calcIOSize(s.BckRebuildReadBwc),
		calcIOSize(s.BckRebuildWriteBwc),
	)
	m.Normal.IOSize.set(
		calcIOSize(s.NormRebuildReadBwc),
		calcIOSize(s.NormRebuildWriteBwc),
	)
	m.Total.IOSize.set(
		sumGauge(m.Forward.IOSize.Read, m.Backward.IOSize.Read, m.Normal.IOSize.Read),
		sumGauge(m.Forward.IOSize.Write, m.Backward.IOSize.Write, m.Normal.IOSize.Write),
	)

	// --Pending Capacity--

	m.Forward.Pending.Set(s.PendingFwdRebuildCapacityInKb)
	m.Backward.Pending.Set(s.PendingBckRebuildCapacityInKb)
	m.Normal.Pending.Set(s.PendingNormRebuildCapacityInKb)
	m.Total.Pending.Set(sumGauge(m.Forward.Pending, m.Backward.Pending, m.Normal.Pending))
}

func collectSystemRebalance(mx *metrics, stats selectedStatistics) {
	var (
		m = &mx.SystemOverview.Rebalance
		s = stats.System
	)

	// --Bandwidth--
	m.BW.set(
		calcBW(s.RebalanceReadBwc),
		calcBW(s.RebalanceWriteBwc),
	)

	// --IOPS--

	m.IOPS.set(
		calcIOPS(s.RebalanceReadBwc),
		calcIOPS(s.RebalanceWriteBwc),
	)

	// --IO SIZE--

	m.IOSize.set(
		calcIOSize(s.RebalanceReadBwc),
		calcIOSize(s.RebalanceWriteBwc),
	)

	// --Pending Capacity--

	m.Pending.Set(s.PendingRebalanceCapacityInKb)

}

func div(a, b int64) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func calcBW(bwc client.Bwc) float64 { return div(bwc.TotalWeightInKb, bwc.NumSeconds) }

func calcIOPS(bwc client.Bwc) float64 { return div(bwc.NumOccured, bwc.NumSeconds) }

func calcIOSize(bwc client.Bwc) float64 { return div(bwc.TotalWeightInKb, bwc.NumOccured) }

func sumGauge(vs ...mtx.Gauge) (res float64) {
	for _, v := range vs {
		res += v.Value()
	}
	return res
}

func sum(vs ...float64) (res float64) {
	for _, v := range vs {
		res += v
	}
	return res
}
