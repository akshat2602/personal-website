---
title: 11) Linearizability
date: 16 Oct 2023
tags: [subject/distributed-system, consistency]
---

# Notes

Consistency is a spectrum with weak to strong levels.

GFS/Eventual consistency is an example of weak consistency.

Linearizability is an example of strong consistency.

Consistency can be defined as the following properties:

1. Completion to Invocation(C2I)
    1. Globally
    2. Locally
2. Sequential ordering

If all of these properties are satisfied then the system is linearizable.
If only local C2I and sequential ordering is guaranteed then the system is sequential consistency(second strongest consistency model). Sequential consistency is the one being used in x86 CPUs' memory model.
