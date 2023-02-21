const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
const { expect } = require("chai");

describe("Token", function () {
  async function CCSFixture() {
    // Set the expiration time of the authority
    const experiationTime = (await time.latest()) + 365 * 24 * 60 * 60;

    // Get the contract's signers
    const [owner, authority, user, others] = await ethers.getSigners();

    // Deploy a new instance of the `CCS` contract
    const CCS = await ethers.getContractFactory("CCS");
    const ccs = await CCS.deploy();

    // Return the contract instance and other variables as an object
    return { ccs, owner, authority, user, others, experiationTime };
  }

  describe("Mintage", function () {
    it("Should successfully mint a token", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      // Define a new slot
      await ccs.slotDefine(3525, "ERC3525");

      // Register the authority with the contract
      const name = "Authority";
      const domain = "authority.com";
      await ccs.authorityRegister(authority.address, name, domain);

      // Allocate the slot to the authority
      await ccs.slotAllocate(3525, authority.address);

      // Mint a new token and get its ID
      const slot = 3525;
      const value = 10;
      const txResponse = await ccs.connect(authority).mint(slot, value);
      const txRecipt = await txResponse.wait();
      const tokenId = txRecipt.events.find(
        (event) => event.event === "TransferValue"
      ).args._toTokenId;

      // Assert that the balance of the authority has increased
      expect(await ccs["balanceOf(address)"](authority.address)).to.equal(1);

      // Assert that the balance of the token is set to the value that was minted
      expect(await ccs["balanceOf(uint256)"](tokenId)).to.equal(value);

      // Assert that the token owner is the authority
      expect(await ccs.ownerOf(tokenId)).to.equal(authority.address);

      // Assert that the slot associated with the token is the one that was minted
      expect(await ccs.slotOf(tokenId)).to.equal(slot);
    });

    it("Should be reverted if called by non-authority", async function () {
      const { ccs, others } = await loadFixture(CCSFixture);

      // Define a new slot
      const slot = 3525;
      await ccs.slotDefine(slot, "ERC3525");

      const value = 10;
      await expect(ccs.connect(others).mint(slot, value)).to.be.revertedWith(
        "authority is never registered"
      );
    });

    it("Should revert when authority is not valid", async function () {
      // Load the fixture
      const { ccs, authority, experiationTime } = await loadFixture(CCSFixture);

      // Define a new slot
      const slot = 3525;
      await ccs.slotDefine(slot, "ERC3525");

      // Register a new authority
      const name = "Authority";
      const domain = "authority.com";
      await ccs.authorityRegister(authority.address, name, domain);

      // Allocate the new slot to the authority
      await ccs.slotAllocate(3525, authority.address);

      // Increase time past the authority's expiration time
      await time.increaseTo(experiationTime + 60);

      // Verify that a non-valid authority cannot mint a new token
      const value = 10;
      await expect(ccs.connect(authority).mint(slot, value)).to.be.revertedWith(
        "authority is not vaild"
      );
    });

    it("Should be reverted if the slot is not allocated to the authority", async function () {
      const { ccs, authority } = await loadFixture(CCSFixture);

      // Define a new slot
      const slot = 3525;
      await ccs.slotDefine(slot, "ERC3525");

      // Register the authority
      const name = "Authority";
      const domain = "authority.com";
      await ccs.authorityRegister(authority.address, name, domain);

      // Mint the token without slot allocation
      const value = 10;
      await expect(ccs.connect(authority).mint(slot, value)).to.be.revertedWith(
        "slot is not alloctaed to authority"
      );
    });
  });

  describe("Distribution", function () {
    it("Should successfully transfer token value to a user", async function () {
      const { ccs, authority, user } = await loadFixture(CCSFixture);

      await ccs.slotDefine(3525, "ERC3525");

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await ccs.slotAllocate(3525, authority.address);

      const slot = 3525;
      const originalValue = 10;

      // Get the token id by emitted event
      const mintTxResponse = await ccs
        .connect(authority)
        .mint(slot, originalValue);
      const mintTxRecipt = await mintTxResponse.wait();
      const mintTxEvents = mintTxRecipt.events.filter(
        (event) => event.event === "TransferValue"
      );
      const tokenId = mintTxEvents[0].args._toTokenId;

      const transferedValue = 1;

      // Get the token ids by emitted event
      const transferTxResponse = await ccs
        .connect(authority)
        ["transferFrom(uint256,address,uint256)"](
          tokenId,
          user.address,
          transferedValue
        );
      const transferTxRecipt = await transferTxResponse.wait();
      const transferTxEvents = transferTxRecipt.events.filter(
        (event) => event.event === "TransferValue"
      );

      const fromTokenId = transferTxEvents[1].args._fromTokenId;
      const toTokenId = transferTxEvents[1].args._toTokenId;

      expect(tokenId).to.equal(fromTokenId);

      expect(await ccs["balanceOf(uint256)"](fromTokenId)).to.equal(
        originalValue - transferedValue
      );

      expect(await ccs["balanceOf(uint256)"](toTokenId)).to.equal(
        transferedValue
      );
    });

    it("Should be reverted if called by non-authority", async function () {
      const { ccs, authority, user, others } = await loadFixture(CCSFixture);

      await ccs.slotDefine(3525, "ERC3525");

      const name = "Authority";
      const domain = "authority.com";

      await ccs.authorityRegister(authority.address, name, domain);

      await ccs.slotAllocate(3525, authority.address);

      const slot = 3525;
      const originalValue = 10;

      // Get the token id by emitted event
      const mintTxResponse = await ccs
        .connect(authority)
        .mint(slot, originalValue);
      const mintTxRecipt = await mintTxResponse.wait();
      const mintTxEvents = mintTxRecipt.events.filter(
        (event) => event.event === "TransferValue"
      );
      const tokenId = mintTxEvents[0].args._toTokenId;

      const transferedValue = 1;

      // Get the token ids by emitted event
      const transferTxResponse = await ccs
        .connect(authority)
        ["transferFrom(uint256,address,uint256)"](
          tokenId,
          user.address,
          transferedValue
        );
      const transferTxRecipt = await transferTxResponse.wait();
      const transferTxEvents = transferTxRecipt.events.filter(
        (event) => event.event === "TransferValue"
      );

      const toTokenId = transferTxEvents[1].args._toTokenId;

      await expect(
        ccs
          .connect(user)
          ["transferFrom(uint256,address,uint256)"](
            toTokenId,
            others.address,
            transferedValue
          )
      ).to.be.revertedWith("authority is never registered");
    });
  });
});