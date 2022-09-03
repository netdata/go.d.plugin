// SPDX-License-Identifier: GPL-3.0-or-later

package v5_0_0

const ReplSetGetStatus = `
{
  "date": "2000-01-01T00:00:00.000Z",
  "members": [
    {
      "name": "node1",
      "state": 1,
      "optimeDate": "2000-01-01T00:00:00.000Z"
    },
    {
      "name": "node2",
      "state": 2,
      "optimeDate": "2000-01-01T00:00:00.000Z",
      "lastHeartbeat": "2000-01-01T00:00:00.000Z",
      "lastHeartbeatRecv": "2000-01-01T00:00:00.000Z",
      "pingMs": 0
    },
    {
      "name": "node3",
      "state": 2,
      "optimeDate": "2000-01-01T00:00:00.000Z",
      "lastHeartbeat": "2000-01-01T00:00:00.000Z",
      "lastHeartbeatRecv": "2000-01-01T00:00:00.000Z",
      "pingMs": 0
    }
  ]
}
`

const ReplSetGetStatusNode1 = `
{
  "date": "2000-01-01T00:00:00.000Z",
  "members": [
    {
      "name": "node1",
      "state": 1,
      "optimeDate": "2000-01-01T00:00:00.000Z",
      "lastHeartbeat": "2000-01-01T00:00:00.000Z",
      "lastHeartbeatRecv": "2000-01-01T00:00:00.000Z",
      "pingMs": 0
    }
  ]
}
`

const ReplSetGetStatusNode2 = `
{
  "date": "2000-01-01T00:00:00.000Z",
  "members": [
    {
      "name": "node2",
      "state": 1,
      "optimeDate": "2000-01-01T00:00:00.000Z",
      "lastHeartbeat": "2000-01-01T00:00:00.000Z",
      "lastHeartbeatRecv": "2000-01-01T00:00:00.000Z",
      "pingMs": 0
    }
  ]
}
`
