// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/utils/Strings.sol";
import "@solvprotocol/erc-3525/ERC3525.sol";
import "hardhat/console.sol";

contract CCS is ERC3525 {
    address public owner;

    constructor() ERC3525("Community Credit System", "CCS", 16) {
        owner = msg.sender;
    }

    /**
     * @notice This part is for slot-related features, including
     *  - Defination
     *  - Allocation
     */

    mapping(uint256 => string) private slots;

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

    struct Authority {
        string name;
        string domain;
        /// @dev Timestamp of the last time authority register or renew its ownership of the corresponding domain
        uint256 lastCheck;
        /// @dev Slots where the authority can mint corresponding tokens
        uint256[] slots;
    }

    mapping(address => Authority) private authorities;

    uint256 public expirationTime = 365 * 24 * 60 * 60;

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
            "empty authority name or domain"
        );

        // TODO: Add verification procedures using oracles to validate the ownership of the domain

        authorities[_account] = Authority({
            name: _name,
            domain: _domain,
            lastCheck: block.timestamp,
            slots: new uint256[](0)
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

    /**
     * @notice Changes the expiration time of the authorities' validity
     * @param _expirationTime The new expiration time, in seconds
     * @dev Only the contract owner is allowed to change the expiration time
     * @dev The expiration time must be within a reasonable range, i.e., between 1 day and 3 years
     */
    function changeExpirationTime(uint256 _expirationTime) external {
        require(
            msg.sender == owner,
            "only the owner can change expiration time"
        );
        require(
            _expirationTime >= 24 * 60 * 60 &&
                _expirationTime <= 3 * 365 * 24 * 60 * 60,
            "expiration time should be reasonable"
        );

        expirationTime = _expirationTime;
    }

    /**
     * @notice This part is for token-related features, including
     *  - Mintage
     *  - Distribution
     */

    /// @dev This mapping is used to record the authority of each token distributed to the user
    mapping(uint256 => address) private tokenFrom;

    /**
     * @notice Mints a new token with the specified value and assigns it to the calling authority in the specified slot
     * @param _slot The slot in which the token will be created
     * @param _value The value of the token to be created
     * @return The ID of the newly minted token
     * @dev Only an authority can mint a new token and assign it to a slot
     * @dev The authority must be registered and valid
     * @dev The slot must be allocated to the calling authority
     */
    function mint(uint256 _slot, uint256 _value) external returns (uint256) {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");
        require(
            isSlotAllocatedTo(_slot, msg.sender),
            "slot is not alloctaed to authority"
        );

        uint256 tokenId = _mint(msg.sender, _slot, _value);
        tokenFrom[tokenId] = msg.sender;

        return tokenId;
    }

    /**
     * @notice These following functions should be overrided to deny non-suthority callers
     */

    /**
     * @notice Approves another address to transfer the specified token on behalf of the token owner
     * @param _tokenId The ID of the token to be approved for transfer
     * @param _to The address that will be approved to transfer the token
     * @param _value The amount of tokens to be approved for transfer
     * @dev Only an authority can approve a transfer of tokens on behalf of the owner
     * @dev The authority must be registered and valid
     */
    function approve(
        uint256 _tokenId,
        address _to,
        uint256 _value
    ) public payable virtual override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        super.approve(_tokenId, _to, _value);
    }

    /**
     * @notice Approves another address to transfer the specified token on behalf of the token owner
     * @param _to The address that will be approved to transfer the token
     * @param _tokenId The ID of the token to be approved for transfer
     * @dev Only an authority can approve a transfer of tokens on behalf of the owner
     * @dev The authority must be registered and valid
     */
    function approve(
        address _to,
        uint256 _tokenId
    ) public payable virtual override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        super.approve(_to, _tokenId);
    }

    /**
     * @notice Gets the approved address for a token ID
     * @param _tokenId The ID of the token to get the approved address for
     * @return The approved address for the given token ID
     * @dev Only an authority can get the approved address for a token
     * @dev The authority must be registered and valid
     */
    function getApproved(
        uint256 _tokenId
    ) public view virtual override returns (address) {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        return super.getApproved(_tokenId);
    }

    /**
     * @notice Sets or unsets the approval of a given operator
     * @param _operator The operator whose approval is to be set or unset
     * @param _approved True if the operator is to be approved, false if the approval is to be revoked
     * @dev Only an authority can set or unset the approval of an operator
     * @dev The authority must be registered and valid
     */
    function setApprovalForAll(
        address _operator,
        bool _approved
    ) public virtual override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        super.setApprovalForAll(_operator, _approved);
    }

    /**
     * @notice Returns whether the given operator is approved by a given owner
     * @param _owner The owner of the tokens
     * @param _operator The operator to check for approval
     * @return True if the operator is approved by the owner, false otherwise
     * @dev Only an authority can check whether an operator is approved by an owner
     * @dev The authority must be registered and valid
     */
    function isApprovedForAll(
        address _owner,
        address _operator
    ) public view virtual override returns (bool) {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        return super.isApprovedForAll(_owner, _operator);
    }

    /**
     * @notice Returns the amount of tokens owned by a given address that are approved for transfer by another address
     * @param _tokenId The ID of the token to check the allowance for
     * @param _operator The address of the operator to check the allowance for
     * @return The amount of tokens approved for transfer by the operator for the given token ID
     * @dev Only an authority can check the allowance for a token
     * @dev The authority must be registered and valid
     */
    function allowance(
        uint256 _tokenId,
        address _operator
    ) public view virtual override returns (uint256) {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        return super.allowance(_tokenId, _operator);
    }

    /**
     * @notice Transfers some value of a given token to another address
     * @param _fromTokenId The ID of the token to be transferred
     * @param _to The address to transfer the token to
     * @param _value The amount of value to transfer
     * @return The new token ID of the transferred token
     * @dev Only an authority can transfer a token
     * @dev The authority must be registered and valid
     */
    function transferFrom(
        uint256 _fromTokenId,
        address _to,
        uint256 _value
    ) public payable virtual override returns (uint256) {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        uint256 tokenId = super.transferFrom(_fromTokenId, _to, _value);
        tokenFrom[tokenId] = msg.sender;

        return tokenId;
    }

    /**
     * @notice Transfers some value of a given token to another token
     * @param _fromTokenId The ID of the token to be transferred
     * @param _toTokenId The ID of the token to transfer to
     * @param _value The amount of value to transfer
     * @dev Only an authority can transfer a token
     * @dev The authority must be registered and valid
     */
    function transferFrom(
        uint256 _fromTokenId,
        uint256 _toTokenId,
        uint256 _value
    ) public payable virtual override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        tokenFrom[_toTokenId] = msg.sender;

        super.transferFrom(_fromTokenId, _toTokenId, _value);
    }

    /**
     * @notice Transfers the ownership of a given token to another address
     * @param _from The address of the owner of the token
     * @param _to The address to transfer the token to
     * @param _tokenId The ID of the token to be transferred
     * @dev Only an authority can transfer a token
     * @dev The authority must be registered and valid
     */
    function transferFrom(
        address _from,
        address _to,
        uint256 _tokenId
    ) public payable virtual override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        tokenFrom[_tokenId] = msg.sender;

        super.transferFrom(_from, _to, _tokenId);
    }

    /**
     * @notice Safely transfers the ownership of a given token ID to another address
     * @param _from The address of the owner of the token
     * @param _to The address to transfer the token to
     * @param _tokenId The ID of the token to be transferred
     * @param _data Additional data with no specified format to be passed to the receiver contract
     * @dev Only an authority can safely transfer a token
     * @dev The authority must be registered and valid
     */
    function safeTransferFrom(
        address _from,
        address _to,
        uint256 _tokenId,
        bytes memory _data
    ) public payable virtual override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        tokenFrom[_tokenId] = msg.sender;

        super.safeTransferFrom(_from, _to, _tokenId, _data);
    }

    /**
     * @notice Safely transfers the ownership of a given token ID to another address
     * @param _from The address of the owner of the token
     * @param _to The address to transfer the token to
     * @param _tokenId The ID of the token to be transferred
     * @dev Only an authority can safely transfer a token
     * @dev The authority must be registered and valid
     */
    function safeTransferFrom(
        address _from,
        address _to,
        uint256 _tokenId
    ) public payable virtual override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(isAuthorityValid(msg.sender), "authority is not vaild");

        tokenFrom[_tokenId] = msg.sender;

        super.safeTransferFrom(_from, _to, _tokenId);
    }

    /**
     * @notice This part is for user-related features, including
     *  - Slot Info
     *  - Token Info
     */

    /**
     * @dev Returns an array of all the unique slots in which the specified address has some tokens
     * @param _account The address to query.
     * @return An array of uint256 values representing the unique slots in which the specified address has some tokens
     */
    function slotsOf(address _account) public view returns (uint256[] memory) {
        uint256 balance = balanceOf(_account);
        uint256[] memory allSlots = new uint256[](balance);
        uint256 uniqueSlotCount = 0;

        for (uint256 i = 0; i < balance; i++) {
            allSlots[i] = (slotOf(tokenOfOwnerByIndex(_account, i)));

            bool isUnique = true;
            for (uint256 j = 0; j < i; j++) {
                if (allSlots[i] == allSlots[j]) {
                    isUnique = false;
                    break;
                }
            }

            if (isUnique) {
                uniqueSlotCount++;
            }
        }

        uint256[] memory uniqueSlots = new uint256[](uniqueSlotCount);
        uint256 idx = 0;

        for (uint256 i = 0; i < balance; i++) {
            bool isUnique = true;
            for (uint256 j = 0; j < i; j++) {
                if (allSlots[i] == allSlots[j]) {
                    isUnique = false;
                    break;
                }
            }

            if (isUnique) {
                uniqueSlots[idx] = allSlots[i];
                idx++;
            }
        }

        return uniqueSlots;
    }

    /**
     * @dev Returns an array of all the tokens owned by the specified address
     * @param _account The address to query.
     * @return An array of uint256 values representing the tokens owned by the specified address
     */
    function tokensOf(address _account) public view returns (uint256[] memory) {
        uint256 balance = balanceOf(_account);
        uint256[] memory allTokens = new uint256[](balance);

        for (uint256 i = 0; i < balance; i++) {
            allTokens[i] = tokenOfOwnerByIndex(_account, i);
        }

        return allTokens;
    }
}
