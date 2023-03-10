const { loadFixture } = require("@nomicfoundation/hardhat-network-helpers");
const { expect } = require("chai");

describe("User", function () {
  async function CCSFixture() {
    // Get the contract's signers
    const [owner, authority, user] = await ethers.getSigners();

    // Deploy a new instance of the `CCS` contract
    const CCS = await ethers.getContractFactory("CCS");
    const ccs = await CCS.deploy();

    // Return the contract instance and other variables as an object
    return { ccs, owner, authority, user };
  }

  describe.skip("Slot Info", function () {
    it("Should get the slot info of the user", async function () {
      const { ccs, authority, user } = await loadFixture(CCSFixture);

      // Define some slots
      await ccs.slotDefine(1155, "ERC1155");
      await ccs.slotDefine(3525, "ERC3525");

      // Register an authority
      const name = "Authority";
      const domain = "authority.com";
      await ccs.authorityRegister(authority.address, name, domain);

      // Allocate slots to the authority
      await ccs.slotAllocate(1155, authority.address);
      await ccs.slotAllocate(3525, authority.address);

      // Mint some tokens in each slot
      const mint1TxResponse = await ccs.connect(authority).mint(1155, 10);
      const mint1TxReceipt = await mint1TxResponse.wait();
      const mint1TxEvents = mint1TxReceipt.events.filter(
        (event) => event.event === "TransferValue"
      );
      const originalTokenId1 = mint1TxEvents[0].args._toTokenId;

      const mint2TxResponse = await ccs.connect(authority).mint(3525, 10);
      const mint2TxReceipt = await mint2TxResponse.wait();
      const mint2TxEvents = mint2TxReceipt.events.filter(
        (event) => event.event === "TransferValue"
      );
      const originalTokenId2 = mint2TxEvents[0].args._toTokenId;

      const mint3TxResponse = await ccs.connect(authority).mint(3525, 10);
      const mint3TxReceipt = await mint3TxResponse.wait();
      const mint3TxEvents = mint3TxReceipt.events.filter(
        (event) => event.event === "TransferValue"
      );
      const originalTokenId3 = mint3TxEvents[0].args._toTokenId;

      // Transfer tokens to the user
      await ccs
        .connect(authority)
        ["transferFrom(uint256,address,uint256)"](
          originalTokenId1,
          user.address,
          1
        );

      await ccs
        .connect(authority)
        ["transferFrom(uint256,address,uint256)"](
          originalTokenId2,
          user.address,
          2
        );

      await ccs
        .connect(authority)
        ["transferFrom(uint256,address,uint256)"](
          originalTokenId3,
          user.address,
          3
        );

      const slots = await ccs.slotsOf(user.address);

      expect(slots.length).to.equal(2);

      expect(slots).to.deep.equal([1155, 3525]);
    });
  });

  describe("Token Info", function () {
    it("Should get the token info of the user", async function () {
      const { ccs, authority, user } = await loadFixture(CCSFixture);

      // Define some slots
      await ccs.slotDefine(1155, "ERC1155");
      await ccs.slotDefine(3525, "ERC3525");

      // Register an authority
      const name = "Authority";
      const domain = "authority.com";
      await ccs.authorityRegister(authority.address, name, domain);

      // Allocate slots to the authority
      await ccs.slotAllocate(1155, authority.address);
      await ccs.slotAllocate(3525, authority.address);

      // Mint some tokens in each slot
      const mint1TxResponse = await ccs.connect(authority).mint(1155, 10);
      const mint1TxReceipt = await mint1TxResponse.wait();
      const mint1TxEvents = mint1TxReceipt.events.filter(
        (event) => event.event === "TransferValue"
      );
      const originalTokenId1 = mint1TxEvents[0].args._toTokenId;

      const mint2TxResponse = await ccs.connect(authority).mint(3525, 10);
      const mint2TxReceipt = await mint2TxResponse.wait();
      const mint2TxEvents = mint2TxReceipt.events.filter(
        (event) => event.event === "TransferValue"
      );
      const originalTokenId2 = mint2TxEvents[0].args._toTokenId;

      // Transfer tokens to the user
      await ccs
        .connect(authority)
        ["transferFrom(uint256,address,uint256)"](
          originalTokenId1,
          user.address,
          1
        );

      await ccs
        .connect(authority)
        ["transferFrom(uint256,address,uint256)"](
          originalTokenId2,
          user.address,
          2
        );

      const tokenIds = await ccs.connect(user).tokensOf(user.address);

      expect(await ccs.slotOf(tokenIds[0])).to.equal(1155);

      expect(await ccs["balanceOf(uint256)"](tokenIds[0])).to.equal(1);

      expect(await ccs.slotOf(tokenIds[1])).to.equal(3525);

      expect(await ccs["balanceOf(uint256)"](tokenIds[1])).to.equal(2);
    });
  });
});
