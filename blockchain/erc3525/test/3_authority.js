const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
const { expect } = require("chai");

describe.skip("Authority", function () {
  async function CCSFixture() {
    // Set the expiration time of the authority
    const expirationTime = (await time.latest()) + 365 * 24 * 60 * 60;

    // Get the contract's signers
    const [owner, authority, others] = await ethers.getSigners();

    // Deploy a new instance of the `CCS` contract
    const CCS = await ethers.getContractFactory("CCS");
    const ccs = await CCS.deploy();

    // Return the contract instance and other variables as an object
    return { ccs, owner, authority, others, expirationTime };
  }

  describe("Registration", function () {
    it("Should register an authority by the owner", async function () {
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

    it("Should be reverted if called by a non-owner", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      await expect(
        ccs
          .connect(others)
          .authorityRegister(others.address, "Others", "others.com")
      ).to.be.revertedWith("only the owner can register authorities");
    });

    it("Should be reverted if given an empty name or domain", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      await expect(
        ccs.authorityRegister(authority.address, "", "")
      ).to.be.revertedWith("empty authority name or domain");
    });
  });

  describe("Validation", function () {
    it("Should recognize an unregistered account", async function () {
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
    it("Should renew the authority ownership by the owner", async function () {
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

    it("Should be reverted if called by a non-owner", async function () {
      const { ccs, authority, others } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";
      await ccs.authorityRegister(authority.address, name, domain);

      await expect(
        ccs.connect(others).authorityRenew(authority.address)
      ).to.be.revertedWith("only the owner can renew authorities");
    });

    it("Should be reverted if given an unregistered account", async function () {
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

  describe("Expiration Time", function () {
    it("Should change the expiration time by the owner", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      const name = "Authority";
      const domain = "authority.com";
      await ccs.authorityRegister(authority.address, name, domain);

      const newExpirationPeriod = 30 * 24 * 60 * 60;
      const newExpirationTime = (await time.latest()) + 30 * 24 * 60 * 60;

      await ccs.changeExpirationTime(newExpirationPeriod);

      await time.increaseTo(newExpirationTime - 60);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(true);

      await time.increaseTo(newExpirationTime + 60);
      expect(await ccs.isAuthorityValid(authority.address)).to.equal(false);
    });

    it("Should be reverted if called by a non-owner", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      const newExpirationPeriod = 30 * 24 * 60 * 60;

      await expect(
        ccs.connect(others).changeExpirationTime(newExpirationPeriod)
      ).to.be.revertedWith("only the owner can change the expiration time");
    });

    it("Should be reverted if the new expiration time is non-reasonable", async function () {
      const { ccs } = await loadFixture(CCSFixture);

      const smallExpirationPeriod = 60 * 60;
      const largeExpirationPeriod = 5 * 365 * 24 * 60 * 60;

      await expect(
        ccs.changeExpirationTime(smallExpirationPeriod)
      ).to.be.revertedWith("expiration time should be reasonable");

      await expect(
        ccs.changeExpirationTime(largeExpirationPeriod)
      ).to.be.revertedWith("expiration time should be reasonable");
    });
  });
});
