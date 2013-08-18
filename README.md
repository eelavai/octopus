Octopus
=======

A network simulator built on UNIX sockets.


Usage
-----

    octopus [--seed=42] node1sockname node1dir [node2sockname node2dir...]

Each node's directory is the namespace under which they listen on the
appropriately named socket and expect to find other servers. For instance, if
`nodeNsockname` is `nodeN`, and `nodeNdir` is `/tmp/dirN`, the first node should
listen on the socket `/tmp/dir1/node1` and would attempt to connect to `node3`
through the UNIX socket `/tmp/dir1/node3`. Octopus will listen on
`/tmp/dir1/node3` and establish a corresponding connection to `/tmp/dir3/node3`
on your behalf.

One of the characteristic features of distributed systems is their
nondeterminism, and Octopus cannot produce completely deterministic network
simulations. That being said, Octopus will make an effort to provide similar
results when run twice with the same seed value and other parameters.


Design
------

Octopus attempts to simulate a variable-latency and lossy fully-connected
network topology between N hosts. It will never violate TCP-like behavior: all
bytes that arrive at the destination are guaranteed to arrive in the order they
were sent, without any data corruption along the way. That being said, Octopus
is free to split up the byte stream in whatever way it desires, to arbitrarily
delay the stream, or to close the stream at any point (discarding any traffic
that was "in-flight" at the time).

Octopus has two primary design goals:
1. To simulate both connection-level and network-level events. A simple network
   simulator might simulate each network link individually, with every
   connection failing independently. However, this poorly models real network
   events, where failures are highly correlated. For instance, a network split
   will cause some of the nodes to be able to communicate amongst themselves,
   but not with the nodes on the other side of the split.
2. To provide roughly-reproducible results. In particular, the system should be
   architected in such a way that running it twice with the same seed on agents
   with roughly the same communication pattern will exercise similar code paths
   and trigger similar application-level bugs.

To satisfy these goals, Octopus has two primary object types: a single network
director, and several point-to-point connections.

We model the network as a completely connected graph between N nodes, where each
edge represents a connection. Each connection has an associated delay and queue
size, which represents the latency of the connection between the nodes and the
(rough) number of bytes that are allowed to be in-flight at any point in time
respectively. Furthermore, each connection has a flag that represents whether it
is currently connected. The network as a whole, then, can be described as a
state machine over the states of each of its N(N-1)/2 connections.

The network director, then, is simply a process which randomly selects a
sequence of state transitions and times at which they occur. These network
events can be point mutations (changing the latency of a single link, for
instance) or bulk operations (disrupting the network along some split). And
since the director's state transitions represent the bulk of the nondeterminism
in the system, seeding its random number generator is sufficient to produce
roughly reproducible network traces.

One side effect of modeling node-to-node connection state as opposed to the
state of individual transport-level connections is that if node A makes a lot of
connections to node B, they will all exhibit similar latency, and all fail at
roughly the same time. It's unclear if this it at all realistic, but the
reproducibility benefits we get from node-to-node state outweighs the cost of
the unrealism in my mind.
