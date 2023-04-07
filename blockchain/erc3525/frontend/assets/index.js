import { ethers } from "https://cdnjs.cloudflare.com/ajax/libs/ethers/6.2.3/ethers.min.js";

let provider = null;
let signer = null;
// Check if the MetaMask is installed
if (window.ethereum == null) {
  document.getElementById("wallet-status").innerText = "Not installed";
  // If not, use the default provider for read-only purposes
  provider = ethers.getDefaultProvider();
} else {
  document.getElementById("wallet-status").innerText =
    "Installed but not connected";
  // If so, try to connect the wallet
  provider = new ethers.BrowserProvider(window.ethereum);
  signer = await provider.getSigner();
  // Check if the wallet is connected
  if (signer != null) {
    document.getElementById("wallet-status").innerText = "Connected";
  } else {
    alert("You must connect your wallet to this app before proceeding!");
  }
  // Refresh the page when network changed
  window.ethereum.on("chainChanged", (chainId) => {
    location.reload();
  });
}

// Display the current network
const currentNetwork = await provider.getNetwork();
document.getElementById("current-network").innerText =
  currentNetwork.name.charAt(0).toUpperCase() +
  currentNetwork.name.slice(1) +
  " (" +
  currentNetwork.chainId +
  ")";
// Display the current block number
const blockNumber = await provider.getBlockNumber();
document.getElementById("block-number").innerText = blockNumber;
// Display the current account address
const accountAddress = await signer.getAddress();
document.getElementById("account-address").innerText = accountAddress;
// Display the current account balance
const accountBalance = await provider.getBalance(accountAddress);
document.getElementById("account-balance").innerText =
  ethers.formatEther(accountBalance);

// Define the ABI for some functions manually
// const abiJSON = [
//   // Contract-related APIs
//   "function name() view returns (string)",
//   "function symbol() view returns (string)",
//   "function owner() view returns (address)",
//   "function contractURI() view returns (string)",
//   // Slot-related APIs
//   "function slotDefine(uint256 _slot, string _name)",
//   "function slotInfo(uint256 _slot) view returns (string)",
//   "function slotURI(uint256 _slot) view returns (string)",
//   "function slotAllocate(uint256 _slot, address _account)",
//   "function isSlotAllocatedTo(uint256 _slot, address _account) view returns (bool)",
//   // Authority-related APIs
//   "function isAuthority(address _account) view returns (bool)",
//   "function isAuthorityValid(address _account) view returns (bool)",
//   "function authorityInfo(address _account) view returns (string, string, uint256[], bool, bool, uint256)",
// ];

// Define the ABI for all functions automatically
const fullRawJSON = await fetch("assets/abi.json");
const fullJSON = await fullRawJSON.json();
const abiJSON = fullJSON.abi;

const abi = abiJSON;

// Define the contract address
const contractAddress = "0xc5d59500a8fef16017F2F01D4286dB17C3C18D07";
// Define the contract
const contract = new ethers.Contract(contractAddress, abi, signer);

// Contract-related APIs
// Display the contract address
document.getElementById("contract-address").innerText = contractAddress;
// Display the contract name
const contractName = await contract.name();
document.getElementById("contract-name").innerText = contractName;
// Display the contract symbol
const contractSymbol = await contract.symbol();
document.getElementById("contract-symbol").innerText = contractSymbol;
// Display the contract decimals
const contractDecimals = await contract.valueDecimals();
document.getElementById("contract-decimals").innerText = contractDecimals;
// Display the contract owner
const contractOwner = await contract.owner();
document.getElementById("contract-owner").innerText = contractOwner;
// Get and display the contract URI
document
  .getElementById("contract-uri-submit")
  .addEventListener("click", async function () {
    const contractURI = await contract.contractURI();

    toggleModal("Contract URI", contractURI);
  });

// Slot-related APIs
// Define a slot
document
  .getElementById("slot-definition-submit")
  .addEventListener("click", async function () {
    const id = document.getElementById("slot-definition-id").value;
    const name = document.getElementById("slot-definition-name").value;

    const tx = await contract.slotDefine(id, name);
    const receipt = await tx.wait();

    // console.log(JSON.stringify(receipt));
    if (receipt.status == 1) {
      toggleToast(true, `Slot ${id} has been defined as ${name}.`);
    } else {
      toggleToast(false, `Slot ${id} failed to be defined as ${name}.`);
    }
  });
// Get and display the slot info
document
  .getElementById("slot-info-submit")
  .addEventListener("click", async function () {
    const id = document.getElementById("slot-info-id").value;
    const name = await contract.slotInfo(id);
    document.getElementById("slot-info-result").innerText = name;
  });
// Get and display the slot URI
document
  .getElementById("slot-uri-submit")
  .addEventListener("click", async function () {
    const slot = document.getElementById("slot-uri-id").value;
    const slotURI = await contract.slotURI(slot);

    toggleModal("Slot URI", slotURI);
  });
