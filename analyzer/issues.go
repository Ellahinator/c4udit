package analyzer

// AllIssues returns the list of all issues.
func AllIssues() []Issue {
	return append(append(GasOpIssues(), LowRiskIssues()...), NonCriticalIssues()...)
}

// GasOpIssues returns the list of all gas optimization issues.
func GasOpIssues() []Issue {
	return []Issue{
		// G-01 - Don't Initialize Variables with Default Value
		{
			"G-01",
			GASOP,
			"Cache Array Length Outside of Loop",
			"Reading array length at each iteration of the loop takes 6 gas (3 for mload and 3 to place memory_offset) in the stack. Caching the array length in the stack saves around 3 gas per iteration.",
			// `(uint[0-9]*[[:blank:]][a-z,A-Z,0-9]*.?=.?0;)|(bool.[a-z,A-Z,0-9]*.?=.?false;)|(int[0-9]*[[:blank:]][a-z,A-Z,0-9]*.?=.?0;)`,
			`(for.*\.length)`,
			"Store the array’s length in a variable before the for-loop.",
		},
		// G-02 - Cache Array Length Outside of Loop
		{
			"G-02",
			GASOP,
			"Use `!= 0` instead of `> 0` for Unsigned Integer Comparison in require statements",
			"`!= 0` is cheapear than `> 0` when comparing unsigned integers in require statements.",
			`(require.*>0|require.*> 0)`,
			"Use `!= 0` instead of `> 0`.",
		},
		// G-03 - Use != 0 instead of > 0 for Unsigned Integer Comparison
		{
			"G-03",
			GASOP,
			"Reduce the size of error messages (Long revert Strings).",
			"Shortening revert strings to fit in 32 bytes will decrease deployment time gas and will decrease runtime gas when the revert condition is met. Revert strings that are longer than 32 bytes require at least one additional mstore, along with additional overhead for computing memory offset, etc.",
			"require.*\".{33,}\"|require.*'.{33,}'",
			"Shorten the revert strings to fit in 32 bytes, or use custom errors if >0.8.4.",
		},
		// G-04 - Use Custom Errors instead of Revert Strings.
		{
			"G-04",
			GASOP,
			"Use Custom Errors instead of Revert Strings.",
			"Custom errors from Solidity 0.8.4 are cheaper than revert strings (cheaper deployment cost and runtime cost when the revert condition is met)",
			// `(pragma solidity \^0.[8-9].[0-9]|pragma solidity \>0.[8-9].[0-9]|pragma solidity 0.[8-9].[4-9])`,
			"(pragma solidity \\^0.[8-9].[0-9]|pragma solidity >0.[8-9].[0-9]|pragma solidity 0.[8-9].[4-9])?(require.*\"|require.*\\')",
			"Use custom errors instead of revert strings.",
		},

		//G-05
		{
			"G-05",
			GASOP,
			"No need to initialize variables with default values",
			"If a variable is not set/initialized, it is assumed to have the default value (0, false, 0x0 etc depending on the data type). Explicitly initializing it with its default value is an anti-pattern and wastes gas.",
			`(uint[0-9]*[[:blank:]][a-z,A-Z,0-9]*.?=.?0;)|(bool.[a-z,A-Z,0-9]*.?=.?false;)|(int[0-9]*[[:blank:]][a-z,A-Z,0-9]*.?=.?0;)`,
			"Remove explicit default initializations.",
		},
		// G-06 - ++i costs less gas compared to i++ or i += 1
		{
			"G-06",
			GASOP,
			"`++i` costs less gas compared to `i++` or `i += 1`",
			"`++i` costs less gas compared to `i++` or `i += 1` for unsigned integer, as pre-increment is cheaper (about 5 gas per iteration). This statement is true even with the optimizer enabled.",
			`(i\++|i \+= 1|i\--|[a-z,A-Z]*\++\)|[a-z,A-Z]*\++[[:blank:]]\)|[a-z,A-Z]*\--|i \-= 1)`,
			"Use `++i` instead of `i++` to increment the value of an uint variable. Same thing for `--i` and `i--`.",
		},

		// G-07 - Use Shift Right/Left instead of Division/Multiplication if possible
		{
			"G-07",
			GASOP,
			"Use Shift Right/Left instead of Division/Multiplication if possible",
			"A division/multiplication by any number `x` being a power of 2 can be calculated by shifting `log2(x)` to the right/left. While the `DIV` opcode uses 5 gas, the `SHR` opcode only uses 3 gas. Furthermore, Solidity's division operation also includes a division-by-0 prevention which is bypassed using shifting.",
			`(/[2,4,8]|/ [2,4,8]|\*[2,4,8]|\* [2,4,8])`,
			"Use SHR/SHL.\nBad\n```solidity\nuint256 b = a / 2;\nuint256 c = a / 4;\nuint256 d = a * 8;\n```\nGood\n```solidity\nuint256 b = a >> 1;\nuint256 c = a >> 2;\nuint256 d = a << 3;\n```",
		},
		// G-08 - Contracts using unlocked pragma.
		{
			"G-08",
			GASOP,
			"Contracts using unlocked pragma.",
			"Contracts in scope use `pragma solidity ^0.X.Y` or `pragma solidity >0.X.Y`, allowing wide enough range of versions.",
			`pragma solidity \^|pragma solidity >`,
			"Consider locking compiler version, for example `pragma solidity 0.8.6`. This can have additional benefits, for example using custom errors to save gas and so forth.",
		},
		// G-09 - Empty blocks should be removed or emit something
		{
			"G-09",
			GASOP,
			"Empty blocks should be removed or emit something",
			"Empty blocks should be removed or emit something. Waste of gas.",
			`(function.*{*})`,
			"The code should be refactored such that they no longer exist, or the block should do something useful, such as emitting an event or reverting.",
		},
		// G-10 - Use `calldata` instead of `memory` for read-only arguments in `external` functions.
		{
			"G-10",
			GASOP,
			"Use `calldata` instead of `memory` for read-only arguments in `external` functions.",
			"When a function with a `memory` array is called externally, the `abi.decode()` step has to use a for-loop to copy each index of the `calldata` to the `memory` index. Each iteration of this for-loop costs at least 60 gas (i.e. 60 * <mem_array>.length). Using calldata directly, obliviates the need for such a loop in the contract code and runtime execution.",
			`(function.*memory.*external)`,
			"Use `calldata` instead of `memory`.",
		},
		// G-11 - Use `storage` instead of `memory` for structs/arrays.
		{
			"G-11",
			GASOP,
			"Use `storage` instead of `memory` for structs/arrays.",
			"When fetching data from a `storage` location, assigning the data to a `memory` variable causes all fields of the struct/array to be read from `storage`, which incurs a Gcoldsload (2100 gas) for each field of the struct/array. If the fields are read from the new `memory` variable, they incur an additional MLOAD rather than a cheap stack read. Instead of declearing the variable with the `memory` keyword, declaring the variable with the `storage` keyword and caching any fields that need to be re-read in stack variables, will be much cheaper, only incuring the Gcoldsload for the fields actually read. The only time it makes sense to read the whole struct/array into a `memory` variable, is if the full struct/array is being returned by the function, is being passed to a function that requires `memory`, or if the array/struct is being read from another `memory` array/struct.",
			`memory.*\=.*\[.*\]`,
			"Use `storage` instead of `memory` for findings above",
		},
		// G-12 - `x += y` costs more gas than `x = x + y` for state variables.
		{
			"G-12",
			GASOP,
			"`x += y` costs more gas than `x = x + y` for state variables.",
			"Same thing applies for subtraction",
			`.*\+=|.*\-=`,
			"Use `x = x + y` instead of `x += y",
		},
		// G-13 - Don't use `SafeMath` if solidity version  >0.8.0.
		// {
		// 	"G-13",
		// 	GASOP,
		// 	"Don't use `SafeMath` if solidity version  >0.8.0.",
		// 	"Version 0.8.0 introduces internal overflow/underflow checks, so using SafeMath is redundant and adds overhead.",
		// 	`SafeMath`,
		// 	"Remove `SafeMath`.",
		// },
		
	}
}

