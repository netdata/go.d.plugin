package api

/*
The REST API is served from the VxFlex OS Gateway.
The FxFlex Gateway connects to a single MDM and serves requests by querying the MDM
and reformatting the answers it receives from the MDM in s RESTful manner, back to a REST client.
The Gateway is stateless. It requires the MDM username and password for the login requests.
The login returns a token in the response, that is used for later authentication for other requests.

The token is valid for 8 hours from the time it was created, unless there has been no activity
for 10 minutes, of if the client has sent a logout request.
*/

const (
	// URL PATHs

	PATHVersion                       = "/api/version"                           // GET
	PATHLogin                         = "/api/login"                             // GET
	PATHLogout                        = "/api/logout"                            // GET
	PATHAllObjects                    = "/api/instances"                         // GET
	PATHObjectsByType                 = "/api/types/%s/instances"                // GET
	PATHObjectByTypeByID              = "/api/types/%s::%s"                      // GET
	PATHObjectRelationshipsByTypeByID = "/api/instances/%s::%s/relationships/%s" // GET
	PATHActionOnObjectsByType         = "/api/types/%s/action/%s"                // POST
	PATHActionOnObjectByTypeByID      = "/api/types/%s::%s/action/%s"            // POST
	PATHSelectedStatistics            = "/api/instances/querySelectedStatistics" // POST
)

//type (
//	Type         string
//	Action       string
//	Relationship string
//)

const (
	// Types

	TypeSystem           = "System"
	TypeProtectionDomain = "ProtectionDomain"
	TypeSds              = "Sds"
	TypeStoragePool      = "StoragePool"
	TypeDevice           = "Device"
	TypeVolume           = "Volume"
	TypeVTree            = "VTree"
	TypeSdc              = "Sdc"
	TypeUser             = "User"
	TypeFaultSet         = "FaultSet"
	TypeRfcacheDevice    = "RfcacheDevice"
	TypeAlerts           = "Alerts"

	// Actions (useful for monitoring)

	ActionQuerySelectedStatistics      = "querySelectedStatistics"      // All types except Alarm and User
	ActionQuerySystemLimits            = "querySystemLimits"            // System
	ActionQueryDisconnectedSdss        = "queryDisconnectedSdss"        // Sds
	ActionQuerySdsNetworkLatencyMeters = "querySdsNetworkLatencyMeters" // Sds
	ActionQueryFailedDevices           = "queryFailedDevices"           // Device. Note: works strange!

	// Relationships

	RelStatistics       = "Statistics"         // All types except Alarm and User
	RelProtectionDomain = TypeProtectionDomain // System
	RelSdc              = TypeSdc              // System
	RelUser             = TypeUser             // System
	RelStoragePool      = TypeStoragePool      // ProtectionDomain
	RelFaultSet         = TypeFaultSet         // ProtectionDomain
	RelSds              = TypeSds              // ProtectionDomain
	RelRfcacheDevice    = TypeRfcacheDevice    // Sds
	RelDevice           = TypeDevice           // Sds, StoragePool
	RelVolue            = TypeVolume           // Sdc, StoragePool
	RelVTree            = TypeVTree            // StoragePool
)
