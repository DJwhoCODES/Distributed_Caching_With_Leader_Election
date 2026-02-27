Roadmap

PHASE 1 — Networking Basics

- Create a minimal TCP server
- Create a minimal TCP client
- Send/receive plain text over the connection
- Understand buffers, reads, partial reads, packet splits
- Convert this to sending/receiving raw bytes

PHASE 2 — Binary Protocol

- Define our binary message format
- Write encoder (struct → bytes)
- Write decoder (bytes → struct)
- Test sending a SET command end-to-end
- Add error handling, io.ReadFull, etc.

PHASE 3 — Cache

- Create a simple in-memory map
- Add TTL logic
- Connect deserialized commands to cache operations

PHASE 4 — Client Library

- Build client functions → SET / GET / DEL
- Ensure round-trip works

PHASE 5 — Distributed Behavior

- Add leader/follower roles
- Forward writes to leader
- Add simple leader election (static or round-robin)
