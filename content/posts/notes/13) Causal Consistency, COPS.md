---
title: Causal Consistency
date: 23 Oct 2023
tags: [subject/distributed-system, distributed-storage, consistency]
---

# Notes

Requirements for consistencies

### Linearizability

1. Global ordering
2. Completion to Invocation globally

### Sequential

1. Global ordering
2. Completion to Invocation per client

### Causal

1. Completion to invocation per client
2. Write and then later read dependency
3. Transitivity for 1) [#Causal]({{< ref "#causal" >}}) and 2) [#Causal]({{< ref "#causal" >}})

### Fork-join Causal

### Eventual
