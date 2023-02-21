const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
const { expect } = require("chai");

describe("Appearance", function () {
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

  describe("Contract URI", function () {
    it("Should successfully generate the URI of the contract", async function () {
      const { ccs } = await loadFixture(CCSFixture);

      const contractURI = await ccs.contractURI();

      const fs = require("fs");

      fs.writeFile("./build/assets/contract.svg", contractURI, function (err) {
        if (err) throw err;
      });
    });
  });

  describe("Slot URI", function () {
    it("Should successfully generate the URI of the slot", async function () {
      const { ccs } = await loadFixture(CCSFixture);

      await ccs.slotDefine(3525, "ERC3525");

      const slotURI = await ccs.slotURI(3525);

      const fs = require("fs");
      fs.writeFile("./build/assets/slot.svg", slotURI, function (err) {
        if (err) throw err;
      });
    });
  });

  describe("Token URI", function () {
    it("Should successfully generate the URI of the token owned by the authority", async function () {
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

      const tokenURI = await ccs.tokenURI(tokenId);

      const fs = require("fs");
      fs.writeFile("./build/assets/authorityToken.svg", tokenURI, function (err) {
        if (err) throw err;
      });
    });

    it("Should successfully generate the URI of the token owned by the user", async function () {
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

      const toTokenId = transferTxEvents[1].args._toTokenId;

      const tokenURI = await ccs.tokenURI(toTokenId);

      const fs = require("fs");
      fs.writeFile("./build/assets/userToken.svg", tokenURI, function (err) {
        if (err) throw err;
      });
    });
  });
});
