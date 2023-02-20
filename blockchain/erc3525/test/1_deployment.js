const { loadFixture } = require("@nomicfoundation/hardhat-network-helpers");
const { expect } = require("chai");

describe("Deployment", function () {
  async function CCSFixture() {
    const [owner] = await ethers.getSigners();

    const CCS = await ethers.getContractFactory("CCS");
    const ccs = await CCS.deploy();

    return { ccs, owner };
  }

  describe("Deployment", function () {
    it("Should set the correct owner", async function () {
      const { ccs, owner } = await loadFixture(CCSFixture);

      expect(await ccs.owner()).to.equal(owner.address);
    });
  });
});
