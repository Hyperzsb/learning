const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
const { expect } = require("chai");

describe("Slot", function () {
  async function CCSFixture() {
    const experiationTime = (await time.latest()) + 365 * 24 * 60 * 60;

    const [owner, authority, user, others] = await ethers.getSigners();

    const CCS = await ethers.getContractFactory("CCS");
    const ccs = await CCS.deploy();

    return { ccs, owner, authority, user, others, experiationTime };
  }

  describe("Defination", function () {
    it("Should successfully define a slot", async function () {
      const { ccs } = await loadFixture(CCSFixture);

      await ccs.slotDefine(3525, "ERC3525");

      expect(await ccs.slotInfo(3525)).to.equal("ERC3525");
      expect(await ccs.slotInfo(721)).to.equal("UNDEFINED");
    });

    it("Should be reverted if called by non-owner", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      await expect(
        ccs.connect(others).slotDefine(3525, "ERC3525")
      ).to.be.revertedWith("only the owner can define slots");
    });

    it("Should be reverted if given empty name or domain", async function () {
      const { ccs } = await loadFixture(CCSFixture);

      await expect(ccs.slotDefine(3525, "")).to.be.revertedWith(
        "empty slot name"
      );
    });
  });

  describe("Allocation", function () {
    it("Should successfully allocate a slot to the corresponding authority", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await ccs.slotDefine(3525, "ERC3525");

      await ccs.slotAllocate(3525, authority.address);

      // Also test authorityInfo here
      const info = await ccs.authorityInfo(authority.address);

      expect(info.slots).to.deep.equal([3525]);
    });

    it("Should be reverted if called by non-owner", async function () {
      const { ccs, authority, others } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await ccs.slotDefine(3525, "ERC3525");

      await expect(
        ccs.connect(others).slotAllocate(3525, authority.address)
      ).to.be.revertedWith("only the owner can allocate slots");
    });

    it("Should be reverted if given the unregistered account", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      await ccs.slotDefine(3525, "ERC3525");

      await expect(ccs.slotAllocate(3525, others.address)).to.be.revertedWith(
        "authority is never registered"
      );
    });

    it("Should be reverted if given the undefined slot", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await expect(ccs.slotAllocate(3525, authority.address)).to.be.revertedWith(
        "slot is never defined"
      );
    });
  });
});
