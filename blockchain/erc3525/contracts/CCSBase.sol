// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.17;

import "@solvprotocol/erc-3525/ERC3525.sol";

contract CCSBase is ERC3525 {
    address payable public owner;

    constructor() ERC3525("Community Credit System", "CCS", 16) {
        owner = payable(msg.sender);
    }

    /**
     * @notice This part is for slot-related features, including
     *  - Definition
     *  - Allocation
     */

    mapping(uint256 => string) internal slots;

    /**
     * @dev Maps a slot number to a slot name, which can be used to provide a more meaningful label for a specific slot
     * @param _slot The slot number to define the name of
     * @param _name The name to associate with the given slot number
     * @dev Only the contract owner can define slot names
     * @dev The slot name must not be an empty string
     */
    function slotDefine(uint256 _slot, string memory _name) external {
        require(msg.sender == owner, "only the owner can define slots");
        require(bytes(_name).length > 0, "empty slot name");

        slots[_slot] = _name;
    }

    /**
     * @notice Retrieves the name associated with a specific slot number
     * @param _slot The slot number to retrieve the name of
     * @return The name associated with the given slot number, or "UNDEFINED" if the slot has not been defined
     * @dev If the slot has not been defined, the function returns "UNDEFINED"
     */
    function slotInfo(uint256 _slot) public view returns (string memory) {
        if (bytes(slots[_slot]).length == 0) {
            return "UNDEFINED";
        } else {
            return slots[_slot];
        }
    }

    /**
     * @notice Allocates the specified slot to the specified authority
     * @param _slot The slot to be allocated
     * @param _account The authority to allocate the slot to
     * @dev Only the owner can allocate slots to authorities
     * @dev The authority must be registered
     * @dev The slot must be defined before it can be allocated
     */
    function slotAllocate(uint256 _slot, address _account) external {
        require(msg.sender == owner, "only the owner can allocate slots");
        require(isAuthority(_account), "authority is never registered");
        require(bytes(slots[_slot]).length > 0, "slot is never defined");

        authorities[_account].slots.push(_slot);
    }

    /**
     * @notice Checks whether the specified slot is allocated to the specified authority
     * @param _slot The slot to check
     * @param _account The authority to check
     * @return A boolean indicating whether the slot is allocated to the authority
     * @dev The slot must be defined before it can be checked
     * @dev The authority must be registered
     */
    function isSlotAllocatedTo(
        uint256 _slot,
        address _account
    ) public view returns (bool) {
        require(bytes(slots[_slot]).length > 0, "slot is never defined");
        require(isAuthority(_account), "authority is never registered");

        for (uint256 i = 0; i < authorities[_account].slots.length; i++) {
            if (authorities[_account].slots[i] == _slot) {
                return true;
            }
        }

        return false;
    }

    /**
     * @notice This part is for authority-related features, including
     *  - Registration
     *  - Validation
     *  - Renewal
     */

    uint256 public expirationTime = 365 * 24 * 60 * 60;

    struct Authority {
        string name;
        string domain;
        /// @dev Slots where the authority can mint corresponding tokens
        uint256[] slots;
        /// @dev Status indicating whether the authority passes the registration process
        bool registered;
        /// @dev Status indicating whether the authority passes the renewal process
        bool renewed;
        /// @dev Timestamp of the last time authority register or renew its ownership of the corresponding domain
        uint256 lastCheck;
    }

    mapping(address => Authority) internal authorities;

    /**
     * @notice Changes the expiration time of the authorities' validity
     * @param _expirationTime The new expiration time, in seconds
     * @dev Only the contract owner is allowed to change the expiration time
     * @dev The expiration time must be within a reasonable range, i.e., between 1 day and 3 years
     */
    function changeExpirationTime(uint256 _expirationTime) external {
        require(
            msg.sender == owner,
            "only the owner can change the expiration time"
        );
        require(
            _expirationTime >= 24 * 60 * 60 &&
                _expirationTime <= 3 * 365 * 24 * 60 * 60,
            "expiration time should be reasonable"
        );

        expirationTime = _expirationTime;
    }

    function authorityRegistrationRequest(
        string memory _name,
        string memory _domain
    ) external payable virtual {}

    function authorityRegistrationRetrieve()
        external
        view
        virtual
        returns (string memory)
    {}

    function authorityRegistrationConfirm() external payable virtual {}

    /**
     * @notice Checks whether the given account is registered as an authority
     * @param _account The address to check
     * @return True if the account is registered as an authority, false otherwise
     */
    function isAuthority(address _account) public view returns (bool) {
        if (authorities[_account].registered) {
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
            authorities[_account].lastCheck + expirationTime < block.timestamp
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
    function authorityInfo(
        address _account
    ) public view returns (Authority memory) {
        require(isAuthority(_account), "authority is never registered");

        return authorities[_account];
    }

    function authorityRenewalRequest() external payable virtual {}

    function authorityRenewalRetrieve()
        external
        view
        virtual
        returns (string memory)
    {}

    function authorityRenewalConfirm() external payable virtual {}
}
