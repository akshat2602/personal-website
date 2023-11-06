---
title: PAXOS Made Live
date: 27th September 2023
tags: [subject/distributed-system, consensus]
---
Personal Notes - [[PAXOS Made Live]]
# Notes
Runs in 2 phases -
1) Phase 1 - get a promise from majority of servers
2) Phase 2 - get a majority of servers to commit

## Two important rules
1) Two numbers should be used to track the system
	1) Highest seen 
	2) Highest accepted
2) Majority will only happen if the values getting selected are happening at the same number
	1) It can't be that one value got accepted at 1,v1 and another got accepted at 3,v1. It won't be a majority.

<mark style="background: #FFF3A3A6;">RAFT's up-to date term and index matching can be mapped to the two numbers rule in paxos.
RAFT's figure 2 last rule where current term is stored while appending(so a majority cannot be reached as term and index both must be matched) can be matched to the majority rule where values are selected that happened at the same number.</mark>


PAXOS can be converted to RAFT:
1) Comparison of numbers happen in receiver instead of proposer
2) Super prepare(prepare all log entries at the same time)
3) Super accept(Accept all the log entries at the same time)