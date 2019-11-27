package scaleio

var (
	selectedSdcStatsQuery = `
{
    "selectedStatisticsList":[
        {
            "type":"Sdc",
            "allIds":[
            ],
            "properties":[
                "numOfMappedVolumes",
                "userDataReadBwc",
                "userDataWriteBwc",
                "volumeIds"
            ]
        }
    ]
}
`
)
