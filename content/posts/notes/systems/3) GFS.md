---
title: 3) GFS
date: 06 Sep 2023
tags: [subject/distributed-system, distributed-storage]
---

File systems that already existed during GFS - **GPFS/Lustre/AFS**

### Why GFS was built instead of using existing ones?

1. GPFS and Lustre are paid systems(could be the main reason)
2. Specialized workflows

### Most common solutions to problems that come up:

1. Batching
2. Pipelining
3. Separation of control and data flow

# Consistency

Weakly consistent
Mitigates consistency maintenance to its applications

# Takeaway:

Satisfying the need is the main goal and to have that industries try to create simple system designs.
