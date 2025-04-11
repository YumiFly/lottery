const hre = require("hardhat");

async function main() {
  // 部署 VRFCoordinatorV2 合约
  const VRFCoordinatorV2 = await hre.ethers.getContractFactory("VRFCoordinatorV2");
  const vrfCoordinator = await VRFCoordinatorV2.deploy();
  await vrfCoordinator.waitForDeployment();
  console.log("VRFCoordinatorV2 deployed to:", vrfCoordinator.target);

  // 部署 LOTToken 合约
  const LOTToken = await hre.ethers.getContractFactory("LOTToken");
  const initialSupply = hre.ethers.parseEther("100000000"); // 1,000,000 个代币
  const lotToken = await LOTToken.deploy(initialSupply);
  await lotToken.waitForDeployment();
  console.log("LOTToken deployed to:", lotToken.target);

  // 部署 SimpleRollout 合约const simple = await Simple.attach("0x5FbDB2315678afecb367f032d93F642f64180aa3"); 
  //// 调用合约函数  const result = await simple.s_results("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266");
  const SimpleRollout = await hre.ethers.getContractFactory("SimpleRollout");
  const subscriptionId = 1; // 替换为你的 VRF 订阅 ID
  const simpleRollout = await SimpleRollout.deploy(subscriptionId, vrfCoordinator.target, (await hre.ethers.getSigners())[0].address);
  await simpleRollout.waitForDeployment();
  console.log("SimpleRollout deployed to:", simpleRollout.target);

  // 部署 LotteryManager 合约
  const LotteryManager = await hre.ethers.getContractFactory("LotteryManager");
  const admin = (await hre.ethers.getSigners())[0].address;
  const owner = (await hre.ethers.getSigners())[1].address;
  const rolloutContract = simpleRollout.target;
  const name = "MyLottery";
  const supply = 100;
  const price = hre.ethers.parseEther("1"); // 1 个代币
  const tokenContract = lotToken.target;
  const lotteryManager = await LotteryManager.deploy(admin, owner, rolloutContract, name, supply, price, tokenContract);
  await lotteryManager.waitForDeployment();
  console.log("LotteryManager deployed to:", lotteryManager.target);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });