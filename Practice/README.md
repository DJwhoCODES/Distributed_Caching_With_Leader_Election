Distributed In-Memory Cache:

1. A fast key-value store
2. Data stored in RAM
3. Running on multiple machines/processes
4. Clients can connect to any server
5. One server coordinates writes

---

Advantages:

1. Speed 🚀
2. Scalability 📈
3. Reduce DB load
4. Learn systems + networking

---

TCP over HTTP

1. This is a low-level, high-performance, long-lived connection system

=> What HTTP gives you (but you don’t need here)

HTTP is great for:

1. Browsers
2. REST APIs
3. Stateless request/response
4. Human-readable traffic

But HTTP comes with:

1. Headers (lots of extra bytes)
2. Parsing overhead
3. Request lifecycle overhead
4. Mostly stateless design

What TCP gives you (and why it’s chosen)

TCP gives you:

1. Raw byte stream
2. Full control over protocol
3. Persistent connections
4. Minimal overhead
5. Maximum speed

---

Binary Encoded over JSON:

1. JSON is larger → more network bandwidth
2. JSON parsing is slower
3. This project wants high-performance cache behavior like Redis/Memcached
4. Binary encoding is how real caches work

5. {"key":"mykey","value":"myvalue","ttl":5}
   5.1. Raw text size = 28 bytes (actually 33 if you include quotes, braces, colons, etc.)
   5.2. Must be parsed as string → converted to Go types (json.Unmarshal)
   5.3. Parsing overhead = CPU + memory

6. Binary version:
   6.1. 25 bytes (less than JSON)
   6.2. No parsing overhead → binary.Read is very fast

7. Bandwidth saved: ~25–30%
   7.1. CPU saved: JSON parsing is much heavier
   7.2. Scalable: For 1000 SETs/sec, difference multiplies

---

conn.Write(cmd.Bytes()):

1. conn.Write(cmd.Bytes())
2. Write() sends raw bytes over the network
3. Server receives bytes, reconstructs struct using binary.Read

---

encoding/binary (binary.Write(buf, binary.LittleEndian, CmdSet)):

1. LittleEndian → how multi-byte numbers (like int32) are stored in memory.
2. TCP only sends raw bytes. TCP doesn’t understand Go structs, maps, slices, etc.
3. Encodes numbers, lengths, status codes compactly
4. Converts numbers (int32, byte) into binary representation.

5. binary is a package:
   5.1. binary.Write(buf, binary.LittleEndian, int32(len(c.Key))):
   1. it writes the binary representation of the int32 value into the buffer, using little-endian byte order.

6. bytes.Buffer (buf := new(bytes.Buffer)):
   6.1. Buffer is a temporary in-memory storage of bytes.
   6.2. You can Write many things (command type, lengths, values) in order.
   6.3. At the end, return a single contiguous byte slice.
   6.4. Buffer lets you build this sequence efficiently in memory before sending.

---

encoding/binary -> binary.Write() and binary.Read():

1. binary.Write takes a Go value (int32, float64, struct, etc.), converts it into raw bytes
   1.1. binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
   1.2. binary.Write(buf, binary.LittleEndian, []byte(c.Key))
   1.3. [00 00 00 05] [48 65 6C 6C 6F]
   length=5 "Hello"

2. binary.Read
   2.1. it reads raw bytes from a Reader
   2.2. it fills a variable with the decoded value

   2.3. var keyLen int32; binary.Read(conn, binary.LittleEndian, &keyLen)
   2.4. data := make([]byte, keyLen); io.ReadFull(conn, data)

   2.5. Step-by-step:
   2.5.1. binary.Read reads 4 bytes from conn
   2.5.2. decodes them as int32
   2.5.3. stores the result in keyLen
   2.5.4. then we allocate a byte slice (data) of that size
   2.5.5. io.ReadFull reads exactly keyLen bytes of content

3. In Go, if you want a function to fill or update a variable, you must pass a pointer.
   3.1. var keyLen int32; binary.Read(conn, binary.LittleEndian, &keyLen)
   3.2. That's why we pass address of the variable

---

Visual Representation:

[CmdSet] [KeyLen] [Key bytes] [ValueLen] [Value bytes] [TTL]
(0x01) (0x05 00 00 00) ('m','y','k','e','y') (0x07 00 00 00) ('m','y','v','a','l','u','e') (0x05 00 00 00)