// Allocate a slot to an address
document
  .getElementById("slot-allocation-submit")
  .addEventListener("click", async function () {
    const id = document.getElementById("slot-allocation-id").value;
    const address = document.getElementById("slot-allocation-address").value;

    const tx = await contract.slotAllocate(id, address);
    const receipt = await tx.wait();

    if (receipt.status == 1) {
      toggleToast(true, `Slot ${id} has been allocated to ${address}.`);
    } else {
      toggleToast(false, `Slot ${id} failed to be allocated to ${address}.`);
    }
  });
// Check the allocation status of a slot
document
  .getElementById("slot-allocation-status-submit")
  .addEventListener("click", async function () {
    const id = document.getElementById("slot-allocation-status-id").value;
    const address = document.getElementById(
      "slot-allocation-status-address"
    ).value;

    const status = await contract.isSlotAllocatedTo(id, address);

    document.getElementById("slot-allocation-status-result").innerText = status;
  });

// Authority-related APIs
// Check the existence of an authority
document
  .getElementById("authority-existence-submit")
  .addEventListener("click", async function () {
    const address = document.getElementById(
      "authority-existence-address"
    ).value;
    const existence = await contract.isAuthority(address);

    document.getElementById("authority-existence-result").innerText = existence;
  });
// Check the validity of an authority
document
  .getElementById("authority-validity-submit")
  .addEventListener("click", async function () {
    const address = document.getElementById("authority-validity-address").value;
    const validity = await contract.isAuthorityValid(address);

    document.getElementById("authority-validity-result").innerText = validity;
  });
// Get the info of an authority
document
  .getElementById("authority-info-submit")
  .addEventListener("click", async function () {
    const address = document.getElementById("authority-info-address").value;
    const info = await contract.authorityInfo(address);

    document.getElementById("authority-info-result").innerText =
      info.toString();
  });

// Token-related APIs
// Mint a token
document
  .getElementById("token-mintage-submit")
  .addEventListener("click", async function () {
    const slot = document.getElementById("token-mintage-slot").value;
    const value = document.getElementById("token-mintage-value").value;
    const tx = await contract.mint(slot, value);
    const receipt = await tx.wait();

    // console.log(JSON.stringify(receipt));
    if (receipt.status == 1) {
      toggleToast(true, `A new token has been minted.`);
    } else {
      toggleToast(false, `A new token failed to be minted.`);
    }
  });
// Get the balance of an address
document
  .getElementById("account-balance-submit")
  .addEventListener("click", async function () {
    const address = document.getElementById("account-balance-address").value;
    const balance = await contract["balanceOf(address)"](address);

    document.getElementById("account-balance-result").innerText = balance;
  });
// Get the tokens of an address
document
  .getElementById("account-tokens-submit")
  .addEventListener("click", async function () {
    const address = document.getElementById("account-tokens-address").value;
    const tokens = await contract.tokensOf(address);

    document.getElementById("account-tokens-result").innerText = tokens;
  });
// Get the balance of a Token
document
  .getElementById("token-balance-submit")
  .addEventListener("click", async function () {
    const token = document.getElementById("token-balance-token").value;
    const balance = await contract["balanceOf(uint256)"](token);

    document.getElementById("token-balance-result").innerText = balance;
  });
// Transfer a token to another address
document
  .getElementById("token-transfer-submit")
  .addEventListener("click", async function () {
    const token = document.getElementById("token-transfer-token").value;
    const address = document.getElementById("token-transfer-address").value;
    const value = document.getElementById("token-transfer-value").value;
    const tx = await contract["transferFrom(uint256,address,uint256)"](
      token,
      address,
      value
    );
    const receipt = await tx.wait();

    // console.log(JSON.stringify(receipt));
    if (receipt.status == 1) {
      toggleToast(
        true,
        `A value of ${value} from token ${token} has been transferred to account ${address}.`
      );
    } else {
      toggleToast(
        false,
        `A value of ${value} from token ${token} failed to be transferred to account ${address}.`
      );
    }
  });
// Get the owner of a token
document
  .getElementById("token-owner-submit")
  .addEventListener("click", async function () {
    const token = document.getElementById("token-owner-token").value;
    const owner = await contract.ownerOf(token);

    document.getElementById("token-owner-result").innerText = owner;
  });
// Get the slot of a token
document
  .getElementById("token-slot-submit")
  .addEventListener("click", async function () {
    const token = document.getElementById("token-slot-token").value;
    const slot = await contract.slotOf(token);

    document.getElementById("token-slot-result").innerText = slot;
  });
// Get and display the token URI
document
  .getElementById("token-uri-submit")
  .addEventListener("click", async function () {
    const slot = document.getElementById("token-uri-token").value;
    const tokenURI = await contract.tokenURI(slot);

    toggleModal("Token URI", tokenURI);
  });
