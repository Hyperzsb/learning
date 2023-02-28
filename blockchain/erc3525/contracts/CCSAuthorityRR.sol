// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/utils/Strings.sol";
import "@chainlink/contracts/src/v0.8/interfaces/VRFCoordinatorV2Interface.sol";
import "@chainlink/contracts/src/v0.8/VRFConsumerBaseV2.sol";
import "./CCSBase.sol";

contract CCSAuthorityRR is CCSBase, VRFConsumerBaseV2 {
    using Strings for uint256;

    VRFCoordinatorV2Interface private VRFCOORDINATOR;

    struct VRFConfig {
        address coordinator;
        address link;
        bytes32 keyHash;
        uint64 subscriptionId;
        uint32 callbackGasLimit;
        uint16 requestConfirmations;
        uint32 numWords;
    }

    /// @dev Currently using hardcode for simplicity
    VRFConfig private vrfConfig =
        VRFConfig({
            coordinator: 0x2Ca8E0C643bDe4C2E08ab1fA0da3401AdAD7734D,
            link: 0x326C977E6efc84E512bB9C30f76E30c160eD06FB,
            keyHash: 0x79d3d8832d904592c0bf9818b621522c988bb8b0c05cdc3b15aea1b6e8db0c15,
            subscriptionId: 10381,
            callbackGasLimit: 100000,
            requestConfirmations: 3,
            numWords: 1
        });

    constructor() VRFConsumerBaseV2(vrfConfig.coordinator) {
        VRFCOORDINATOR = VRFCoordinatorV2Interface(vrfConfig.coordinator);
    }

    mapping(uint256 => address) private requestToAddress;
    mapping(address => uint256) private addressToRequest;
    mapping(address => string) private randomStrings;

    function authorityRegistrationRequest(
        string memory _name,
        string memory _domain
    ) external payable override {
        require(!isAuthority(msg.sender), "authority is already registered");
        require(
            addressToRequest[msg.sender] == 0,
            "registration is already initialized"
        );

        authorities[msg.sender] = Authority({
            name: _name,
            domain: _domain,
            slots: new uint256[](0),
            registered: false,
            renewed: false,
            lastCheck: 0
        });

        uint256 requestId = VRFCOORDINATOR.requestRandomWords(
            vrfConfig.keyHash,
            vrfConfig.subscriptionId,
            vrfConfig.requestConfirmations,
            vrfConfig.callbackGasLimit,
            vrfConfig.numWords
        );

        requestToAddress[requestId] = msg.sender;
        addressToRequest[msg.sender] = requestId;
    }

    function authorityRegistrationRetrieve()
        external
        view
        override
        returns (string memory)
    {
        require(!isAuthority(msg.sender), "authority is already registered");
        require(
            addressToRequest[msg.sender] > 0,
            "registration is not initialized"
        );
        require(
            bytes(randomStrings[msg.sender]).length > 0,
            "random string is not ready"
        );

        return randomStrings[msg.sender];
    }

    function authorityRegistrationConfirm() external payable override {
        require(!isAuthority(msg.sender), "authority is already registered");
        require(
            addressToRequest[msg.sender] > 0,
            "registration is not initialized"
        );
        require(
            bytes(randomStrings[msg.sender]).length > 0,
            "random string is not ready"
        );
        require(isDNSRecordMatched(msg.sender), "DNS record is not matched");

        delete requestToAddress[addressToRequest[msg.sender]];
        delete addressToRequest[msg.sender];
        delete randomStrings[msg.sender];

        authorities[msg.sender].renewed = true;
        authorities[msg.sender].lastCheck = block.timestamp;
        authorities[msg.sender].registered = true;
    }

    function authorityRenewalRequest() external payable override {
        require(isAuthority(msg.sender), "authority is never registered");
        require(!isAuthorityValid(msg.sender), "authority is still valid");
        require(
            addressToRequest[msg.sender] == 0,
            "renewal is already initialized"
        );

        authorities[msg.sender].renewed = false;

        uint256 requestId = VRFCOORDINATOR.requestRandomWords(
            vrfConfig.keyHash,
            vrfConfig.subscriptionId,
            vrfConfig.requestConfirmations,
            vrfConfig.callbackGasLimit,
            vrfConfig.numWords
        );

        requestToAddress[requestId] = msg.sender;
        addressToRequest[msg.sender] = requestId;
    }

    function authorityRenewalRetrieve()
        external
        view
        override
        returns (string memory)
    {
        require(isAuthority(msg.sender), "authority is never registered");
        require(!isAuthorityValid(msg.sender), "authority is still valid");
        require(
            !authorities[msg.sender].renewed &&
                addressToRequest[msg.sender] > 0,
            "renewal is not initialized"
        );
        require(
            bytes(randomStrings[msg.sender]).length > 0,
            "random string is not ready"
        );

        return randomStrings[msg.sender];
    }

    function authorityRenewalConfirm() external payable override {
        require(isAuthority(msg.sender), "authority is already registered");
        require(!isAuthorityValid(msg.sender), "authority is still valid");
        require(
            !authorities[msg.sender].renewed &&
                addressToRequest[msg.sender] > 0,
            "renewal is not initialized"
        );
        require(
            bytes(randomStrings[msg.sender]).length > 0,
            "random string is not ready"
        );
        require(isDNSRecordMatched(msg.sender), "DNS record is not matched");

        delete requestToAddress[addressToRequest[msg.sender]];
        delete addressToRequest[msg.sender];
        delete randomStrings[msg.sender];

        authorities[msg.sender].renewed = true;
        authorities[msg.sender].lastCheck = block.timestamp;
    }

    function fulfillRandomWords(
        uint256 _requestId,
        uint256[] memory _randomWords
    ) internal override {
        randomStrings[requestToAddress[_requestId]] = _randomWords[0].toString();
    }

    function isDNSRecordMatched(address _account) private pure returns (bool) {
        _account;
        return true;
    }
}
