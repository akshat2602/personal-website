---
title: Spanner
date: 01 Nov 2023
tags: [subject/distributed-system, distributed-storage]
---
# Notes
Two strategies for implementing deadlock prevention used in 2-phase locking
1) Wound wait - Force the lock to release from another transaction
2) Wait die - Wait for lock to be released or die
https://stackoverflow.com/questions/32794142/what-is-the-difference-between-wait-die-and-wound-wait-deadlock-prevention-a
## Distributed Transactions 
Two problems with distributed transactions
1) Write ahead logging needs to happen on each shard
2) There needs to be a flag on each shard indicating that its in the commit phase (because otherwise locks would be wounded)

This is solved using <mark style="background: #FFF3A3A6;">2-phase commit</mark>
1) Prepare for commit (checks if data is actually written and locks were acquired)
2) Commit

2-phase commit is inherently not a fault tolerant protocol and that is the reason it isn't used so much. But, in Spanner they replicated the coordinator making it fault tolerant.
## Spanner
Spanner doesn't replicate locks because of performance issues(imagine replicating locks across 100s of servers which are in different datacenters).
### Leader change issues
To solve issues of leader changes, spanner makes use of checking whether leader has changed(using epoch/term numbers for leaders)