// LowRiskIssues returns the list of all low risk issues.
func LowRiskIssues() []Issue {
	return []Issue{
		// L-01 - Unsafe ERC20 Operation(s)
		{
			"L-01",
			LOW,
			"Unsafe ERC20 Operation(s)",
			"The return value of an external `transfer`/`transferFrom` call is not checked",
			`\.transfer\(|\.transferFrom\(|\.approve\(`, // ".tranfer(", ".transferFrom(" or ".approve("
			"Use `SafeERC20`, or ensure that the `transfer`/`transferFrom` return value is checked.",
		},
		// L-02 - Unspecific Compiler Version Pragma
		{
			"L-02",
			LOW,
			"Unspecific Compiler Version Pragma",
			"A known vulnerable compiler version may accidentally be selected or security tools might fall-back to an older compiler version ending up checking a different EVM compilation that is ultimately deployed on the blockchain.",
			"pragma solidity (\\^|>)", // "pragma solidity ^" or "pragma solidity >"
			"Avoid floating pragmas for non-library contracts. It is recommended to pin to a concrete compiler version.",
		},
		// L-03 - Do not use Deprecated Library Functions
		{
			"L-03",
			LOW,
			"Do not use Deprecated Library Functions",
			"The usage of deprecated library functions should be discouraged.",
			`_setupRole\(|safeApprove\(|latestAnswer`, // _setupRole and safeApprove are common deprecated lib functions
			"Use `safeIncreaseAllowance` / `safeDecreaseAllowance` instead of `safeApprove`.",
		},
		// L-04 - Open TODOs
		{
			"L-04",
			LOW,
			"Open TODOs",
			"There are many open TODOs throughout the various test files, but also some among the code files.",
			`TODO`,
			"Remove TODO's before deployment",
		},
		// L-05 - ecrecover()
		{
			"L-05",
			LOW,
			"`ecrecover()` not checked for signer address of zero",
			"The `ecrecover()` function returns an address of zero when the signature does not match. This can cause problems if address zero is ever the owner of assets, and someone uses the permit function on address zero. If that happens, any invalid signature will pass the checks, and the assets will be stealable. ",
			`(address*[[:blank:]][a-z,A-Z,0-9]*.?=.?ecrecover.*;)`,
			"Add a check to ensure `ecrecover()` does not return an address of zero.",
		},
		// L-06 - `_safeMint()` should be used rather than `_mint()` wherever possible.
		{
			"L-06",
			LOW,
			"`_safeMint()` should be used rather than `_mint()` wherever possible.",
			"`_mint()` is [discouraged](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/d4d8d2ed9798cc3383912a23b5e8d5cb602f7d4b/contracts/token/ERC721/ERC721.sol#L271) in favor of `_safeMint()` which ensures that the recipient is either an EOA or implements `IERC721Receiver`.",
			`\_mint\(.*\)`,
			"Use either [OpenZeppelin's](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/d4d8d2ed9798cc3383912a23b5e8d5cb602f7d4b/contracts/token/ERC721/ERC721.sol#L238-L250) or [solmate's](https://github.com/transmissions11/solmate/blob/4eaf6b68202e36f67cab379768ac6be304c8ebde/src/tokens/ERC721.sol#L180) version of this function.",
		},
		// L-07 - Expressions for constant values such as a call to `keccak256()`, should use `immutable` rather than `constant`.
		{
			"L-07",
			LOW,
			"Expressions for constant values such as a call to `keccak256()`, should use `immutable` rather than `constant`.",
			"",
			`.*constant.*\=.*keccak256\(.*\)`,
			"",
		},
	}
}

//non critical
func NonCriticalIssues() []Issue {
	return []Issue{
		{
			"N-01",
			NC,
			"Use of `ecrecover()` is susceptible to signature malleability",
			"", // Impact should be empty.
			`ecrecover`,
			"Use OpenZeppelin's `ECDSA` contract rather than calling `ecrecover()` directly.",
		},
		{
			"N-02",
			NC,
			"Declare `uint` as `uint256`",
			"",
			` uint | int `,
			"To favor explicitness, all instances of `uint`/`int` should be declared as `uint256`/`int256`.",
		},
	}
}
