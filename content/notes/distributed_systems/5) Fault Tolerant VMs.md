---
title: Fault Tolerant VMs
date: 13th September 2023
tags: [subject/distributed-system]
---
Personal Notes - [[Fault Tolerant Systems]]
# Notes
While reading a paper, also pay attention to the authors and why they encountered this problem.

While replicating, there are 2 options for replication:
1) Sending the **input** to the replicas
2) Sending the **changes** to the replicas
More often than not, sending the input is more **performant** and **easier** to do.

VMware hypervisor is a binary translation hypervisor(Another hypervisor type is TRAP)

Incase of this paper, not every single state(of the VM) is replicated, but only some key states are recorded and replayed(as the performance tradeoff is too much).

Thread scheduling is deterministic for single core systems(mostly). 
For multi-core systems thread scheduling is non-deterministic(because threads can be scheduled to run simultaneously on two different cores).

# TO-DO:
Search up hypervisor types and levels.
