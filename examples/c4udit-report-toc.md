# c4udit Report

## Files analyzed
- dummy.sol
# Table of Contents
Low
- [1. Unsafe ERC20 Operation(s)](#1-unsafe-erc20-operations)
- [2. Unspecific Compiler Version Pragma](#2-unspecific-compiler-version-pragma)
- [3. Open TODOs](#3-open-todos)
- [4. `ecrecover()` not checked for signer address of zero](#4-ecrecover-not-checked-for-signer-address-of-zero)
- [5. `_safeMint()` should be used rather than `_mint()` wherever possible.](#5-_safemint-should-be-used-rather-than-_mint-wherever-possible)
- [6. Expressions for constant values such as a call to `keccak256()`, should use `immutable` rather than `constant`.](#6-expressions-for-constant-values-such-as-a-call-to-keccak256-should-use-immutable-rather-than-constant)

Non-Critical
- [1. Use of `ecrecover()` is susceptible to signature malleability](#1-use-of-ecrecover-is-susceptible-to-signature-malleability)
- [2. Declare uint as uint256](#2-declare-uint-as-uint256)


## Low Findings

### 1. Unsafe ERC20 Operation(s)
#### Impact
The return value of an external `transfer`/`transferFrom` call is not checked
#### Findings:
```solidity
dummy.sol::47 => token.transferFrom(msg.sender, address(this), 100);
dummy.sol::72 => IERC721(_token).transferFrom(address(this), _to, _tokenId);
```
#### Recommendation
Use `SafeERC20`, or ensure that the `transfer`/`transferFrom` return value is checked.

### 2. Unspecific Compiler Version Pragma
#### Impact
A known vulnerable compiler version may accidentally be selected or security tools might fall-back to an older compiler version ending up checking a different EVM compilation that is ultimately deployed on the blockchain.
#### Findings:
```solidity
dummy.sol::1 => pragma solidity ^0.8.0;
dummy.sol::2 => pragma solidity >0.8.0;
```
#### Recommendation
Avoid floating pragmas for non-library contracts. It is recommended to pin to a concrete compiler version.

### 3. Open TODOs
#### Impact
There are many open TODOs throughout the various test files, but also some among the code files.
#### Findings:
```solidity
dummy.sol::23 => // TODO
```
#### Recommendation
Remove TODO's before deployment

### 4. `ecrecover()` not checked for signer address of zero
#### Impact
The `ecrecover()` function returns an address of zero when the signature does not match. This can cause problems if address zero is ever the owner of assets, and someone uses the permit function on address zero. If that happens, any invalid signature will pass the checks, and the assets will be stealable. 
#### Findings:
```solidity
dummy.sol::28 => address signer = ecrecover(aiwdd);
```
#### Recommendation
Add a check to ensure `ecrecover()` does not return an address of zero.

### 5. `_safeMint()` should be used rather than `_mint()` wherever possible.
#### Impact
`_mint()` is [discouraged](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/d4d8d2ed9798cc3383912a23b5e8d5cb602f7d4b/contracts/token/ERC721/ERC721.sol#L271) in favor of `_safeMint()` which ensures that the recipient is either an EOA or implements `IERC721Receiver`.
#### Findings:
```solidity
dummy.sol::30 => _mint(_curator, 0);
```
#### Recommendation
Use either [OpenZeppelin's](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/d4d8d2ed9798cc3383912a23b5e8d5cb602f7d4b/contracts/token/ERC721/ERC721.sol#L238-L250) or [solmate's](https://github.com/transmissions11/solmate/blob/4eaf6b68202e36f67cab379768ac6be304c8ebde/src/tokens/ERC721.sol#L180) version of this function.

### 6. Expressions for constant values such as a call to `keccak256()`, should use `immutable` rather than `constant`.
#### Findings:
```solidity
dummy.sol::110 => bytes32 public constant FEE_ROLE = keccak256("FEE_ROLE");
dummy.sol::111 => bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
dummy.sol::112 => bytes32 public constant IMPLEMENTER_ROLE = keccak256("IMPLEMENTER_ROLE");
```
#### Recommendation


## Non-Critical Findings

### 1. Use of `ecrecover()` is susceptible to signature malleability
#### Findings:
```solidity
dummy.sol::28 => address signer = ecrecover(aiwdd);
```
#### Recommendation
Use OpenZeppelin's `ECDSA` contract rather than calling `ecrecover()` directly.

### 2. Declare uint as uint256
#### Findings:
```solidity
dummy.sol::9 => for(uint index = 0; something.length; index++) {}
dummy.sol::10 => for(uint index = 0; something.length; index--) {}
dummy.sol::11 => for(uint index = 0; something.length; index++) {}
dummy.sol::12 => for(uint index = 0; something.length; index--) {}
dummy.sol::15 => uint x = y / 2;
dummy.sol::16 => uint z > 0;
dummy.sol::24 => uint x;
dummy.sol::42 => uint array[] = [1, 2, 3];
```
#### Recommendation
To favor explicitness, all instances of uint should be declared as uint256.


# Table of Contents
Gas
- [1. Cache Array Length Outside of Loop](#1-cache-array-length-outside-of-loop)
- [2. Use `!= 0` instead of `> 0` for Unsigned Integer Comparison in require statements](#2-use--0-instead-of--0-for-unsigned-integer-comparison-in-require-statements)
- [3. Reduce the size of error messages (Long revert Strings).](#3-reduce-the-size-of-error-messages-long-revert-strings)
- [4. Use Custom Errors instead of Revert Strings.](#4-use-custom-errors-instead-of-revert-strings)
- [5. No need to initialize variables with default values](#5-no-need-to-initialize-variables-with-default-values)
- [6. `++i` costs less gas compared to `i++` or `i += 1`](#6-i-costs-less-gas-compared-to-i-or-i--1)
- [7. Use Shift Right/Left instead of Division/Multiplication if possible](#7-use-shift-rightleft-instead-of-divisionmultiplication-if-possible)
- [8. Contracts using unlocked pragma.](#8-contracts-using-unlocked-pragma)
- [9. Empty blocks should be removed or emit something](#9-empty-blocks-should-be-removed-or-emit-something)
- [10. Use `calldata` instead of `memory` for read-only arguments in `external` functions.](#10-use-calldata-instead-of-memory-for-read-only-arguments-in-external-functions)
- [11. Use `storage` instead of `memory` for structs/arrays.](#11-use-storage-instead-of-memory-for-structsarrays)
- [12. `x += y` costs more gas than `x = x + y` for state variables.](#12-x--y-costs-more-gas-than-x--x--y-for-state-variables)

## Gas Findings


### 1. Cache Array Length Outside of Loop
#### Impact
Reading array length at each iteration of the loop takes 6 gas (3 for mload and 3 to place memory_offset) in the stack. Caching the array length in the stack saves around 3 gas per iteration.
#### Findings:
```solidity
dummy.sol::9 => for(uint index = 0; something.length; index++) {}
dummy.sol::10 => for(uint index = 0; something.length; index--) {}
dummy.sol::11 => for(uint index = 0; something.length; index++) {}
dummy.sol::12 => for(uint index = 0; something.length; index--) {}
dummy.sol::43 => for (uint256 i = 0; i < array.length; i++) {
dummy.sol::61 => for (uint256 i = 0; i < _tokens.length; i++) {
dummy.sol::89 => for (uint256 i = 0; i < _tokens.length; i++) {
```
#### Recommendation
Store the array’s length in a variable before the for-loop.

### 2. Use `!= 0` instead of `> 0` for Unsigned Integer Comparison in require statements
#### Impact
`!= 0` is cheapear than `> 0` when comparing unsigned integers in require statements.
#### Findings:
```solidity
dummy.sol::17 => require(z > 0);
```
#### Recommendation
Use `!= 0` instead of `> 0`.

### 3. Reduce the size of error messages (Long revert Strings).
#### Impact
Shortening revert strings to fit in 32 bytes will decrease deployment time gas and will decrease runtime gas when the revert condition is met. Revert strings that are longer than 32 bytes require at least one additional mstore, along with additional overhead for computing memory offset, etc.
#### Findings:
```solidity
dummy.sol::21 => require(x = 2, "This message is more than thirty-two characters.");
dummy.sol::22 => require(x = 2, 'This message is more than thirty-two characters.');
```
#### Recommendation
Shorten the revert strings to fit in 32 bytes, or use custom errors if >0.8.4.

### 4. Use Custom Errors instead of Revert Strings.
#### Impact
Custom errors from Solidity 0.8.4 are cheaper than revert strings (cheaper deployment cost and runtime cost when the revert condition is met)
#### Findings:
```solidity
dummy.sol::21 => require(x = 2, "This message is more than thirty-two characters.");
dummy.sol::22 => require(x = 2, 'This message is more than thirty-two characters.');
dummy.sol::60 => require(_isApprovedOrOwner(msg.sender, 0), "withdraw:not allowed");
dummy.sol::71 => require(_isApprovedOrOwner(msg.sender, 0), "withdraw:not allowed");
dummy.sol::80 => require(_isApprovedOrOwner(msg.sender, 0), "withdraw:not allowed");
dummy.sol::88 => require(_isApprovedOrOwner(msg.sender, 0), "withdraw:not allowed");
```
#### Recommendation
Use custom errors instead of revert strings.

### 5. No need to initialize variables with default values
#### Impact
If a variable is not set/initialized, it is assumed to have the default value (0, false, 0x0 etc depending on the data type). Explicitly initializing it with its default value is an anti-pattern and wastes gas.
#### Findings:
```solidity
dummy.sol::9 => for(uint index = 0; something.length; index++) {}
dummy.sol::10 => for(uint index = 0; something.length; index--) {}
dummy.sol::11 => for(uint index = 0; something.length; index++) {}
dummy.sol::12 => for(uint index = 0; something.length; index--) {}
dummy.sol::18 => bool test = false;
dummy.sol::26 => int y = 0;
dummy.sol::27 => int8 y = 0;
dummy.sol::35 => uint256 a = 0;
dummy.sol::43 => for (uint256 i = 0; i < array.length; i++) {
dummy.sol::61 => for (uint256 i = 0; i < _tokens.length; i++) {
dummy.sol::89 => for (uint256 i = 0; i < _tokens.length; i++) {
```
#### Recommendation
Remove explicit default initializations.

### 6. `++i` costs less gas compared to `i++` or `i += 1`
#### Impact
`++i` costs less gas compared to `i++` or `i += 1` for unsigned integer, as pre-increment is cheaper (about 5 gas per iteration). This statement is true even with the optimizer enabled.
#### Findings:
```solidity
dummy.sol::9 => for(uint index = 0; something.length; index++) {}
dummy.sol::10 => for(uint index = 0; something.length; index--) {}
dummy.sol::11 => for(uint index = 0; something.length; index++) {}
dummy.sol::12 => for(uint index = 0; something.length; index--) {}
dummy.sol::14 => for(uint256 i; length; i++ ) {}
dummy.sol::43 => for (uint256 i = 0; i < array.length; i++) {
dummy.sol::61 => for (uint256 i = 0; i < _tokens.length; i++) {
dummy.sol::89 => for (uint256 i = 0; i < _tokens.length; i++) {
```
#### Recommendation
Use `++i` instead of `i++` to increment the value of an uint variable. Same thing for `--i` and `i--`.

### 7. Use Shift Right/Left instead of Division/Multiplication if possible
#### Impact
A division/multiplication by any number `x` being a power of 2 can be calculated by shifting `log2(x)` to the right/left. While the `DIV` opcode uses 5 gas, the `SHR` opcode only uses 3 gas. Furthermore, Solidity's division operation also includes a division-by-0 prevention which is bypassed using shifting.
#### Findings:
```solidity
dummy.sol::15 => uint x = y / 2;
dummy.sol::44 => i = i / 2;
```
#### Recommendation
Use SHR/SHL.
Bad
```solidity
uint256 b = a / 2
uint256 c = a / 4;
uint256 d = a * 8;
```
Good
```solidity
uint256 b = a >> 1;
uint256 c = a >> 2;
uint256 d = a << 3;
```

### 8. Contracts using unlocked pragma.
#### Impact
Contracts in scope use `pragma solidity ^0.X.Y` or `pragma solidity >0.X.Y`, allowing wide enough range of versions.
#### Findings:
```solidity
dummy.sol::1 => pragma solidity ^0.8.0;
dummy.sol::2 => pragma solidity >0.8.0;
```
#### Recommendation
Consider locking compiler version, for example `pragma solidity 0.8.6`. This can have additional benefits, for example using custom errors to save gas and so forth.

### 9. Empty blocks should be removed or emit something
#### Impact
Empty blocks should be removed or emit something. Waste of gas.
#### Findings:
```solidity
dummy.sol::52 => function() public {}
dummy.sol::53 => function() private { }
```
#### Recommendation
The code should be refactored such that they no longer exist, or the block should do something useful, such as emitting an event or reverting.

### 10. Use `calldata` instead of `memory` for read-only arguments in `external` functions.
#### Impact
When a function with a `memory` array is called externally, the `abi.decode()` step has to use a for-loop to copy each index of the `calldata` to the `memory` index. Each iteration of this for-loop costs at least 60 gas (i.e. 60 * <mem_array>.length). Using calldata directly, obliviates the need for such a loop in the contract code and runtime execution.
#### Findings:
```solidity
dummy.sol::59 => function withdrawMultipleERC721(address[] memory _tokens, uint256[] memory _tokenId, address _to) external override {
```
#### Recommendation
Use `calldata` instead of `memory`.

### 11. Use `storage` instead of `memory` for structs/arrays.
#### Impact
When fetching data from a `storage` location, assigning the data to a `memory` variable causes all fields of the struct/array to be read from `storage`, which incurs a Gcoldsload (2100 gas) for each field of the struct/array. If the fields are read from the new `memory` variable, they incur an additional MLOAD rather than a cheap stack read. Instead of declearing the variable with the `memory` keyword, declaring the variable with the `storage` keyword and caching any fields that need to be re-read in stack variables, will be much cheaper, only incuring the Gcoldsload for the fields actually read. The only time it makes sense to read the whole struct/array into a `memory` variable, is if the full struct/array is being returned by the function, is being passed to a function that requires `memory`, or if the array/struct is being read from another `memory` array/struct.
#### Findings:
```solidity
dummy.sol::100 => TwavObservation memory _twavObservationCurrent = twavObservations[(_index)];
dummy.sol::101 => TwavObservation memory _twavObservationPrev = twavObservations[(_index + 1) % TWAV_BLOCK_NUMBERS];
```
#### Recommendation
Use `storage` instead of `memory` for findings above

### 12. `x += y` costs more gas than `x = x + y` for state variables.
#### Impact
Same thing applies for subtraction
#### Findings:
```solidity
dummy.sol::107 => number += 1;
dummy.sol::108 => number -= 1;
```
#### Recommendation
Use `x = x + y` instead of `x += y

#### Tools used
manual, c4udit, slither

