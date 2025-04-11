const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Simple", function () {
  it("Should roll dice and get a result", async function () {
    const VRFCoordinatorV2 = await ethers.getContractFactory("VRFCoordinatorV2");
    const vrfCoordinatorV2 = await VRFCoordinatorV2.deploy();
    await vrfCoordinatorV2.waitForDeployment();

    const Simple = await ethers.getContractFactory("contracts/randGen/SampleV2.sol:Simple");
    const simple = await Simple.deploy(1, vrfCoordinatorV2.target);
    await simple.waitForDeployment();

    console.log("Simple address:", simple.target);
    const [owner] = await ethers.getSigners();
    console.log("Owner address:", owner.address);

    try {
      const tx = await simple.rollDice(owner.address);
      const receipt = await tx.wait();

      console.log("RollDice Receipt:", receipt);
      console.log("Receipt Events:", receipt.events); // 添加调试信息

      const requestId = receipt.events[0].args.requestId;
      await simple.fulfillRandomWords(requestId, [123]);

      const result = await simple.s_results(owner.address);
      expect(result).to.equal(124); // 123 % 2000 + 1 = 124
    } catch (error) {
      console.error("RollDice Error:", error);
      throw error; // 重新抛出错误以便测试失败
    }
  });

  it("test vrfCoordinatorV2", async function() {
    const VRFCoordinatorV2 = await ethers.getContractFactory("VRFCoordinatorV2");
    const vrfCoordinatorV2 = await VRFCoordinatorV2.deploy();
    await vrfCoordinatorV2.waitForDeployment();

    console.log("vrfCoordinatorV2 address:", vrfCoordinatorV2.target);

    const Simple = await ethers.getContractFactory("contracts/randGen/SampleV2.sol:Simple");
    const simple = await Simple.deploy(1, vrfCoordinatorV2.target);
    await simple.waitForDeployment();

    console.log("Simple address:", simple.target);
    const [owner] = await ethers.getSigners();
    console.log("Owner address:", owner.address);

    try {
      const tx = await simple.rollDice(owner.address);
      const receipt = await tx.wait();

      console.log("RollDice Receipt:", receipt);
      console.log("Receipt Events:", receipt.events); // 添加调试信息

      const requestId = receipt.events[0].args.requestId;

      await vrfCoordinatorV2.CallFullfillRandomWords(requestId, [123]);

      const result = await simple.s_results(owner.address);
      expect(result).to.equal(124);
    } catch (error) {
      console.error("RollDice Error:", error);
      throw error; // 重新抛出错误以便测试失败
    }
  });
});