// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.17;

import "./CCSAuthorityRR.sol";

contract CCSToken is CCSAuthorityRR {
    /**
     * @notice This part is for token-related features, including
     *  - Mintage
     *  - Distribution
     */

    /// @dev This mapping is used to record the authority of each token distributed to the user
    mapping(uint256 => address) internal tokenFrom;

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
        require(isAuthorityValid(msg.sender), "authority is not valid");
        require(
            isSlotAllocatedTo(_slot, msg.sender),
            "slot is not allocated to authority"
        );

        uint256 tokenId = _mint(msg.sender, _slot, _value);
        tokenFrom[tokenId] = msg.sender;

        return tokenId;
    }

    /**
     * @notice These following functions should be overrided to deny non-authority callers
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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

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
        require(isAuthorityValid(msg.sender), "authority is not valid");

        tokenFrom[_tokenId] = msg.sender;

        super.safeTransferFrom(_from, _to, _tokenId);
    }
}
