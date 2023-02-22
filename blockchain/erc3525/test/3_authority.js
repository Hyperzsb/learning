const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
const { expect } = require("chai");

describe("Authority", function () {
  async function CCSFixture() {
    const expirationTime = (await time.latest()) + 365 * 24 * 60 * 60;

    const [owner, authority, user, others] = await ethers.getSigners();

    const CCS = await ethers.getContractFactory("CCS");
    const ccs = await CCS.deploy();

    return { ccs, owner, authority, user, others, expirationTime };
  }

  describe("Registration", function () {
    it("Should successfully register an authority by the owner", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      expect(await ccs.isAuthority(authority.address)).to.equal(true);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(true);

      // Also test authorityInfo here
      const info = await ccs.authorityInfo(authority.address);

      expect(info.name).to.equal(name);
      expect(info.domain).to.equal(domain);
      expect(info.lastCheck)
        .to.be.above(0)
        .and.to.be.below((await time.latest()) + 60);
      expect(info.slots).to.deep.equal([]);
    });

    it("Should be reverted if called by non-owner", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      await expect(
        ccs
          .connect(others)
          .authorityRegister(others.address, "Others", "others.com")
      ).to.be.revertedWith("only the owner can register authorities");
    });

    it("Should be reverted if given empty name or domain", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      await expect(
        ccs.authorityRegister(authority.address, "", "")
      ).to.be.revertedWith("empty authority name or domain");
    });
  });

  describe("Validation", function () {
    it("Should recognize the unregistered account", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      expect(await ccs.isAuthority(others.address)).to.equal(false);
      await expect(ccs.isAuthorityValid(others.address)).to.be.revertedWith(
        "authority is never registered"
      );
    });

    it("Should be valid until the expiration time", async function () {
      const { ccs, authority, expirationTime } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await time.increaseTo(expirationTime - 60);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(true);

      await time.increaseTo(expirationTime + 60);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(false);
    });
  });

  describe("Renewal", function () {
    it("Should renew the authority ownership", async function () {
      const { ccs, authority, expirationTime } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await time.increaseTo(expirationTime - 60);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(true);

      await time.increaseTo(expirationTime + 60);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(false);

      await ccs.authorityRenew(authority.address);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(true);
    });

    it("Should be reverted if called by non-owner", async function () {
      const { ccs, authority, others } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await expect(
        ccs.connect(others).authorityRenew(authority.address)
      ).to.be.revertedWith("only the owner can renew authorities");
    });

    it("Should be reverted if given the unregistered account", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      await expect(ccs.authorityRenew(others.address)).to.be.revertedWith(
        "authority is never registered"
      );
    });

    it("Should be reverted if the account is still valid", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await expect(ccs.authorityRenew(authority.address)).to.be.revertedWith(
        "authority is still valid"
      );
    });
  });
});
