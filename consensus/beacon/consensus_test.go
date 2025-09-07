package beacon

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/holiman/uint256"
)

func createTestStateDB() (*state.StateDB, error) {
	// Create an in-memory database
	db := rawdb.NewMemoryDatabase()
	
	// Create state database
	trieDB := triedb.NewDatabase(db, triedb.HashDefaults)
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(trieDB, nil))
	return stateDB, err
}

func TestIsCarbonSystemWithdrawal(t *testing.T) {
	// Test carbon system withdrawal detection
	carbonWithdrawal := &types.Withdrawal{
		Index:     0,
		Validator: CarbonSystemValidatorIndex, // 0xFFFFFFFF
		Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Amount:    5096, // Gwei
	}

	// Test normal validator withdrawal
	normalWithdrawal := &types.Withdrawal{
		Index:     1,
		Validator: 12345, // Normal validator index
		Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Amount:    32000000000, // 32 ETH in Gwei
	}

	// Test carbon withdrawal detection
	if !isCarbonSystemWithdrawal(carbonWithdrawal) {
		t.Error("Carbon system withdrawal not detected correctly")
	}

	// Test normal withdrawal detection  
	if isCarbonSystemWithdrawal(normalWithdrawal) {
		t.Error("Normal withdrawal incorrectly detected as carbon withdrawal")
	}
}

func TestProcessCarbonWithdrawal(t *testing.T) {
	stateDB, err := createTestStateDB()
	if err != nil {
		t.Fatalf("Failed to create test state DB: %v", err)
	}

	treasuryAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	
	// Create a carbon withdrawal
	carbonWithdrawal := &types.Withdrawal{
		Index:     0,
		Validator: CarbonSystemValidatorIndex, // 0xFFFFFFFF
		Address:   treasuryAddr,
		Amount:    5096, // Gwei
	}

	// Get initial balance
	initialBalance := stateDB.GetBalance(treasuryAddr)

	// Process carbon withdrawal
	err = processCarbonWithdrawal(carbonWithdrawal, stateDB)
	if err != nil {
		t.Fatalf("Failed to process carbon withdrawal: %v", err)
	}

	// Check balance increase
	finalBalance := stateDB.GetBalance(treasuryAddr)
	expectedIncrease := new(uint256.Int).SetUint64(carbonWithdrawal.Amount)
	expectedIncrease = expectedIncrease.Mul(expectedIncrease, uint256.NewInt(params.GWei))
	
	expectedFinalBalance := new(uint256.Int).Add(initialBalance, expectedIncrease)

	if finalBalance.Cmp(expectedFinalBalance) != 0 {
		t.Errorf("Balance not updated correctly. Expected: %v, Got: %v", expectedFinalBalance, finalBalance)
	}
}

func TestProcessValidatorWithdrawal(t *testing.T) {
	stateDB, err := createTestStateDB()
	if err != nil {
		t.Fatalf("Failed to create test state DB: %v", err)
	}

	validatorAddr := common.HexToAddress("0x9876543210987654321098765432109876543210")
	
	// Create a normal validator withdrawal
	validatorWithdrawal := &types.Withdrawal{
		Index:     1,
		Validator: 12345, // Normal validator index
		Address:   validatorAddr,
		Amount:    32000000000, // 32 ETH in Gwei
	}

	// Get initial balance
	initialBalance := stateDB.GetBalance(validatorAddr)

	// Process validator withdrawal
	err = processValidatorWithdrawal(validatorWithdrawal, stateDB)
	if err != nil {
		t.Fatalf("Failed to process validator withdrawal: %v", err)
	}

	// Check balance increase
	finalBalance := stateDB.GetBalance(validatorAddr)
	expectedIncrease := new(uint256.Int).SetUint64(validatorWithdrawal.Amount)
	expectedIncrease = expectedIncrease.Mul(expectedIncrease, uint256.NewInt(params.GWei))
	
	expectedFinalBalance := new(uint256.Int).Add(initialBalance, expectedIncrease)

	if finalBalance.Cmp(expectedFinalBalance) != 0 {
		t.Errorf("Balance not updated correctly. Expected: %v, Got: %v", expectedFinalBalance, finalBalance)
	}
}

func TestCarbonWithdrawalIntegration(t *testing.T) {
	stateDB, err := createTestStateDB()
	if err != nil {
		t.Fatalf("Failed to create test state DB: %v", err)
	}

	treasuryAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	validatorAddr := common.HexToAddress("0x9876543210987654321098765432109876543210")

	// Create mixed withdrawals
	withdrawals := []*types.Withdrawal{
		{
			Index:     0,
			Validator: CarbonSystemValidatorIndex, // Carbon withdrawal
			Address:   treasuryAddr,
			Amount:    5096, // Gwei
		},
		{
			Index:     1,
			Validator: 12345, // Normal validator withdrawal
			Address:   validatorAddr,
			Amount:    16000000000, // 16 ETH in Gwei
		},
		{
			Index:     2,
			Validator: CarbonSystemValidatorIndex, // Another carbon withdrawal
			Address:   treasuryAddr,
			Amount:    2048, // Gwei
		},
	}

	// Get initial balances
	initialTreasuryBalance := stateDB.GetBalance(treasuryAddr)
	initialValidatorBalance := stateDB.GetBalance(validatorAddr)

	// Process all withdrawals (simulating the main loop)
	for _, w := range withdrawals {
		if isCarbonSystemWithdrawal(w) {
			err := processCarbonWithdrawal(w, stateDB)
			if err != nil {
				t.Fatalf("Failed to process carbon withdrawal: %v", err)
			}
		} else {
			err := processValidatorWithdrawal(w, stateDB)
			if err != nil {
				t.Fatalf("Failed to process validator withdrawal: %v", err)
			}
		}
	}

	// Check final balances
	finalTreasuryBalance := stateDB.GetBalance(treasuryAddr)
	finalValidatorBalance := stateDB.GetBalance(validatorAddr)

	// Calculate expected treasury balance (5096 + 2048 = 7144 Gwei)
	expectedTreasuryIncrease := new(uint256.Int).SetUint64(5096 + 2048)
	expectedTreasuryIncrease = expectedTreasuryIncrease.Mul(expectedTreasuryIncrease, uint256.NewInt(params.GWei))
	expectedTreasuryBalance := new(uint256.Int).Add(initialTreasuryBalance, expectedTreasuryIncrease)

	// Calculate expected validator balance (16 ETH in Gwei)
	expectedValidatorIncrease := new(uint256.Int).SetUint64(16000000000)
	expectedValidatorIncrease = expectedValidatorIncrease.Mul(expectedValidatorIncrease, uint256.NewInt(params.GWei))
	expectedValidatorBalance := new(uint256.Int).Add(initialValidatorBalance, expectedValidatorIncrease)

	if finalTreasuryBalance.Cmp(expectedTreasuryBalance) != 0 {
		t.Errorf("Treasury balance not updated correctly. Expected: %v, Got: %v", expectedTreasuryBalance, finalTreasuryBalance)
	}

	if finalValidatorBalance.Cmp(expectedValidatorBalance) != 0 {
		t.Errorf("Validator balance not updated correctly. Expected: %v, Got: %v", expectedValidatorBalance, finalValidatorBalance)
	}
}