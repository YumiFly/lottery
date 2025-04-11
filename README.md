# lottery
solc --abi --bin --base-path . --include-path ./node_modules/ ./lottery.sol -o ./build

abigen --abi ./build/LotteryManager.abi --bin ./build/LotteryManager.bin --pkg lottery --type LotteryManager --out ./lottery.go



