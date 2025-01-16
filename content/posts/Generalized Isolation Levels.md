---
title: Generalized Isolation Level Definitions
tags: [subject/distributed-system, database, isolation, paper-summary]
date: 03 Sep 2024
description: Personal notes on the Generalized Isolation Level Definitions paper that was published in 2000.
bsky: "https://bsky.app/profile/akshatsharma.xyz/post/3lbuwd3hldk2u"
---

Personal notes on [Generalized Isolation Level Definitions](https://pmg.csail.mit.edu/papers/icde00.pdf)

## What are isolation levels? 
SQL Isolation levels(part of ACID) - 
1) Introduced to increase performance
2) By improving concurrency
## Standardized? 
- ANSI/ISO SQL isolation level definitions were created to define and standardize different isolation levels.
- Goal was to make them implementation independent.

## Critique
- A paper by the name of "A critique of ANSI/ISO SQL Isolation levels" was published that talked systematically about why the ANSI/ISO SQL isolation level definitions were ambiguous by bringing up anomalies they didn't account for and vague wordings in the definitions that could cause problems. 
- The paper came proposed different definitions that solved the ambiguity problem but they explicitly used locking mechanisms to achieve those definitions.
- Techniques like Optimistic Concurrency Control and Multi-Version mechanisms couldn't work with the definitions provided in the "Critique" paper. 

## Do we really care about OCC and MVCC? 
- Short Answer - YES.
- OCC can outperform PCC in low-contention geo-distributed workloads/datastore settings which is precisely how a lot of databases are being set up now-a-days.
- CockroachDB for the longest time was running on OCC.

## How do the original definitions work? 
- Define a bad behavior as phenomena
- More restrictive consistency levels disallow more phenomena and serializabilty does not permit any phenomena.
- The "Critique" paper demonstrated problems in the phenomena definitions and missing phenomenas. Similar to original isolation level definitions, they then defined isolation levels using the new and improved phenomenas definitions.

## Phenomenas/Anomalies
Refer to the paper summary for [Critique of ANSI SQL Isolation Levels]({{< ref "18) Critique of ANSI SQL Isolation Levels" >}}) which goes over ANSI anomalies and remedies. Alternatively, you can refer to the [paper](https://dl.acm.org/doi/10.1145/223784.223785) directly as well. 


## Problem with original definitions
- The real problem with the preventative approach is that the phenomena are expressed in terms of single-object histories. 
- However, the properties of interest are often multi-object constraints. To avoid problems with such constraints, the phenomena need to restrict what can be done with individual objects more than is necessary.

## Database model and Transaction Histories
### Database model
Better to refer to the paper for this section as summarizing it is pointless, almost everything mentioned is important to understand the crux of the paper.
### Transaction Histories
Better to refer to the paper for this section as summarizing it is pointless, almost everything mentioned is important to understand the crux of the paper.
### Predicates
Q) What are predicates? 
- A predicate P identifies a Boolean condition (e.g., as in the WHERE clause of a SQL statement) and the relations on which the condition has to be applied; one or more relations can be specified in P. All tuples that match this condition are read or modified depending on whether a predicate-based read or write is being considered.

**Definition 1 : Version set of a predicate-based operation.**
When a transaction executes a read or write based on a predicate P, the system selects a version for each tuple in P’s relations. The set of selected versions is called the
Version set of this predicate-based operation and is denoted by Vset(P).

#### Predicate-based reads
- Any reads that get served by the predicates will show up as a separate events in the transaction history. Ex - ri (P: Vset(P)) ri (xj ) ri (yk )
- Any count based predicate won't make separate read events show up in the history as you are not actually accessing any object/version.
#### Predicate-based modifications
- A modification based on a predicate P is modeled as a predicate-based read followed by write operations on tuples that match P.

### Conflicts and Serialization Graphs
Better to refer to the paper for this section as summarizing it is pointless, almost everything mentioned is important to understand the crux of the paper.

## New Isolation Level Definitions
### Isolation Level PL-1
**Phenomena -** 
**G0: Write Cycles** - A history H exhibits phenomenon G0 if DSG(H) contains a directed cycle consisting entirely of write-dependency edges.

**PL-1 disallows G0.**
- Non-serializable interleaving of write operations is possible among uncommitted transactions as long as such interleavings are disallowed among committed transactions (e.g., by aborting some transactions)
- Note that predicate-writes might slip through this since there's no write-dependency, its a read and then write making it a read-write edge.

### Isolation Level PL-2
Only disallowing G0 doesn't handle any kind of read serialization, reads technically can happen on uncommitted or even aborted transactions. To avoid this problem, more phenomenas are defined and then circumvented.

Phenomenas - G1
- **G1a: Aborted Reads.** A history H shows phenomenon G1a if it contains an aborted transaction T1 and a committed transaction T2 such that T2 has read some object (maybe via a predicate) modified by T1 .
- **G1b: Intermediate Reads.** A history H shows phenomenon G1b if it contains a committed transaction T2 that has read a version of object x (maybe via a predicate) written by transaction T1 that was not T1 ’s final modification of x. 
- **G1c: Circular Information Flow.** A history H exhibits phenomenon G1c if DSG(H) contains a directed cycle consisting entirely of dependency edges.
**G1c implies G0. How?**
**PL-2 disallows G1.**

### Isolation Level PL-3
In a system that proscribes only G1, it is possible for a transaction to read inconsistent data and therefore to make inconsistent updates.

**G2: Anti-dependency Cycles.** A history H exhibits phenomenon G2 if DSG(H) contains a directed cycle with one or more anti-dependency edges.

**PL-3 disallows G1 and G2.**

### Isolation Level PL-2.99(Repeatable Read)
**G2-item: Item Anti-dependency Cycles.** A history H exhibits phenomenon G2-item if DSG(H) contains a directed cycle having one or more item-anti-dependency edges.

**PL-2.99 disallows G1 and G2-item.**

> {{< figure src="/media/images/generalized_isolation_level_definitions.png" align="center" >}}


It turns out, the phenomenas defined in this paper are the ones that are generally used while defining isolation levels but instead of using the names of the isolation levels(ex - PL1, PL2), databases use the older names to define what level of isolation is supported due to which there's a certain ambuigity that pops up while talking about isolation in popular databases.

A complete list of what phenomenas are avoided at what isolation level in popular databases was compiled by Martin Klepmann(of DDIA fame) in [this GitHub repository](https://github.com/ept/hermitage).


# References: 
1) Generalized Isolation Level Definitions https://pmg.csail.mit.edu/papers/icde00.pdf
2) Hermitage: Testing transaction isolation levels https://github.com/ept/hermitage