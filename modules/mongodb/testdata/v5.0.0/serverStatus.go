// SPDX-License-Identifier: GPL-3.0-or-later

package v5_0_0

var ServerStatus = `
{
    "host": "578b738612b0",
    "version": "5.0.0",
    "process": "mongod",
    "pid": 1,
    "uptime": 68,
    "uptimeMillis": 67710,
    "uptimeEstimate": 67,
    "localTime": {},
    "asserts": {
        "regular": 1234567890123456789,
        "warning": 0,
        "msg": 0,
        "user": 22,
        "tripwire": 0,
        "rollovers": 0
    },
    "catalogStats": {
        "collections": 0,
        "capped": 0,
        "timeseries": 0,
        "views": 0,
        "internalCollections": 3,
        "internalViews": 0
    },
    "connections": {
        "current": 3,
        "available": 838857,
        "totalCreated": 6,
        "active": 2,
        "threaded": 3,
        "exhaustIsMaster": 1,
        "exhaustHello": 0,
        "awaitingTopologyChanges": 1
    },
    "electionMetrics": {
        "stepUpCmd": {
            "called": 0,
            "successful": 0
        },
        "priorityTakeover": {
            "called": 0,
            "successful": 0
        },
        "catchUpTakeover": {
            "called": 0,
            "successful": 0
        },
        "electionTimeout": {
            "called": 0,
            "successful": 0
        },
        "freezeTimeout": {
            "called": 0,
            "successful": 0
        },
        "numStepDownsCausedByHigherTerm": 0,
        "numCatchUps": 0,
        "numCatchUpsSucceeded": 0,
        "numCatchUpsAlreadyCaughtUp": 0,
        "numCatchUpsSkipped": 0,
        "numCatchUpsTimedOut": 0,
        "numCatchUpsFailedWithError": 0,
        "numCatchUpsFailedWithNewTerm": 0,
        "numCatchUpsFailedWithReplSetAbortPrimaryCatchUpCmd": 0,
        "averageCatchUpOps": 0
    },
    "extra_info": {
        "note": "fields vary by platform",
        "user_time_us": 897283,
        "system_time_us": 268586,
        "maximum_resident_set_kb": 108096,
        "input_blocks": 27048,
        "output_blocks": 760,
        "page_reclaims": 17040,
        "page_faults": 211,
        "voluntary_context_switches": 3281,
        "involuntary_context_switches": 423
    },
    "flowControl": {
        "enabled": true,
        "targetRateLimit": 1000000000,
        "timeAcquiringMicros": 171,
        "locksPerKiloOp": 0,
        "sustainerRate": 0,
        "isLagged": false,
        "isLaggedCount": 0,
        "isLaggedTimeMicros": 0
    },
    "freeMonitoring": {
        "state": "undecided"
    },
    "globalLock": {
        "totalTime": 67708000,
        "currentQueue": {
            "total": 0,
            "readers": 0,
            "writers": 0
        },
        "activeClients": {
            "total": 0,
            "readers": 0,
            "writers": 0
        }
    },
    "locks": {
        "ParallelBatchWriterMode": {
            "acquireCount": {
                "r": 45
            }
        },
        "ReplicationStateTransition": {
            "acquireCount": {
                "w": 315
            }
        },
        "Global": {
            "acquireCount": {
                "r": 324,
                "w": 9,
                "W": 5
            }
        },
        "Database": {
            "acquireCount": {
                "r": 13,
                "w": 8,
                "W": 1
            }
        },
        "Collection": {
            "acquireCount": {
                "r": 13,
                "w": 6,
                "W": 2
            }
        },
        "Mutex": {
            "acquireCount": {
                "r": 47
            }
        }
    },
    "logicalSessionRecordCache": {
        "activeSessionsCount": 2,
        "sessionsCollectionJobCount": 1,
        "lastSessionsCollectionJobDurationMillis": 29,
        "lastSessionsCollectionJobTimestamp": {},
        "lastSessionsCollectionJobEntriesRefreshed": 0,
        "lastSessionsCollectionJobEntriesEnded": 0,
        "lastSessionsCollectionJobCursorsClosed": 0,
        "transactionReaperJobCount": 1,
        "lastTransactionReaperJobDurationMillis": 0,
        "lastTransactionReaperJobTimestamp": {},
        "lastTransactionReaperJobEntriesCleanedUp": 0,
        "sessionCatalogSize": 0
    },
    "network": {
        "bytesIn": 7532,
        "bytesOut": 269548,
        "physicalBytesIn": 6827,
        "physicalBytesOut": 269548,
        "numSlowDNSOperations": 0,
        "numSlowSSLOperations": 0,
        "numRequests": 58,
        "tcpFastOpen": {
            "kernelSetting": 1,
            "serverSupported": true,
            "clientSupported": true,
            "accepted": 0
        },
        "compression": {
            "snappy": {
                "compressor": {
                    "bytesIn": 0,
                    "bytesOut": 0
                },
                "decompressor": {
                    "bytesIn": 0,
                    "bytesOut": 0
                }
            },
            "zstd": {
                "compressor": {
                    "bytesIn": 0,
                    "bytesOut": 0
                },
                "decompressor": {
                    "bytesIn": 0,
                    "bytesOut": 0
                }
            },
            "zlib": {
                "compressor": {
                    "bytesIn": 0,
                    "bytesOut": 0
                },
                "decompressor": {
                    "bytesIn": 0,
                    "bytesOut": 0
                }
            }
        },
        "serviceExecutors": {
            "passthrough": {
                "threadsRunning": 3,
                "clientsInTotal": 3,
                "clientsRunning": 3,
                "clientsWaitingForData": 0
            },
            "fixed": {
                "threadsRunning": 1,
                "clientsInTotal": 0,
                "clientsRunning": 0,
                "clientsWaitingForData": 0
            }
        }
    },
    "opLatencies": {
        "reads": {
            "latency": 3336,
            "ops": 6
        },
        "writes": {
            "latency": 0,
            "ops": 0
        },
        "commands": {
            "latency": 12272,
            "ops": 50
        },
        "transactions": {
            "latency": 0,
            "ops": 0
        }
    },
    "opcounters": {
        "insert": 0,
        "query": 6,
        "update": 0,
        "delete": 0,
        "getmore": 0,
        "command": 55
    },
    "opcountersRepl": {
        "insert": 0,
        "query": 0,
        "update": 0,
        "delete": 0,
        "getmore": 0,
        "command": 0
    },
    "readConcernCounters": {
        "nonTransactionOps": {
            "none": 6,
            "noneInfo": {
                "CWRC": {
                    "local": 0,
                    "available": 0,
                    "majority": 0
                },
                "implicitDefault": {
                    "local": 6,
                    "available": 0
                }
            },
            "local": 0,
            "available": 0,
            "majority": 0,
            "snapshot": {
                "withClusterTime": 0,
                "withoutClusterTime": 0
            },
            "linearizable": 0
        },
        "transactionOps": {
            "none": 0,
            "noneInfo": {
                "CWRC": {
                    "local": 0,
                    "majority": 0
                },
                "implicitDefault": {
                    "local": 0
                }
            },
            "local": 0,
            "majority": 0,
            "snapshot": {
                "withClusterTime": 0,
                "withoutClusterTime": 0
            }
        }
    },
    "security": {
        "authentication": {
            "saslSupportedMechsReceived": 0,
            "mechanisms": {
                "MONGODB-X509": {
                    "speculativeAuthenticate": {
                        "received": 0,
                        "successful": 0
                    },
                    "clusterAuthenticate": {
                        "received": 0,
                        "successful": 0
                    },
                    "authenticate": {
                        "received": 0,
                        "successful": 0
                    }
                },
                "SCRAM-SHA-1": {
                    "speculativeAuthenticate": {
                        "received": 0,
                        "successful": 0
                    },
                    "clusterAuthenticate": {
                        "received": 0,
                        "successful": 0
                    },
                    "authenticate": {
                        "received": 0,
                        "successful": 0
                    }
                },
                "SCRAM-SHA-256": {
                    "speculativeAuthenticate": {
                        "received": 0,
                        "successful": 0
                    },
                    "clusterAuthenticate": {
                        "received": 0,
                        "successful": 0
                    },
                    "authenticate": {
                        "received": 0,
                        "successful": 0
                    }
                }
            }
        }
    },
    "storageEngine": {
        "name": "wiredTiger",
        "supportsCommittedReads": true,
        "oldestRequiredTimestampForCrashRecovery": {},
        "supportsPendingDrops": true,
        "dropPendingIdents": 0,
        "supportsSnapshotReadConcern": true,
        "readOnly": false,
        "persistent": true,
        "backupCursorOpen": false,
        "supportsResumableIndexBuilds": true
    },
    "tcmalloc": {
        "generic": {
            "current_allocated_bytes": 91076600,
            "heap_size": 93360128
        },
        "tcmalloc": {
            "pageheap_free_bytes": 892928,
            "pageheap_unmapped_bytes": 0,
            "max_total_thread_cache_bytes": 260046848,
            "current_total_thread_cache_bytes": 979248,
            "total_free_bytes": 1390600,
            "central_cache_free_bytes": 265944,
            "transfer_cache_free_bytes": 145408,
            "thread_cache_free_bytes": 979248,
            "aggressive_memory_decommit": 0,
            "pageheap_committed_bytes": 93360128,
            "pageheap_scavenge_count": 0,
            "pageheap_commit_count": 55,
            "pageheap_total_commit_bytes": 93360128,
            "pageheap_decommit_count": 0,
            "pageheap_total_decommit_bytes": 0,
            "pageheap_reserve_count": 55,
            "pageheap_total_reserve_bytes": 93360128,
            "spinlock_total_delay_ns": 0,
            "release_rate": 1,
            "formattedString": "------------------------------------------------\nMALLOC:       91077176 (   86.9 MiB) Bytes in use by application\nMALLOC: +       892928 (    0.9 MiB) Bytes in page heap freelist\nMALLOC: +       265944 (    0.3 MiB) Bytes in central cache freelist\nMALLOC: +       145408 (    0.1 MiB) Bytes in transfer cache freelist\nMALLOC: +       978672 (    0.9 MiB) Bytes in thread cache freelists\nMALLOC: +      2752512 (    2.6 MiB) Bytes in malloc metadata\nMALLOC:   ------------\nMALLOC: =     96112640 (   91.7 MiB) Actual memory used (physical + swap)\nMALLOC: +            0 (    0.0 MiB) Bytes released to OS (aka unmapped)\nMALLOC:   ------------\nMALLOC: =     96112640 (   91.7 MiB) Virtual address space used\nMALLOC:\nMALLOC:            813              Spans in use\nMALLOC:             33              Thread heaps in use\nMALLOC:           4096              Tcmalloc page size\n------------------------------------------------\nCall ReleaseFreeMemory() to release freelist memory to the OS (via madvise()).\nBytes released to the OS take up virtual address space but no physical memory.\n"
        }
    },
    "tenantMigrations": {
        "currentMigrationsDonating": 0,
        "currentMigrationsReceiving": 0,
        "totalSuccessfulMigrationsDonated": 0,
        "totalSuccessfulMigrationsReceived": 0,
        "totalFailedMigrationsDonated": 0,
        "totalFailedMigrationsReceived": 0
    },
    "trafficRecording": {
        "running": false
    },
    "transactions": {
        "retriedCommandsCount": 0,
        "retriedStatementsCount": 0,
        "transactionsCollectionWriteCount": 0,
        "currentActive": 0,
        "currentInactive": 0,
        "currentOpen": 0,
        "totalAborted": 0,
        "totalCommitted": 0,
        "totalStarted": 0,
        "totalPrepared": 0,
        "totalPreparedThenCommitted": 0,
        "totalPreparedThenAborted": 0,
        "currentPrepared": 0
    },
    "transportSecurity": {
        "1.0": 0,
        "1.1": 0,
        "1.2": 0,
        "1.3": 0,
        "unknown": 0
    },
    "twoPhaseCommitCoordinator": {
        "totalCreated": 0,
        "totalStartedTwoPhaseCommit": 0,
        "totalAbortedTwoPhaseCommit": 0,
        "totalCommittedTwoPhaseCommit": 0,
        "currentInSteps": {
            "writingParticipantList": 0,
            "waitingForVotes": 0,
            "writingDecision": 0,
            "waitingForDecisionAcks": 0,
            "deletingCoordinatorDoc": 0
        }
    },
    "wiredTiger": {
        "uri": "statistics:",
        "block-manager": {
            "blocks pre-loaded": 0,
            "blocks read": 0,
            "blocks written": 24,
            "bytes read": 0,
            "bytes read via memory map API": 0,
            "bytes read via system call API": 0,
            "bytes written": 114688,
            "bytes written for checkpoint": 114688,
            "bytes written via memory map API": 0,
            "bytes written via system call API": 0,
            "mapped blocks read": 0,
            "mapped bytes read": 0,
            "number of times the file was remapped because it changed size via fallocate or truncate": 0,
            "number of times the region was remapped via write": 0
        },
        "cache": {
            "application threads page read from disk to cache count": 0,
            "application threads page read from disk to cache time (usecs)": 0,
            "application threads page write from cache to disk count": 12,
            "application threads page write from cache to disk time (usecs)": 487,
            "bytes allocated for updates": 46146,
            "bytes belonging to page images in the cache": 0,
            "bytes belonging to the history store table in the cache": 182,
            "bytes not belonging to page images in the cache": 50676,
            "cache overflow score": 0,
            "eviction calls to get a page": 0,
            "eviction calls to get a page found queue empty": 0,
            "eviction calls to get a page found queue empty after locking": 0,
            "eviction currently operating in aggressive mode": 0,
            "eviction empty score": 0,
            "eviction passes of a file": 0,
            "eviction server candidate queue empty when topping up": 0,
            "eviction server candidate queue not empty when topping up": 0,
            "eviction server evicting pages": 0,
            "eviction server slept, because we did not make progress with eviction": 0,
            "eviction server unable to reach eviction goal": 0,
            "eviction server waiting for a leaf page": 0,
            "eviction state": 64,
            "eviction walk target strategy both clean and dirty pages": 0,
            "eviction walk target strategy only clean pages": 0,
            "eviction walk target strategy only dirty pages": 0,
            "eviction worker thread active": 4,
            "eviction worker thread created": 0,
            "eviction worker thread evicting pages": 0,
            "eviction worker thread removed": 0,
            "eviction worker thread stable number": 0,
            "files with active eviction walks": 0,
            "files with new eviction walks started": 0,
            "force re-tuning of eviction workers once in a while": 0,
            "forced eviction - history store pages failed to evict while session has history store cursor open": 0,
            "forced eviction - history store pages selected while session has history store cursor open": 0,
            "forced eviction - history store pages successfully evicted while session has history store cursor open": 0,
            "forced eviction - pages evicted that were clean count": 0,
            "forced eviction - pages evicted that were clean time (usecs)": 0,
            "forced eviction - pages evicted that were dirty count": 0,
            "forced eviction - pages evicted that were dirty time (usecs)": 0,
            "forced eviction - pages selected because of too many deleted items count": 0,
            "forced eviction - pages selected count": 0,
            "forced eviction - pages selected unable to be evicted count": 0,
            "forced eviction - pages selected unable to be evicted time": 0,
            "hazard pointer check calls": 0,
            "hazard pointer check entries walked": 0,
            "hazard pointer maximum array length": 0,
            "history store score": 0,
            "history store table max on-disk size": 0,
            "history store table on-disk size": 0,
            "internal pages queued for eviction": 0,
            "internal pages seen by eviction walk": 0,
            "internal pages seen by eviction walk that are already queued": 0,
            "maximum bytes configured": 505413632,
            "maximum page size at eviction": 0,
            "modified pages evicted by application threads": 0,
            "operations timed out waiting for space in cache": 0,
            "pages currently held in the cache": 19,
            "pages evicted by application threads": 0,
            "pages evicted in parallel with checkpoint": 0,
            "pages queued for eviction": 0,
            "pages queued for eviction post lru sorting": 0,
            "pages queued for urgent eviction": 0,
            "pages queued for urgent eviction during walk": 0,
            "pages queued for urgent eviction from history store due to high dirty content": 0,
            "pages seen by eviction walk that are already queued": 0,
            "pages selected for eviction unable to be evicted": 0,
            "pages selected for eviction unable to be evicted as the parent page has overflow items": 0,
            "pages selected for eviction unable to be evicted because of active children on an internal page": 0,
            "pages selected for eviction unable to be evicted because of failure in reconciliation": 0,
            "pages walked for eviction": 0,
            "percentage overhead": 8,
            "tracked bytes belonging to internal pages in the cache": 3770,
            "tracked bytes belonging to leaf pages in the cache": 46906,
            "tracked dirty pages in the cache": 2,
            "bytes currently in the cache": 50676,
            "bytes dirty in the cache cumulative": 5968,
            "bytes read into cache": 0,
            "bytes written from cache": 22543,
            "checkpoint blocked page eviction": 0,
            "checkpoint of history store file blocked non-history store page eviction": 0,
            "eviction walk target pages histogram - 0-9": 0,
            "eviction walk target pages histogram - 10-31": 0,
            "eviction walk target pages histogram - 128 and higher": 0,
            "eviction walk target pages histogram - 32-63": 0,
            "eviction walk target pages histogram - 64-128": 0,
            "eviction walk target pages reduced due to history store cache pressure": 0,
            "eviction walks abandoned": 0,
            "eviction walks gave up because they restarted their walk twice": 0,
            "eviction walks gave up because they saw too many pages and found no candidates": 0,
            "eviction walks gave up because they saw too many pages and found too few candidates": 0,
            "eviction walks reached end of tree": 0,
            "eviction walks restarted": 0,
            "eviction walks started from root of tree": 0,
            "eviction walks started from saved location in tree": 0,
            "hazard pointer blocked page eviction": 0,
            "history store table insert calls": 0,
            "history store table insert calls that returned restart": 0,
            "history store table out-of-order resolved updates that lose their durable timestamp": 0,
            "history store table out-of-order updates that were fixed up by reinserting with the fixed timestamp": 0,
            "history store table reads": 0,
            "history store table reads missed": 0,
            "history store table reads requiring squashed modifies": 0,
            "history store table truncation by rollback to stable to remove an unstable update": 0,
            "history store table truncation by rollback to stable to remove an update": 0,
            "history store table truncation to remove an update": 0,
            "history store table truncation to remove range of updates due to key being removed from the data page during reconciliation": 0,
            "history store table truncation to remove range of updates due to out-of-order timestamp update on data page": 0,
            "history store table writes requiring squashed modifies": 0,
            "in-memory page passed criteria to be split": 0,
            "in-memory page splits": 0,
            "internal pages evicted": 0,
            "internal pages split during eviction": 0,
            "leaf pages split during eviction": 0,
            "modified pages evicted": 0,
            "overflow pages read into cache": 0,
            "page split during eviction deepened the tree": 0,
            "page written requiring history store records": 0,
            "pages read into cache": 0,
            "pages read into cache after truncate": 8,
            "pages read into cache after truncate in prepare state": 0,
            "pages requested from the cache": 251,
            "pages seen by eviction walk": 0,
            "pages written from cache": 12,
            "pages written requiring in-memory restoration": 0,
            "tracked dirty bytes in the cache": 1449,
            "unmodified pages evicted": 0
        },
        "capacity": {
            "background fsync file handles considered": 0,
            "background fsync file handles synced": 0,
            "background fsync time (msecs)": 0,
            "bytes read": 0,
            "bytes written for checkpoint": 22230,
            "bytes written for eviction": 0,
            "bytes written for log": 24576,
            "bytes written total": 46806,
            "threshold to call fsync": 0,
            "time waiting due to total capacity (usecs)": 0,
            "time waiting during checkpoint (usecs)": 0,
            "time waiting during eviction (usecs)": 0,
            "time waiting during logging (usecs)": 0,
            "time waiting during read (usecs)": 0
        },
        "checkpoint-cleanup": {
            "pages added for eviction": 0,
            "pages removed": 0,
            "pages skipped during tree walk": 0,
            "pages visited": 6
        },
        "connection": {
            "auto adjusting condition resets": 19,
            "auto adjusting condition wait calls": 447,
            "auto adjusting condition wait raced to update timeout and skipped updating": 0,
            "detected system time went backwards": 0,
            "files currently open": 14,
            "hash bucket array size for data handles": 512,
            "hash bucket array size general": 512,
            "memory allocations": 6225,
            "memory frees": 5239,
            "memory re-allocations": 426,
            "pthread mutex condition wait calls": 1119,
            "pthread mutex shared lock read-lock calls": 1417,
            "pthread mutex shared lock write-lock calls": 128,
            "total fsync I/Os": 48,
            "total read I/Os": 17,
            "total write I/Os": 58
        },
        "cursor": {
            "cached cursor count": 16,
            "cursor bulk loaded cursor insert calls": 0,
            "cursor close calls that result in cache": 603,
            "cursor create calls": 62,
            "cursor insert calls": 59,
            "cursor insert key and value bytes": 35428,
            "cursor modify calls": 0,
            "cursor modify key and value bytes affected": 0,
            "cursor modify value bytes modified": 0,
            "cursor next calls": 44,
            "cursor operation restarted": 0,
            "cursor prev calls": 6,
            "cursor remove calls": 0,
            "cursor remove key bytes removed": 0,
            "cursor reserve calls": 0,
            "cursor reset calls": 850,
            "cursor search calls": 172,
            "cursor search history store calls": 0,
            "cursor search near calls": 9,
            "cursor sweep buckets": 84,
            "cursor sweep cursors closed": 0,
            "cursor sweep cursors examined": 0,
            "cursor sweeps": 14,
            "cursor truncate calls": 0,
            "cursor update calls": 0,
            "cursor update key and value bytes": 0,
            "cursor update value size change": 0,
            "cursors reused from cache": 587,
            "Total number of entries skipped by cursor next calls": 0,
            "Total number of entries skipped by cursor prev calls": 0,
            "Total number of entries skipped to position the history store cursor": 0,
            "Total number of times a search near has exited due to prefix config": 0,
            "cursor next calls that skip due to a globally visible history store tombstone": 0,
            "cursor next calls that skip greater than or equal to 100 entries": 0,
            "cursor next calls that skip less than 100 entries": 44,
            "cursor prev calls that skip due to a globally visible history store tombstone": 0,
            "cursor prev calls that skip greater than or equal to 100 entries": 0,
            "cursor prev calls that skip less than 100 entries": 6,
            "open cursor count": 8
        },
        "data-handle": {
            "connection data handle size": 440,
            "connection data handles currently active": 21,
            "connection sweep candidate became referenced": 0,
            "connection sweep dhandles closed": 0,
            "connection sweep dhandles removed from hash list": 0,
            "connection sweep time-of-death sets": 27,
            "connection sweeps": 6,
            "connection sweeps skipped due to checkpoint gathering handles": 0,
            "session dhandles swept": 1,
            "session sweep attempts": 31
        },
        "lock": {
            "checkpoint lock acquisitions": 1,
            "checkpoint lock application thread wait time (usecs)": 0,
            "checkpoint lock internal thread wait time (usecs)": 0,
            "dhandle lock application thread time waiting (usecs)": 0,
            "dhandle lock internal thread time waiting (usecs)": 0,
            "dhandle read lock acquisitions": 291,
            "dhandle write lock acquisitions": 23,
            "durable timestamp queue lock application thread time waiting (usecs)": 0,
            "durable timestamp queue lock internal thread time waiting (usecs)": 0,
            "durable timestamp queue read lock acquisitions": 0,
            "durable timestamp queue write lock acquisitions": 0,
            "metadata lock acquisitions": 1,
            "metadata lock application thread wait time (usecs)": 0,
            "metadata lock internal thread wait time (usecs)": 0,
            "read timestamp queue lock application thread time waiting (usecs)": 0,
            "read timestamp queue lock internal thread time waiting (usecs)": 0,
            "read timestamp queue read lock acquisitions": 0,
            "read timestamp queue write lock acquisitions": 0,
            "schema lock acquisitions": 13,
            "schema lock application thread wait time (usecs)": 0,
            "schema lock internal thread wait time (usecs)": 0,
            "table lock application thread time waiting for the table lock (usecs)": 0,
            "table lock internal thread time waiting for the table lock (usecs)": 0,
            "table read lock acquisitions": 0,
            "table write lock acquisitions": 10,
            "txn global lock application thread time waiting (usecs)": 0,
            "txn global lock internal thread time waiting (usecs)": 0,
            "txn global read lock acquisitions": 14,
            "txn global write lock acquisitions": 6
        },
        "log": {
            "busy returns attempting to switch slots": 0,
            "force archive time sleeping (usecs)": 0,
            "log bytes of payload data": 20121,
            "log bytes written": 24448,
            "log files manually zero-filled": 0,
            "log flush operations": 663,
            "log force write operations": 743,
            "log force write operations skipped": 738,
            "log records compressed": 31,
            "log records not compressed": 2,
            "log records too small to compress": 37,
            "log release advances write LSN": 12,
            "log scan operations": 0,
            "log scan records requiring two reads": 0,
            "log server thread advances write LSN": 5,
            "log server thread write LSN walk skipped": 3071,
            "log sync operations": 15,
            "log sync time duration (usecs)": 11891,
            "log sync_dir operations": 1,
            "log sync_dir time duration (usecs)": 1027,
            "log write operations": 70,
            "logging bytes consolidated": 23936,
            "maximum log file size": 104857600,
            "number of pre-allocated log files to create": 2,
            "pre-allocated log files not ready and missed": 1,
            "pre-allocated log files prepared": 2,
            "pre-allocated log files used": 0,
            "records processed by log scan": 0,
            "slot close lost race": 0,
            "slot close unbuffered waits": 0,
            "slot closures": 17,
            "slot join atomic update races": 0,
            "slot join calls atomic updates raced": 0,
            "slot join calls did not yield": 70,
            "slot join calls found active slot closed": 0,
            "slot join calls slept": 0,
            "slot join calls yielded": 0,
            "slot join found active slot closed": 0,
            "slot joins yield time (usecs)": 0,
            "slot transitions unable to find free slot": 0,
            "slot unbuffered writes": 0,
            "total in-memory size of compressed records": 35650,
            "total log buffer size": 33554432,
            "total size of compressed records": 17206,
            "written slots coalesced": 0,
            "yields waiting for previous log file close": 0
        },
        "perf": {
            "file system read latency histogram (bucket 1) - 10-49ms": 0,
            "file system read latency histogram (bucket 2) - 50-99ms": 0,
            "file system read latency histogram (bucket 3) - 100-249ms": 0,
            "file system read latency histogram (bucket 4) - 250-499ms": 0,
            "file system read latency histogram (bucket 5) - 500-999ms": 0,
            "file system read latency histogram (bucket 6) - 1000ms+": 0,
            "file system write latency histogram (bucket 1) - 10-49ms": 0,
            "file system write latency histogram (bucket 2) - 50-99ms": 0,
            "file system write latency histogram (bucket 3) - 100-249ms": 0,
            "file system write latency histogram (bucket 4) - 250-499ms": 0,
            "file system write latency histogram (bucket 5) - 500-999ms": 0,
            "file system write latency histogram (bucket 6) - 1000ms+": 0,
            "operation read latency histogram (bucket 1) - 100-249us": 0,
            "operation read latency histogram (bucket 2) - 250-499us": 0,
            "operation read latency histogram (bucket 3) - 500-999us": 0,
            "operation read latency histogram (bucket 4) - 1000-9999us": 0,
            "operation read latency histogram (bucket 5) - 10000us+": 0,
            "operation write latency histogram (bucket 1) - 100-249us": 0,
            "operation write latency histogram (bucket 2) - 250-499us": 0,
            "operation write latency histogram (bucket 3) - 500-999us": 0,
            "operation write latency histogram (bucket 4) - 1000-9999us": 0,
            "operation write latency histogram (bucket 5) - 10000us+": 0
        },
        "reconciliation": {
            "internal-page overflow keys": 0,
            "leaf-page overflow keys": 0,
            "maximum seconds spent in a reconciliation call": 0,
            "page reconciliation calls that resulted in values with prepared transaction metadata": 0,
            "page reconciliation calls that resulted in values with timestamps": 0,
            "page reconciliation calls that resulted in values with transaction ids": 1,
            "pages written including at least one prepare state": 0,
            "pages written including at least one start timestamp": 0,
            "records written including a prepare state": 0,
            "split bytes currently awaiting free": 0,
            "split objects currently awaiting free": 0,
            "approximate byte size of timestamps in pages written": 0,
            "approximate byte size of transaction IDs in pages written": 96,
            "fast-path pages deleted": 0,
            "page reconciliation calls": 12,
            "page reconciliation calls for eviction": 0,
            "pages deleted": 0,
            "pages written including an aggregated newest start durable timestamp ": 0,
            "pages written including an aggregated newest stop durable timestamp ": 0,
            "pages written including an aggregated newest stop timestamp ": 0,
            "pages written including an aggregated newest stop transaction ID": 0,
            "pages written including an aggregated newest transaction ID ": 0,
            "pages written including an aggregated oldest start timestamp ": 0,
            "pages written including an aggregated prepare": 0,
            "pages written including at least one start durable timestamp": 0,
            "pages written including at least one start transaction ID": 1,
            "pages written including at least one stop durable timestamp": 0,
            "pages written including at least one stop timestamp": 0,
            "pages written including at least one stop transaction ID": 0,
            "records written including a start durable timestamp": 0,
            "records written including a start timestamp": 0,
            "records written including a start transaction ID": 12,
            "records written including a stop durable timestamp": 0,
            "records written including a stop timestamp": 0,
            "records written including a stop transaction ID": 0
        },
        "session": {
            "flush state races": 0,
            "flush_tier busy retries": 0,
            "flush_tier operation calls": 0,
            "open session count": 15,
            "session query timestamp calls": 0,
            "table alter failed calls": 0,
            "table alter successful calls": 0,
            "table alter unchanged and skipped": 0,
            "table compact failed calls": 0,
            "table compact successful calls": 0,
            "table create failed calls": 0,
            "table create successful calls": 9,
            "table drop failed calls": 0,
            "table drop successful calls": 0,
            "table rename failed calls": 0,
            "table rename successful calls": 0,
            "table salvage failed calls": 0,
            "table salvage successful calls": 0,
            "table truncate failed calls": 0,
            "table truncate successful calls": 0,
            "table verify failed calls": 0,
            "table verify successful calls": 0,
            "tiered operations dequeued and processed": 0,
            "tiered operations scheduled": 0,
            "tiered storage local retention time (secs)": 0,
            "tiered storage object size": 0
        },
        "thread-state": {
            "active filesystem fsync calls": 0,
            "active filesystem read calls": 0,
            "active filesystem write calls": 0
        },
        "thread-yield": {
            "application thread time evicting (usecs)": 0,
            "application thread time waiting for cache (usecs)": 0,
            "connection close blocked waiting for transaction state stabilization": 0,
            "connection close yielded for lsm manager shutdown": 0,
            "data handle lock yielded": 0,
            "get reference for page index and slot time sleeping (usecs)": 0,
            "log server sync yielded for log write": 0,
            "page access yielded due to prepare state change": 0,
            "page acquire busy blocked": 0,
            "page acquire eviction blocked": 0,
            "page acquire locked blocked": 0,
            "page acquire read blocked": 0,
            "page acquire time sleeping (usecs)": 0,
            "page delete rollback time sleeping for state change (usecs)": 0,
            "page reconciliation yielded due to child modification": 0
        },
        "transaction": {
            "Number of prepared updates": 0,
            "Number of prepared updates committed": 0,
            "Number of prepared updates repeated on the same key": 0,
            "Number of prepared updates rolled back": 0,
            "prepared transactions": 0,
            "prepared transactions committed": 0,
            "prepared transactions currently active": 0,
            "prepared transactions rolled back": 0,
            "query timestamp calls": 70,
            "rollback to stable calls": 0,
            "rollback to stable pages visited": 0,
            "rollback to stable tree walk skipping pages": 0,
            "rollback to stable updates aborted": 0,
            "sessions scanned in each walk of concurrent sessions": 2015,
            "set timestamp calls": 0,
            "set timestamp durable calls": 0,
            "set timestamp durable updates": 0,
            "set timestamp oldest calls": 0,
            "set timestamp oldest updates": 0,
            "set timestamp stable calls": 0,
            "set timestamp stable updates": 0,
            "transaction begins": 35,
            "transaction checkpoint currently running": 0,
            "transaction checkpoint currently running for history store file": 0,
            "transaction checkpoint generation": 2,
            "transaction checkpoint history store file duration (usecs)": 181,
            "transaction checkpoint max time (msecs)": 18,
            "transaction checkpoint min time (msecs)": 18,
            "transaction checkpoint most recent duration for gathering all handles (usecs)": 1126,
            "transaction checkpoint most recent duration for gathering applied handles (usecs)": 1022,
            "transaction checkpoint most recent duration for gathering skipped handles (usecs)": 0,
            "transaction checkpoint most recent handles applied": 10,
            "transaction checkpoint most recent handles skipped": 0,
            "transaction checkpoint most recent handles walked": 21,
            "transaction checkpoint most recent time (msecs)": 18,
            "transaction checkpoint prepare currently running": 0,
            "transaction checkpoint prepare max time (msecs)": 1,
            "transaction checkpoint prepare min time (msecs)": 1,
            "transaction checkpoint prepare most recent time (msecs)": 1,
            "transaction checkpoint prepare total time (msecs)": 1,
            "transaction checkpoint scrub dirty target": 0,
            "transaction checkpoint scrub time (msecs)": 0,
            "transaction checkpoint total time (msecs)": 18,
            "transaction checkpoints": 1,
            "transaction checkpoints skipped because database was clean": 0,
            "transaction failures due to history store": 0,
            "transaction fsync calls for checkpoint after allocating the transaction ID": 1,
            "transaction fsync duration for checkpoint after allocating the transaction ID (usecs)": 10447,
            "transaction range of IDs currently pinned": 0,
            "transaction range of IDs currently pinned by a checkpoint": 0,
            "transaction range of timestamps currently pinned": 0,
            "transaction range of timestamps pinned by a checkpoint": 0,
            "transaction range of timestamps pinned by the oldest active read timestamp": 0,
            "transaction range of timestamps pinned by the oldest timestamp": 0,
            "transaction read timestamp of the oldest active reader": 0,
            "transaction rollback to stable currently running": 0,
            "transaction sync calls": 0,
            "transaction walk of concurrent sessions": 136,
            "transactions committed": 7,
            "transactions rolled back": 28,
            "race to read prepared update retry": 0,
            "rollback to stable history store records with stop timestamps older than newer records": 0,
            "rollback to stable inconsistent checkpoint": 0,
            "rollback to stable keys removed": 0,
            "rollback to stable keys restored": 0,
            "rollback to stable restored tombstones from history store": 0,
            "rollback to stable restored updates from history store": 0,
            "rollback to stable sweeping history store keys": 0,
            "rollback to stable updates removed from history store": 0,
            "transaction checkpoints due to obsolete pages": 0,
            "update conflicts": 0
        },
        "concurrentTransactions": {
            "write": {
                "out": 0,
                "available": 128,
                "totalTickets": 128
            },
            "read": {
                "out": 1,
                "available": 127,
                "totalTickets": 128
            }
        },
        "snapshot-window-settings": {
            "total number of SnapshotTooOld errors": 0,
            "minimum target snapshot window size in seconds": 300,
            "current available snapshot window size in seconds": 0,
            "latest majority snapshot timestamp available": "Jan  1 00:00:00:0",
            "oldest majority snapshot timestamp available": "Jan  1 00:00:00:0",
            "pinned timestamp requests": 0,
            "min pinned timestamp": {}
        },
        "oplog": {
            "visibility timestamp": {}
        }
    },
    "mem": {
        "bits": 64,
        "resident": 105,
        "virtual": 1502,
        "supported": true
    },
    "metrics": {
        "apiVersions": {
            "DataGrip": [
                "default"
            ]
        },
        "aggStageCounters": {
            "$_internalConvertBucketIndexStats": 0,
            "$_internalInhibitOptimization": 0,
            "$_internalReshardingIterateTransaction": 0,
            "$_internalSetWindowFields": 0,
            "$_internalSplitPipeline": 0,
            "$_internalUnpackBucket": 0,
            "$_unpackBucket": 0,
            "$addFields": 0,
            "$bucket": 0,
            "$bucketAuto": 0,
            "$changeStream": 0,
            "$collStats": 0,
            "$count": 0,
            "$currentOp": 0,
            "$facet": 0,
            "$geoNear": 0,
            "$graphLookup": 0,
            "$group": 0,
            "$indexStats": 0,
            "$limit": 0,
            "$listLocalSessions": 0,
            "$listSessions": 0,
            "$lookup": 0,
            "$match": 0,
            "$merge": 0,
            "$mergeCursors": 0,
            "$operationMetrics": 0,
            "$out": 0,
            "$planCacheStats": 0,
            "$project": 0,
            "$redact": 0,
            "$replaceRoot": 0,
            "$replaceWith": 0,
            "$sample": 0,
            "$set": 0,
            "$setWindowFields": 0,
            "$skip": 0,
            "$sort": 0,
            "$sortByCount": 0,
            "$unionWith": 0,
            "$unset": 0,
            "$unwind": 0
        },
        "commands": {
            "<UNKNOWN>": 0,
            "_addShard": {
                "failed": 0,
                "total": 0
            },
            "_cloneCollectionOptionsFromPrimaryShard": {
                "failed": 0,
                "total": 0
            },
            "_configsvrAbortReshardCollection": {
                "failed": 0,
                "total": 0
            },
            "_configsvrAddShard": {
                "failed": 0,
                "total": 0
            },
            "_configsvrAddShardToZone": {
                "failed": 0,
                "total": 0
            },
            "_configsvrBalancerCollectionStatus": {
                "failed": 0,
                "total": 0
            },
            "_configsvrBalancerStart": {
                "failed": 0,
                "total": 0
            },
            "_configsvrBalancerStatus": {
                "failed": 0,
                "total": 0
            },
            "_configsvrBalancerStop": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCleanupReshardCollection": {
                "failed": 0,
                "total": 0
            },
            "_configsvrClearJumboFlag": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCommitChunkMerge": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCommitChunkMigration": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCommitChunkSplit": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCommitChunksMerge": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCommitMovePrimary": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCommitReshardCollection": {
                "failed": 0,
                "total": 0
            },
            "_configsvrCreateDatabase": {
                "failed": 0,
                "total": 0
            },
            "_configsvrDropCollection": {
                "failed": 0,
                "total": 0
            },
            "_configsvrDropDatabase": {
                "failed": 0,
                "total": 0
            },
            "_configsvrEnableSharding": {
                "failed": 0,
                "total": 0
            },
            "_configsvrEnsureChunkVersionIsGreaterThan": {
                "failed": 0,
                "total": 0
            },
            "_configsvrMoveChunk": {
                "failed": 0,
                "total": 0
            },
            "_configsvrMovePrimary": {
                "failed": 0,
                "total": 0
            },
            "_configsvrRefineCollectionShardKey": {
                "failed": 0,
                "total": 0
            },
            "_configsvrRemoveShard": {
                "failed": 0,
                "total": 0
            },
            "_configsvrRemoveShardFromZone": {
                "failed": 0,
                "total": 0
            },
            "_configsvrRenameCollectionMetadata": {
                "failed": 0,
                "total": 0
            },
            "_configsvrReshardCollection": {
                "failed": 0,
                "total": 0
            },
            "_configsvrSetAllowMigrations": {
                "failed": 0,
                "total": 0
            },
            "_configsvrShardCollection": {
                "failed": 0,
                "total": 0
            },
            "_configsvrUpdateZoneKeyRange": {
                "failed": 0,
                "total": 0
            },
            "_flushDatabaseCacheUpdates": {
                "failed": 0,
                "total": 0
            },
            "_flushDatabaseCacheUpdatesWithWriteConcern": {
                "failed": 0,
                "total": 0
            },
            "_flushReshardingStateChange": {
                "failed": 0,
                "total": 0
            },
            "_flushRoutingTableCacheUpdates": {
                "failed": 0,
                "total": 0
            },
            "_flushRoutingTableCacheUpdatesWithWriteConcern": {
                "failed": 0,
                "total": 0
            },
            "_getNextSessionMods": {
                "failed": 0,
                "total": 0
            },
            "_getUserCacheGeneration": {
                "failed": 0,
                "total": 0
            },
            "_isSelf": {
                "failed": 0,
                "total": 0
            },
            "_killOperations": {
                "failed": 0,
                "total": 0
            },
            "_mergeAuthzCollections": {
                "failed": 0,
                "total": 0
            },
            "_migrateClone": {
                "failed": 0,
                "total": 0
            },
            "_recvChunkAbort": {
                "failed": 0,
                "total": 0
            },
            "_recvChunkCommit": {
                "failed": 0,
                "total": 0
            },
            "_recvChunkStart": {
                "failed": 0,
                "total": 0
            },
            "_recvChunkStatus": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrAbortReshardCollection": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrCleanupReshardCollection": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrCloneCatalogData": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrCreateCollection": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrCreateCollectionParticipant": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrDropCollection": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrDropCollectionParticipant": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrDropDatabase": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrDropDatabaseParticipant": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrMovePrimary": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrRefineCollectionShardKey": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrRenameCollection": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrRenameCollectionParticipant": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrRenameCollectionParticipantUnblock": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrReshardCollection": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrReshardingOperationTime": {
                "failed": 0,
                "total": 0
            },
            "_shardsvrShardCollection": {
                "failed": 0,
                "total": 0
            },
            "_transferMods": {
                "failed": 0,
                "total": 0
            },
            "abortTransaction": {
                "failed": 0,
                "total": 0
            },
            "aggregate": {
                "failed": 0,
                "total": 0
            },
            "appendOplogNote": {
                "failed": 0,
                "total": 0
            },
            "applyOps": {
                "failed": 0,
                "total": 0
            },
            "authenticate": {
                "failed": 0,
                "total": 0
            },
            "availableQueryOptions": {
                "failed": 0,
                "total": 0
            },
            "buildInfo": {
                "failed": 0,
                "total": 6
            },
            "checkShardingIndex": {
                "failed": 0,
                "total": 0
            },
            "cleanupOrphaned": {
                "failed": 0,
                "total": 0
            },
            "cloneCollectionAsCapped": {
                "failed": 0,
                "total": 0
            },
            "collMod": {
                "failed": 0,
                "total": 0
            },
            "collStats": {
                "failed": 0,
                "total": 0
            },
            "commitTransaction": {
                "failed": 0,
                "total": 0
            },
            "compact": {
                "failed": 0,
                "total": 0
            },
            "connPoolStats": {
                "failed": 0,
                "total": 0
            },
            "connPoolSync": {
                "failed": 0,
                "total": 0
            },
            "connectionStatus": {
                "failed": 0,
                "total": 0
            },
            "convertToCapped": {
                "failed": 0,
                "total": 0
            },
            "coordinateCommitTransaction": {
                "failed": 0,
                "total": 0
            },
            "count": {
                "failed": 0,
                "total": 0
            },
            "create": {
                "failed": 0,
                "total": 0
            },
            "createIndexes": {
                "failed": 0,
                "total": 1
            },
            "createRole": {
                "failed": 0,
                "total": 0
            },
            "createUser": {
                "failed": 0,
                "total": 0
            },
            "currentOp": {
                "failed": 0,
                "total": 0
            },
            "dataSize": {
                "failed": 0,
                "total": 0
            },
            "dbHash": {
                "failed": 0,
                "total": 0
            },
            "dbStats": {
                "failed": 0,
                "total": 0
            },
            "delete": {
                "failed": 0,
                "total": 0
            },
            "distinct": {
                "failed": 0,
                "total": 0
            },
            "donorAbortMigration": {
                "failed": 0,
                "total": 0
            },
            "donorForgetMigration": {
                "failed": 0,
                "total": 0
            },
            "donorStartMigration": {
                "failed": 0,
                "total": 0
            },
            "driverOIDTest": {
                "failed": 0,
                "total": 0
            },
            "drop": {
                "failed": 0,
                "total": 0
            },
            "dropAllRolesFromDatabase": {
                "failed": 0,
                "total": 0
            },
            "dropAllUsersFromDatabase": {
                "failed": 0,
                "total": 0
            },
            "dropConnections": {
                "failed": 0,
                "total": 0
            },
            "dropDatabase": {
                "failed": 0,
                "total": 0
            },
            "dropIndexes": {
                "failed": 0,
                "total": 0
            },
            "dropRole": {
                "failed": 0,
                "total": 0
            },
            "dropUser": {
                "failed": 0,
                "total": 0
            },
            "endSessions": {
                "failed": 0,
                "total": 1
            },
            "explain": {
                "failed": 0,
                "total": 0
            },
            "features": {
                "failed": 0,
                "total": 0
            },
            "filemd5": {
                "failed": 0,
                "total": 0
            },
            "find": {
                "failed": 0,
                "total": 6
            },
            "findAndModify": {
                "arrayFilters": 0,
                "failed": 0,
                "pipeline": 0,
                "total": 0
            },
            "flushRouterConfig": {
                "failed": 0,
                "total": 0
            },
            "fsync": {
                "failed": 0,
                "total": 0
            },
            "fsyncUnlock": {
                "failed": 0,
                "total": 0
            },
            "getCmdLineOpts": {
                "failed": 0,
                "total": 0
            },
            "getDatabaseVersion": {
                "failed": 0,
                "total": 0
            },
            "getDefaultRWConcern": {
                "failed": 0,
                "total": 0
            },
            "getDiagnosticData": {
                "failed": 0,
                "total": 0
            },
            "getFreeMonitoringStatus": {
                "failed": 0,
                "total": 0
            },
            "getLastError": {
                "failed": 0,
                "total": 0
            },
            "getLog": {
                "failed": 0,
                "total": 0
            },
            "getMore": {
                "failed": 0,
                "total": 0
            },
            "getParameter": {
                "failed": 0,
                "total": 0
            },
            "getShardMap": {
                "failed": 0,
                "total": 0
            },
            "getShardVersion": {
                "failed": 0,
                "total": 0
            },
            "getnonce": {
                "failed": 0,
                "total": 0
            },
            "grantPrivilegesToRole": {
                "failed": 0,
                "total": 0
            },
            "grantRolesToRole": {
                "failed": 0,
                "total": 0
            },
            "grantRolesToUser": {
                "failed": 0,
                "total": 0
            },
            "hello": {
                "failed": 0,
                "total": 0
            },
            "hostInfo": {
                "failed": 0,
                "total": 0
            },
            "insert": {
                "failed": 0,
                "total": 0
            },
            "internalRenameIfOptionsAndIndexesMatch": {
                "failed": 0,
                "total": 0
            },
            "invalidateUserCache": {
                "failed": 0,
                "total": 0
            },
            "isMaster": {
                "failed": 1,
                "total": 18
            },
            "killAllSessions": {
                "failed": 0,
                "total": 0
            },
            "killAllSessionsByPattern": {
                "failed": 0,
                "total": 0
            },
            "killCursors": {
                "failed": 0,
                "total": 0
            },
            "killOp": {
                "failed": 0,
                "total": 0
            },
            "killSessions": {
                "failed": 0,
                "total": 0
            },
            "listCollections": {
                "failed": 0,
                "total": 8
            },
            "listCommands": {
                "failed": 0,
                "total": 0
            },
            "listDatabases": {
                "failed": 0,
                "total": 5
            },
            "listIndexes": {
                "failed": 2,
                "total": 8
            },
            "lockInfo": {
                "failed": 0,
                "total": 0
            },
            "logRotate": {
                "failed": 0,
                "total": 0
            },
            "logout": {
                "failed": 0,
                "total": 0
            },
            "mapReduce": {
                "failed": 0,
                "total": 0
            },
            "mergeChunks": {
                "failed": 0,
                "total": 0
            },
            "moveChunk": {
                "failed": 0,
                "total": 0
            },
            "ping": {
                "failed": 0,
                "total": 2
            },
            "planCacheClear": {
                "failed": 0,
                "total": 0
            },
            "planCacheClearFilters": {
                "failed": 0,
                "total": 0
            },
            "planCacheListFilters": {
                "failed": 0,
                "total": 0
            },
            "planCacheSetFilter": {
                "failed": 0,
                "total": 0
            },
            "prepareTransaction": {
                "failed": 0,
                "total": 0
            },
            "profile": {
                "failed": 0,
                "total": 0
            },
            "reIndex": {
                "failed": 0,
                "total": 0
            },
            "recipientForgetMigration": {
                "failed": 0,
                "total": 0
            },
            "recipientSyncData": {
                "failed": 0,
                "total": 0
            },
            "refreshSessions": {
                "failed": 0,
                "total": 0
            },
            "renameCollection": {
                "failed": 0,
                "total": 0
            },
            "repairDatabase": {
                "failed": 0,
                "total": 0
            },
            "replSetAbortPrimaryCatchUp": {
                "failed": 0,
                "total": 0
            },
            "replSetFreeze": {
                "failed": 0,
                "total": 0
            },
            "replSetGetConfig": {
                "failed": 0,
                "total": 0
            },
            "replSetGetRBID": {
                "failed": 0,
                "total": 0
            },
            "replSetGetStatus": {
                "failed": 0,
                "total": 0
            },
            "replSetHeartbeat": {
                "failed": 0,
                "total": 0
            },
            "replSetInitiate": {
                "failed": 0,
                "total": 0
            },
            "replSetMaintenance": {
                "failed": 0,
                "total": 0
            },
            "replSetReconfig": {
                "failed": 0,
                "total": 0
            },
            "replSetRequestVotes": {
                "failed": 0,
                "total": 0
            },
            "replSetResizeOplog": {
                "failed": 0,
                "total": 0
            },
            "replSetStepDown": {
                "failed": 0,
                "total": 0
            },
            "replSetStepDownWithForce": {
                "failed": 0,
                "total": 0
            },
            "replSetStepUp": {
                "failed": 0,
                "total": 0
            },
            "replSetSyncFrom": {
                "failed": 0,
                "total": 0
            },
            "replSetUpdatePosition": {
                "failed": 0,
                "total": 0
            },
            "revokePrivilegesFromRole": {
                "failed": 0,
                "total": 0
            },
            "revokeRolesFromRole": {
                "failed": 0,
                "total": 0
            },
            "revokeRolesFromUser": {
                "failed": 0,
                "total": 0
            },
            "rolesInfo": {
                "failed": 0,
                "total": 0
            },
            "rotateCertificates": {
                "failed": 0,
                "total": 0
            },
            "saslContinue": {
                "failed": 0,
                "total": 0
            },
            "saslStart": {
                "failed": 0,
                "total": 0
            },
            "serverStatus": {
                "failed": 0,
                "total": 6
            },
            "setDefaultRWConcern": {
                "failed": 0,
                "total": 0
            },
            "setFeatureCompatibilityVersion": {
                "failed": 0,
                "total": 0
            },
            "setFreeMonitoring": {
                "failed": 0,
                "total": 0
            },
            "setIndexCommitQuorum": {
                "failed": 0,
                "total": 0
            },
            "setParameter": {
                "failed": 0,
                "total": 0
            },
            "setShardVersion": {
                "failed": 0,
                "total": 0
            },
            "shardingState": {
                "failed": 0,
                "total": 0
            },
            "shutdown": {
                "failed": 0,
                "total": 0
            },
            "splitChunk": {
                "failed": 0,
                "total": 0
            },
            "splitVector": {
                "failed": 0,
                "total": 0
            },
            "startRecordingTraffic": {
                "failed": 0,
                "total": 0
            },
            "startSession": {
                "failed": 0,
                "total": 0
            },
            "stopRecordingTraffic": {
                "failed": 0,
                "total": 0
            },
            "top": {
                "failed": 0,
                "total": 0
            },
            "update": {
                "arrayFilters": 0,
                "failed": 0,
                "pipeline": 0,
                "total": 0
            },
            "updateRole": {
                "failed": 0,
                "total": 0
            },
            "updateUser": {
                "failed": 0,
                "total": 0
            },
            "usersInfo": {
                "failed": 0,
                "total": 0
            },
            "validate": {
                "failed": 0,
                "total": 0
            },
            "validateDBMetadata": {
                "failed": 0,
                "total": 0
            },
            "voteCommitIndexBuild": {
                "failed": 0,
                "total": 0
            },
            "waitForFailPoint": {
                "failed": 0,
                "total": 0
            },
            "whatsmyuri": {
                "failed": 0,
                "total": 0
            }
        },
        "cursor": {
            "moreThanOneBatch": 0,
            "timedOut": 0,
            "totalOpened": 0,
            "lifespan": {
                "greaterThanOrEqual10Minutes": 0,
                "lessThan10Minutes": 0,
                "lessThan15Seconds": 0,
                "lessThan1Minute": 0,
                "lessThan1Second": 0,
                "lessThan30Seconds": 0,
                "lessThan5Seconds": 0
            },
            "open": {
                "noTimeout": 0,
                "pinned": 0,
                "total": 0
            }
        },
        "document": {
            "deleted": 0,
            "inserted": 0,
            "returned": 3,
            "updated": 0
        },
        "dotsAndDollarsFields": {
            "inserts": 0,
            "updates": 0
        },
        "getLastError": {
            "wtime": {
                "num": 0,
                "totalMillis": 0
            },
            "wtimeouts": 0,
            "default": {
                "unsatisfiable": 0,
                "wtimeouts": 0
            }
        },
        "mongos": {
            "cursor": {
                "moreThanOneBatch": 0,
                "totalOpened": 0
            }
        },
        "operation": {
            "scanAndOrder": 0,
            "writeConflicts": 0
        },
        "operatorCounters": {
            "expressions": {
                "$_internalJsEmit": 0,
                "$abs": 0,
                "$acos": 0,
                "$acosh": 0,
                "$add": 0,
                "$allElementsTrue": 0,
                "$and": 0,
                "$anyElementTrue": 0,
                "$arrayElemAt": 0,
                "$arrayToObject": 0,
                "$asin": 0,
                "$asinh": 0,
                "$atan": 0,
                "$atan2": 0,
                "$atanh": 0,
                "$avg": 0,
                "$binarySize": 0,
                "$bsonSize": 0,
                "$ceil": 0,
                "$cmp": 0,
                "$concat": 0,
                "$concatArrays": 0,
                "$cond": 0,
                "$const": 0,
                "$convert": 0,
                "$cos": 0,
                "$cosh": 0,
                "$dateAdd": 0,
                "$dateDiff": 0,
                "$dateFromParts": 0,
                "$dateFromString": 0,
                "$dateSubtract": 0,
                "$dateToParts": 0,
                "$dateToString": 0,
                "$dateTrunc": 0,
                "$dayOfMonth": 0,
                "$dayOfWeek": 0,
                "$dayOfYear": 0,
                "$degreesToRadians": 0,
                "$divide": 0,
                "$eq": 0,
                "$exp": 0,
                "$filter": 0,
                "$first": 0,
                "$floor": 0,
                "$function": 0,
                "$getField": 0,
                "$gt": 0,
                "$gte": 0,
                "$hour": 0,
                "$ifNull": 0,
                "$in": 0,
                "$indexOfArray": 0,
                "$indexOfBytes": 0,
                "$indexOfCP": 0,
                "$isArray": 0,
                "$isNumber": 0,
                "$isoDayOfWeek": 0,
                "$isoWeek": 0,
                "$isoWeekYear": 0,
                "$last": 0,
                "$let": 0,
                "$literal": 0,
                "$ln": 0,
                "$log": 0,
                "$log10": 0,
                "$lt": 0,
                "$lte": 0,
                "$ltrim": 0,
                "$map": 0,
                "$max": 0,
                "$mergeObjects": 0,
                "$meta": 0,
                "$millisecond": 0,
                "$min": 0,
                "$minute": 0,
                "$mod": 0,
                "$month": 0,
                "$multiply": 0,
                "$ne": 0,
                "$not": 0,
                "$objectToArray": 0,
                "$or": 0,
                "$pow": 0,
                "$radiansToDegrees": 0,
                "$rand": 0,
                "$range": 0,
                "$reduce": 0,
                "$regexFind": 0,
                "$regexFindAll": 0,
                "$regexMatch": 0,
                "$replaceAll": 0,
                "$replaceOne": 0,
                "$reverseArray": 0,
                "$round": 0,
                "$rtrim": 0,
                "$second": 0,
                "$setDifference": 0,
                "$setEquals": 0,
                "$setField": 0,
                "$setIntersection": 0,
                "$setIsSubset": 0,
                "$setUnion": 0,
                "$sin": 0,
                "$sinh": 0,
                "$size": 0,
                "$slice": 0,
                "$split": 0,
                "$sqrt": 0,
                "$stdDevPop": 0,
                "$stdDevSamp": 0,
                "$strLenBytes": 0,
                "$strLenCP": 0,
                "$strcasecmp": 0,
                "$substr": 0,
                "$substrBytes": 0,
                "$substrCP": 0,
                "$subtract": 0,
                "$sum": 0,
                "$switch": 0,
                "$tan": 0,
                "$tanh": 0,
                "$toBool": 0,
                "$toDate": 0,
                "$toDecimal": 0,
                "$toDouble": 0,
                "$toHashedIndexKey": 0,
                "$toInt": 0,
                "$toLong": 0,
                "$toLower": 0,
                "$toObjectId": 0,
                "$toString": 0,
                "$toUpper": 0,
                "$trim": 0,
                "$trunc": 0,
                "$type": 0,
                "$unsetField": 0,
                "$week": 0,
                "$year": 0,
                "$zip": 0
            }
        },
        "query": {
            "planCacheTotalSizeEstimateBytes": 0,
            "updateOneOpStyleBroadcastWithExactIDCount": 0
        },
        "queryExecutor": {
            "scanned": 0,
            "scannedObjects": 3,
            "collectionScans": {
                "nonTailable": 6,
                "total": 6
            }
        },
        "record": {
            "moves": 0
        },
        "repl": {
            "executor": {
                "pool": {
                    "inProgressCount": 0
                },
                "queues": {
                    "networkInProgress": 0,
                    "sleepers": 0
                },
                "unsignaledEvents": 0,
                "shuttingDown": false,
                "networkInterface": "DEPRECATED: getDiagnosticString is deprecated in NetworkInterfaceTL"
            },
            "apply": {
                "attemptsToBecomeSecondary": 0,
                "batchSize": 0,
                "batches": {
                    "num": 0,
                    "totalMillis": 0
                },
                "ops": 0
            },
            "buffer": {
                "count": 0,
                "maxSizeBytes": 0,
                "sizeBytes": 0
            },
            "initialSync": {
                "completed": 0,
                "failedAttempts": 0,
                "failures": 0
            },
            "network": {
                "bytes": 0,
                "getmores": {
                    "num": 0,
                    "totalMillis": 0,
                    "numEmptyBatches": 0
                },
                "notPrimaryLegacyUnacknowledgedWrites": 0,
                "notPrimaryUnacknowledgedWrites": 0,
                "oplogGetMoresProcessed": {
                    "num": 0,
                    "totalMillis": 0
                },
                "ops": 0,
                "readersCreated": 0,
                "replSetUpdatePosition": {
                    "num": 0
                }
            },
            "reconfig": {
                "numAutoReconfigsForRemovalOfNewlyAddedFields": 0
            },
            "stateTransition": {
                "lastStateTransition": "",
                "userOperationsKilled": 0,
                "userOperationsRunning": 0
            },
            "syncSource": {
                "numSelections": 0,
                "numSyncSourceChangesDueToSignificantlyCloserNode": 0,
                "numTimesChoseDifferent": 0,
                "numTimesChoseSame": 0,
                "numTimesCouldNotFind": 0
            }
        },
        "ttl": {
            "deletedDocuments": 0,
            "passes": 1
        }
    },
	"repl": "",
    "ok": 1
}`
