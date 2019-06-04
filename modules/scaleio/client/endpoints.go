package client

/*
The REST API is served from the VxFlex OS Gateway.
The FxFlex Gateway connects to a single MDM and serves requests by querying the MDM
and reformatting the answers it receives from the MDM in s RESTful manner, back to a REST API.
The Gateway is stateless. It requires the MDM username and password for the login requests.
The login returns a token in the response, that is used for later authentication for other requests.

The token is valid for 8 hours from the time it was created, unless there has been no activity
for 10 minutes, of if the client has sent a logout request.
*/

const (
	// URL PATHs

	pathVersion = "/api/version" // GET
	pathLogin   = "/api/login"   // GET
	pathLogout  = "/api/logout"  // GET
	//pathAllObjects                    = "/api/instances"                         // GET
	//pathObjectsByType                 = "/api/types/%s/instances"                // GET
	//pathObjectByTypeByID              = "/api/types/%s::%s"                      // GET
	//pathObjectRelationshipsByTypeByID = "/api/instances/%s::%s/relationships/%s" // GET
	//pathActionOnObjectsByType         = "/api/types/%s/action/%s"                // POST
	//pathActionOnObjectByTypeByID      = "/api/types/%s::%s/action/%s"            // POST
	pathSelectedStatistics = "/api/instances/querySelectedStatistics" // POST
)

const (
// Types

//typeSystem           = "System"
//typeProtectionDomain = "ProtectionDomain"
//typeSds              = "Sds"
//typeStoragePool      = "StoragePool"
//typeDevice           = "Device"
//typeVolume           = "Volume"
//typeVTree            = "VTree"
//typeSdc              = "Sdc"
//typeUser             = "User"
//typeFaultSet         = "FaultSet"
//typeRfcacheDevice    = "RfcacheDevice"
//typeAlerts           = "Alerts"

// Actions (useful for monitoring)

//actionQuerySelectedStatistics      = "querySelectedStatistics"      // All types except Alarm and User
//actionQuerySystemLimits            = "querySystemLimits"            // System
//actionQueryDisconnectedSdss        = "queryDisconnectedSdss"        // Sds
//actionQuerySdsNetworkLatencyMeters = "querySdsNetworkLatencyMeters" // Sds
//actionQueryFailedDevices           = "queryFailedDevices"           // Device. Note: works strange!

// Relationships

//relStatistics       = "Statistics"         // All types except Alarm and User
//relProtectionDomain = typeProtectionDomain // System
//relSdc              = typeSdc              // System
//relUser             = typeUser             // System
//relStoragePool      = typeStoragePool      // ProtectionDomain
//relFaultSet         = typeFaultSet         // ProtectionDomain
//relSds              = typeSds              // ProtectionDomain
//relRfcacheDevice    = typeRfcacheDevice    // Sds
//relDevice           = typeDevice           // Sds, StoragePool
//relVolue            = typeVolume           // Sdc, StoragePool
//relVTree            = typeVTree            // StoragePool
)
