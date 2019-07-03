package resources

/*

```
Virtual Datacenter Architecture Representation (partial).

<root>
+-DC0 # Virtual datacenter
   +-datastore # Datastore folder (created by system)
   | +-Datastore1
   |
   +-host # Host folder (created by system)
   | +-Folder1 # Host and Cluster folder
   | | +-NestedFolder1
   | | | +-Cluster1
   | | | | +-Host1
   | +-Cluster2
   | | +-Host2
   | | | +-VM1
   | | | +-VM2
   | | | +-hadoop1
   | +-Host3 # Dummy folder for non-clustered host (created by system)
   | | +-Host3
   | | | +-VM3
   | | | +-VM4
   | | |
   +-vm # VM folder (created by system)
   | +-VM1
   | +-VM2
   | +-Folder2 # VM and Template folder
   | | +-hadoop1
   | | +-NestedFolder1
   | | | +-VM3
   | | | +-VM4
```
*/

type Resources struct {
	Dcs      Dcs
	Folders  Folders
	Clusters Clusters
	Hosts    Hosts
	VMs      VMs
}
