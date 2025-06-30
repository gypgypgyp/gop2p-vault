# GoP2P-Vault

A Decentralized, Encrypted File Storage System with P2P Networking

```
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
```

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

send text:
```
```
Enter command (text <msg> | upload <file>): text Hello from node 9001
```
You will see:
```
[text msg]: Hello from node 9001
```


upload the file:
```
Enter command (text <msg> | upload <file>): upload hello.txt
```

You will find the uploaded file in ./data/ and see:
```
Enter command (text <msg> | upload <file>): upload hello.txt
2025/06/28 22:16:34 [debug] sending message to 127.0.0.1:9000, size=83
```

```
2025/06/29 15:25:43 [debug] Decoded msg: type=upload, len=30
[recv from 127.0.0.1:53289] Type: upload | Len: 30 bytes
[debug] Received upload message
[upload]: File stored with key 48a5d2f0dff82506c5df4e1199b4c0e21938cfd2
```

There will be an encryped file under ./data directory.


download file:
```
Enter command (text <msg> | upload <file>): download 48a5d2f0dff82506c5df4e1199b4c0e21938cfd2
```
You will find the downloaded file in ./data/ and see:
```
[download]: Sent file 48a5d2f0dff82506c5df4e1199b4c0e21938cfd2 to 127.0.0.1:53289
```
```
[recv from 127.0.0.1:9000] Type: download_result | Len: 30 bytes
[download_result]: File saved to ./data/downloaded_1751236160166123000
```

There will be an decryped file under ./data directory.
