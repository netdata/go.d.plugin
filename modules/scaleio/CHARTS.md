`ScaleIO` collector produces the following charts:

Capacity:
  1. **System Total Capacity** in KB
    * total
    
  2. **System Capacity** in KB
    * protected
    * degraded
    * spare
    * failed
    * decreased
    * unavailable
    * in_maintenance
    * unused
    
  3. **Volume Usage By Type** in KB
    * thick
    * thin
    
  4. **Thin Volume Usage** in KB
    * free
    * used

I/O Workload:
  1. **Primary Backend Bandwidth Total (Read and Write)** in KB/s
    * total
    
  2. **Primary Backend Bandwidth** in KB/s
    * read
    * write
    
  3. **Primary Backend IOPS Total (Read and Write)** in iops/s
    * total
    
  4. **Primary Backend I/O Size Total (Read and Write)** in KB
    * io_size
    
  5. **Primary Backend IOPS** in iops/s
    * read
    * write

Rebalance:
  1. **System Rebalance** in KB/s
    * read
    * write
    
  2. **System Rebalance Pending Capacity** in KB
    * left
    
Rebuild:
  1. **System Rebuild Bandwidth Total (Forward, Backward and Normal)** in KB/s
    * read
    * write
    
  2. **System Rebuild Pending Capacity Total (Forward, Backward and Normal)** in KB
    * left
    
Components:
  1. **System Defined Components** in number
    * devices
    * fault_sets
    * protection_domains
    * rfcache_devices
    * scsi_initiators
    * sdc
    * sds
    * snapshots
    * storage_pools
    * volumes
    * vtrees
    
  2. **Volumes By Type** in number
    * thick
    * thin
    
  3. **Volumes By Mapping** in number
    * mapped
    * unmapped
