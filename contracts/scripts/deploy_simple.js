const hre = require("hardhat");

async function main() {
  const Simple = await hre.ethers.getContractFactory("Simple");
  const subscriptionId = 1;
  const simple = await Simple.deploy(subscriptionId);
  await simple.waitForDeployment();

  console.log("Simple deployed to:", simple.target);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });