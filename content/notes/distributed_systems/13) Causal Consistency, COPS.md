---
title: Causal Consistency
date: 23rd October 2023
tags: [subject/distributed-system, distributed-storage, consistency]
---
# Notes
Requirements for consistencies
### Linearizability 
1) Global ordering
2) Completion to Invocation globally

### Sequential 
1) Global ordering
2) Completion to Invocation per client

### Causal 
1) Completion to invocation per client
2) Write and then later read dependency
3) Transitivity for 1) [[#Causal]]  and 2) [[#Causal]]

### Fork-join Causal

### Eventual