# GoP2P-Vault

A Decentralized, Encrypted File Storage System with P2P Networking

gop2p-vault/
├── cmd/                      # Executable entrypoints
│   └── main.go               # CLI startup (node initialization)
│
├── config/                   # Configuration
│   └── config.go             # Load peers, ports, paths
│
├── p2p/                      # P2P Networking
│   ├── transport.go          # Network transport interface
│   ├── tcp_transport.go      # TCP implementation
│   ├── handshake.go          # Peer handshake protocol
│   ├── peer.go               # Peer state management
│   ├── message.go            # Message structs/types
│   └── encoding.go           # GOB encoding/decoding
│
├── store/                    # Local file storage
│   ├── store.go              # File I/O operations
│   ├── path_hash.go          # Path hashing (SHA-1 + segmentation)
│   └── store_test.go         # Unit tests
│
├── crypto/                   # Encryption
│   ├── crypto.go             # AES-CTR encryption/decryption
│   ├── iv.go                 # IV generation
│   └── crypto_test.go        # Encryption tests
│
├── server/                   # File server logic
│   └── server.go             # Upload/download/sync handlers
│
├── protocol/                 # Message protocols (optional)
│   └── types.go              # Request/response structs
│
├── util/                     # Utilities
│   ├── logger.go             # Logging interface
│   └── helper.go             # Helper functions
│
├── testdata/                 # Test files
│   └── sample.txt
│
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
└── README.md                 # Project documentation


## Quick Start Guide
1. Initialize Go Module
```
go mod init gop2p-vault
```
This generates go.mod with:
```
module gop2p-vault
go 1.24
```

2. Test P2P Networking and File Operations

create a new file
```
echo "This is a file from peer 9001" > hello.txt
```

Start receiver node
```
go run cmd/main.go :9000
```

Start sender node
```
go run cmd/main.go :9001 :9000
```

upload the file
```
upload hello.txt
```

You will see:
```
go run cmd/main.go 127.0.0.1:9001 127.0.0.1:9000
2025/06/28 22:16:20 [P2P] Listening on 127.0.0.1:9001
Enter command (text <msg> | upload <file>): text abc
2025/06/28 22:16:28 [debug] sending message to 127.0.0.1:9000, size=54
Enter command (text <msg> | upload <file>): upload hello.txt
2025/06/28 22:16:34 [debug] sending message to 127.0.0.1:9000, size=83
```

```
go run cmd/main.go 127.0.0.1:9000
2025/06/28 22:16:17 [P2P] Listening on 127.0.0.1:9000
2025/06/28 22:16:28 [debug] Decoded msg: type=text, len=3
[recv from 127.0.0.1:62589] Type: text | Len: 3 bytes
[text msg]: abc
2025/06/28 22:16:34 [debug] Decoded msg: type=upload, len=30
[recv from 127.0.0.1:62589] Type: upload | Len: 30 bytes
[debug] Received upload message
[upload]: File stored with key
```
