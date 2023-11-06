---
title: PAXOS
date: 20 Sep 2023
tags: [subject/distributed-system, consensus]
---
# RAFT Notes

### What kind of system is RAFT? Exactly once/At least once/ At most once? 
Two ways to approach this: 
1) Application layer is responsible for de-duplication and the command is just appended (albeit with the same ID).
2) Consensus layer(RAFT) is responsible for checking duplicates.

> [!NOTE]
> Let the client handle ID generation instead of servers.
# PAXOS Notes
Majority concept is for making sure that two different groups of systems cannot reach different consensus as two majorities will always overlap.

> [!IMPORTANT] How is a value chosen(committed)
> Majority need to accept the value->Each server accepts the first value it sees->Obvious problem that each server accepts its own value->Each server can accept multiple values->Obvious problem that multiple majorities->Strict rule that ensures that if there are multiple majorities then there values must be the same->If one value is chosen and other value is accepted then it means that they are same values->If one value is chosen and other value is proposed(proposal accepted) then that means they are same values->If one value is possibly chosen and other is proposed then make sure either the first one is never chosen or it is again proposed.



# Additional Reading
[Viewstamped Replication](https://pmg.csail.mit.edu/papers/vr-revisited.pdf)