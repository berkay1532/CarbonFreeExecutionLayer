# 🌿 Carbon Neutral Ethereum - Execution Layer

## Overview

This repository implements the **execution layer** component of a carbon-neutral Ethereum system. It extends Go Ethereum (Geth) with carbon withdrawal processing capabilities that work in conjunction with the consensus layer to automatically transfer collected carbon offset funds to designated treasury addresses.

## 🎯 What This Does

The execution layer processes special **carbon system withdrawals** that are created by the consensus layer when carbon offset funds are collected from validator rewards. These withdrawals use a special validator index (`0xFFFFFFFF`) to distinguish them from regular validator withdrawals.

## 🔧 Carbon System Architecture

### How It Works

1. **Consensus Layer** collects 1% carbon offset from validator rewards
2. **Consensus Layer** creates carbon withdrawals with special validator index
3. **Execution Layer** (this repo) processes these withdrawals
4. **Treasury Address** receives the carbon offset funds automatically

### Key Components

- **Carbon System Validator Index**: `0xFFFFFFFF` (4294967295) - Special identifier for carbon withdrawals
- **Carbon Withdrawal Processing**: Automatic detection and processing in `consensus/beacon/consensus.go`
- **Treasury Integration**: Direct ETH transfer to configured treasury addresses
- **Metrics & Logging**: Full tracking of carbon fund transfers

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or later
- C compiler (GCC or Clang)

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-carbon-neutrality-eip-repo>
   cd CarbonNeutralityEIP
   ```

2. **Build Geth with Carbon Support**
   ```bash
   make geth
   ```

3. **Verify the build**
   ```bash
   ./build/bin/geth version
   ```

## 🧪 Testing

### Run Carbon System Tests

Test the carbon withdrawal functionality:

```bash
# Test carbon withdrawal integration
go test -v ./consensus/beacon -run TestCarbon

# Test all carbon-related functionality
go test -v ./consensus/beacon -run Test
```

### Expected Test Output

```
=== RUN   TestIsCarbonSystemWithdrawal
--- PASS: TestIsCarbonSystemWithdrawal (0.00s)
=== RUN   TestProcessCarbonWithdrawal
--- PASS: TestProcessCarbonWithdrawal (0.00s)
=== RUN   TestCarbonWithdrawalIntegration
--- PASS: TestCarbonWithdrawal Integration (0.00s)
PASS
```

## 📊 Carbon System Implementation Details

### Carbon Withdrawal Detection

```go
// Special validator index for carbon system withdrawals
const CarbonSystemValidatorIndex = 0xFFFFFFFF

// Check if withdrawal is from carbon system
func isCarbonSystemWithdrawal(withdrawal *types.Withdrawal) bool {
    return withdrawal.Validator == CarbonSystemValidatorIndex
}
```

### Carbon Fund Processing

```go
// Process carbon withdrawal to treasury
func processCarbonWithdrawal(withdrawal *types.Withdrawal, state vm.StateDB) error {
    // Convert amount from gwei to wei
    amount := new(uint256.Int).SetUint64(withdrawal.Amount)
    amount = amount.Mul(amount, uint256.NewInt(params.GWei))
    
    // Add balance to treasury address
    state.AddBalance(withdrawal.Address, amount, tracing.BalanceIncreaseWithdrawal)
    
    // Update metrics and logging
    carbonWithdrawalCount.Inc(1)
    carbonWithdrawalAmount.Inc(int64(withdrawal.Amount))
    
    return nil
}
```

## 🔄 Integration with Consensus Layer

This execution layer works with the **CarbonFreeConsensys** consensus layer:

1. **Consensus Layer Repository**: Contains the beacon chain that:
   - Deducts 1% carbon offset from validator rewards  
   - Creates carbon withdrawals with `ValidatorIndex=0xFFFFFFFF`
   - Sends withdrawals to execution layer

2. **Execution Layer Repository** (this repo): Contains the execution client that:
   - Receives carbon withdrawals from consensus layer
   - Processes them as special system transactions
   - Transfers funds to treasury addresses

## 🌍 Production Deployment

### Configuration

When deploying, ensure you configure:

1. **Treasury Address**: Set the destination for carbon funds
2. **Network Compatibility**: Ensure consensus and execution layers use same network
3. **Monitoring**: Set up metrics collection for carbon fund tracking

### Example Genesis Configuration

```json
{
  "config": {
    "chainId": 1337,
    "carbonEnabled": true,
    "treasuryAddress": "0x1234567890123456789012345678901234567890"
  }
}
```

## 📈 Monitoring Carbon Funds

### Metrics Available

- `carbon/withdrawal/count`: Number of carbon withdrawals processed
- `carbon/withdrawal/amount`: Total amount of carbon funds in Gwei

### Log Output Example

```
INFO Carbon Treasury Transfer
    address=0x1234567890123456789012345678901234567890
    amount_gwei=5096
    amount_wei=5096000000000
    withdrawal_index=0
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Run tests: `go test -v ./consensus/beacon`
4. Commit your changes
5. Push to the branch
6. Create a Pull Request

## 📝 License

This project extends Go Ethereum and maintains the same licensing:
- Library code: GNU Lesser General Public License v3.0
- Binary code: GNU General Public License v3.0

## 🔗 Related Repositories

- **CarbonFreeConsensys**: The consensus layer that generates carbon withdrawals
- **Original Go Ethereum**: https://github.com/ethereum/go-ethereum

## ❓ FAQ

**Q: How much carbon offset is collected?**  
A: The consensus layer deducts 1% from validator rewards, which is then transferred to the treasury.

**Q: Can I change the carbon rate?**  
A: The carbon rate is configured in the consensus layer, not the execution layer.

**Q: How do I verify carbon funds are being collected?**  
A: Check the treasury address balance and monitor the carbon withdrawal logs.

**Q: Is this compatible with mainnet?**  
A: This is a research implementation. Production deployment requires extensive testing and community consensus.

## 🆘 Support

For issues and questions:
- Create an issue in this repository
- Check the test outputs for debugging
- Review the consensus layer repository for the complete carbon system

---

**⚡ Built for ETH Istanbul Hackathon - Making Ethereum Carbon Neutral**