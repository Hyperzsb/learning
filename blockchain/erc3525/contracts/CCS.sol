// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.17;

import "@solvprotocol/erc-3525/ERC3525.sol";
import "hardhat/console.sol";

contract CCS is ERC3525 {
    address public owner;

    constructor() ERC3525("Community Credit System", "CCS", 16) {
        owner = msg.sender;
    }

    /**
     * @notice This part is for authority-related features, including
     *  - Registration
     *  - Validation
     *  - Token mintage
     *  - Token distribution
     */

    struct Authority {
        string name;
        string domain;
        /// @dev Timestamp of the last time authority register or renew its ownership of the corresponding domain
        uint256 lastCheck;
        /// @dev Slots where the authority can mint corresponding tokens
        uint256[] authorizedSlots;
    }

    mapping(address => Authority) private authorities;

    uint public experiationTime = 365 * 24 * 60 * 60;

    /**
     * @notice Registers a new authority with the given account, name, and domain
     * @param _account The address of the authority being registered
     * @param _name The name of the authority being registered
     * @param _domain The domain of the authority being registered
     * @dev This function is currently restricted to be called only by the contract owner
     * @dev Add verification procedures using oracles to validate the ownership of the domain
     */
    function authorityRegister(
        address _account,
        string memory _name,
        string memory _domain
    ) external {
        require(msg.sender == owner, "only the owner can register authorities");
        require(
            bytes(_name).length > 0 && bytes(_domain).length > 0,
            "empty name or domain"
        );

        // TODO: Add verification procedures using oracles to validate the ownership of the domain

        authorities[_account] = Authority({
            name: _name,
            domain: _domain,
            lastCheck: block.timestamp,
            authorizedSlots: new uint256[](0)
        });
    }

    /**
     * @notice Checks whether the given account is registered as an authority
     * @param _account The address to check
     * @return True if the account is registered as an authority, false otherwise
     */
    function isAuthority(address _account) public view returns (bool) {
        if (authorities[_account].lastCheck == 0) {
            return false;
        } else {
            return true;
        }
    }

    /**
     * @notice Checks whether the given account is a registered and valid authority
     * @param _account The address to check
     * @return True if the account is a registered and valid authority, false otherwise
     * @dev The function requires the given account to be registered as an authority
     * @dev The authority is considered valid if its last check timestamp plus the expiration time is greater than the current block timestamp
     */
    function isAuthorityValid(address _account) public view returns (bool) {
        require(isAuthority(_account), "authority is never registered");

        if (
            authorities[_account].lastCheck + experiationTime < block.timestamp
        ) {
            return false;
        } else {
            return true;
        }
    }

    /**
     * @notice Retrieves the authority data associated with the given account
     * @param _account The address of the authority to retrieve data for
     * @return The Authority struct associated with the given account
     * @dev The function requires the given account to be registered as an authority
     */
    function getAuthority(
        address _account
    ) public view returns (Authority memory) {
        require(isAuthority(_account), "it is no authority");

        return authorities[_account];
    }

    /**
     * @notice Renews the registration of the authority associated with the given account
     * @param _account The address of the authority to renew
     * @dev The function is currently restricted to be called only by the contract owner
     * @dev The function requires the given account to be registered as an authority and for the authority to be expired
     * @dev Add verification procedures using oracles to validate the ownership of the domain
     */
    function authorityRenew(address _account) external {
        /// @notice This "owner-only" restriction is only temporary for simplicity
        require(msg.sender == owner, "only the owner can renew authorities");
        require(isAuthority(_account), "authority is never registered");
        require(!isAuthorityValid(_account), "authority is still vaild");

        /// @custom:todo Add verification procedures, using oracles to validate the ownership of the domain

        authorities[_account].lastCheck = block.timestamp;
    }
}
