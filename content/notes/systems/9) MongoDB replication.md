---
title: MongoDB Replication
date: 02 Oct 2023
tags: [subject/distributed-system, distributed-storage]
---
# Notes
## Background
1) Already existing scheme used db query for replication
2) Replication scheme was in such a way that any follower could query any other follower

Due to this when they tried to implement RAFT there was a mismatch between the architecture and algorithm.
So they decided to change the algorithm to implement in existing architecture.

## Problem
1) Querying will not give you the term number meaning you don't know if you're pulling from the leader or not and if the log entries are correct or not.
2) Any other follower can pull from this stale follower.

## Solution
1) Use updateposition RPC to get the term.

## Implementation
1) Reading can be done by batching read requests with other requests.
# After lecture 
1) Look up TiDB
2) What is mmap?
3) Why fsync twice?
4) TLA+ verification
5) FLP Theorem (https://ilyasergey.net/CS6213/week-03-bft.html)