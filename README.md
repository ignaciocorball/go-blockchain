# 🚀 A Go Implementation of Blockchain

A modern blockchain implementation written in Go, featuring Proof of Stake consensus, smart contracts, and a RESTful API. This project demonstrates the core concepts of blockchain technology in a clean, well-documented codebase.

## 🌟 Features

- ⛓️ **Blockchain Core**
  - Immutable block structure
  - Cryptographic block linking
  - Transaction management
  - Proof of Stake consensus

- 📝 **Smart Contracts**
  - Basic smart contract framework
  - Stateful contract execution
  - Contract validation
  - Persistent contract state

- 💾 **Storage**
  - BadgerDB integration
  - ACID transactions
  - High-performance key-value storage
  - Persistent blockchain data

- 🌐 **API Server**
  - RESTful endpoints
  - Transaction creation
  - Block querying
  - Smart contract deployment and execution

## 🛠️ Prerequisites

- Go 1.16 or higher
- BadgerDB
- Echo framework

## 🚀 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/ignaciocorball/go-blockchain.git
   cd go-blockchain
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/transaction` | Create a new transaction |
| GET | `/block/:hash` | Retrieve block information |
| POST | `/contract` | Deploy a new smart contract |
| POST | `/contract/:id/execute` | Execute a deployed contract |

## 🏗️ Project Structure

```
go-blockchain/
├── api/            # API server implementation
├── blockchain/     # Core blockchain logic
├── contracts/      # Smart contract system
├── storage/        # Database layer
└── main.go         # Application entry point
```

## 🔍 Key Components

### Blockchain Core
- Block creation and validation
- Transaction processing
- Proof of Stake consensus
- Cryptographic security

### Smart Contracts
- Contract deployment
- State management
- Execution environment
- Basic validation

### Storage Layer
- BadgerDB integration
- Block persistence
- Transaction storage
- ACID compliance

### API Server
- RESTful interface
- Transaction handling
- Block querying
- Contract management

## 🔐 Security Features

- ECDSA for transaction signing
- Cryptographic block linking
- Proof of Stake consensus
- Secure smart contract execution

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

## 📚 Documentation

- [Blockchain Core](blockchain/README.md)
- [API Documentation](api/README.md)
- [Smart Contracts](contracts/README.md)
- [Storage Layer](storage/README.md)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- Your Name - Initial work

## 🙏 Acknowledgments

- BadgerDB for the storage layer
- Echo framework for the API server
- The Go community for excellent tools and libraries

## 📝 Contact

- Project Link: [https://github.com/ignaciocorball/go-blockchain](https://github.com/ignaciocorball/go-blockchain)
- Email: ignaciocorball@gmail.com

## 🔄 Roadmap

- [ ] Implement full smart contract language
- [ ] Add network layer for P2P communication
- [ ] Implement wallet system
- [ ] Add more consensus mechanisms
- [ ] Improve API documentation
- [ ] Add monitoring and metrics

---

⭐ Star this repository if you find it useful!