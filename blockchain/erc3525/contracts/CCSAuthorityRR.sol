// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/utils/Strings.sol";
import "@chainlink/contracts/src/v0.8/interfaces/VRFCoordinatorV2Interface.sol";
import "@chainlink/contracts/src/v0.8/VRFConsumerBaseV2.sol";
import "@chainlink/contracts/src/v0.8/ChainlinkClient.sol";
import "./CCSBase.sol";

contract CCSAuthorityRR is CCSBase, VRFConsumerBaseV2, ChainlinkClient {
    using Strings for uint256;
    using Chainlink for Chainlink.Request;

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

    /**
     * @dev This struct stores configs for using Chainlink VRF service
     * @dev Currently using hardcode for simplicity
     * @dev This config is only valid on Goerli testnet.
     * @dev For configs of other testnets or mainnet, please refer to https://docs.chain.link/vrf/v2/subscription/supported-networks
     */
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

    function changeSubscriptionId(uint64 _subscriptionId) external {
        require(msg.sender == owner);

        vrfConfig.subscriptionId = _subscriptionId;
    }

    struct APIConfig {
        address oracle;
        address link;
        bytes32 jobId;
        uint256 fee;
    }

    /**
     * @dev This struct stores configs for using Chainlink API call service
     * @dev Currently using hardcode for simplicity
     * @dev This config is only valid on Goerli testnet
     * @dev For configs of other testnets or mainnet, please refer to https://docs.chain.link/any-api/testnet-oracles
     */
    APIConfig private apiConfig =
        APIConfig({
            oracle: 0xCC79157eb46F5624204f47AB42b3906cAA40eaB7,
            link: 0x326C977E6efc84E512bB9C30f76E30c160eD06FB,
            jobId: "7d80a6386ef543a3abb52817f6707e3b",
            fee: (1 * LINK_DIVISIBILITY) / 10
        });

    constructor() VRFConsumerBaseV2(vrfConfig.coordinator) {
        /// @dev This is for VRF service initialization
        /// @dev For more details, please refer to https://docs.chain.link/vrf/v2/subscription/examples/get-a-random-number
        VRFCOORDINATOR = VRFCoordinatorV2Interface(vrfConfig.coordinator);
        /// @dev This is for API call service initialization
        /// @dev For more details, please refer to https://docs.chain.link/any-api/get-request/examples/array-response
        setChainlinkToken(apiConfig.link);
        setChainlinkOracle(apiConfig.oracle);
    }

    mapping(uint256 => address) private vrfRToA;
    mapping(address => uint256) private vrfAToR;
    mapping(address => string) private originalRandomStrings;

    mapping(bytes32 => address) private apiRToA;
    mapping(address => bytes32) private apiAToR;
    mapping(address => string) private targetDNSRecords;

    function authorityRegistrationRequest(
        string memory _name,
        string memory _domain
    ) external payable {
        require(!isAuthority(msg.sender));
        require(vrfAToR[msg.sender] == 0);

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

        vrfRToA[requestId] = msg.sender;
        vrfAToR[msg.sender] = requestId;
    }

    function authorityRegistrationRetrieve()
        external
        view
        returns (string memory)
    {
        require(!isAuthority(msg.sender));
        require(vrfAToR[msg.sender] > 0);
        require(bytes(originalRandomStrings[msg.sender]).length > 0);

        return originalRandomStrings[msg.sender];
    }

    function authorityRegistrationVerify() external payable {
        require(!isAuthority(msg.sender));
        require(vrfAToR[msg.sender] > 0);
        require(bytes(originalRandomStrings[msg.sender]).length > 0);
        require(apiAToR[msg.sender] == 0);

        getDNSRecord(msg.sender);
    }

    function authorityRegistrationConfirm() external payable {
        require(!isAuthority(msg.sender));
        require(vrfAToR[msg.sender] > 0);
        require(bytes(originalRandomStrings[msg.sender]).length > 0);
        require(bytes(targetDNSRecords[msg.sender]).length > 0);
        require(
            keccak256(abi.encodePacked(originalRandomStrings[msg.sender])) ==
                keccak256(abi.encodePacked(targetDNSRecords[msg.sender]))
        );

        delete vrfRToA[vrfAToR[msg.sender]];
        delete vrfAToR[msg.sender];
        delete originalRandomStrings[msg.sender];

        delete apiRToA[apiAToR[msg.sender]];
        delete apiAToR[msg.sender];
        delete targetDNSRecords[msg.sender];

        authorities[msg.sender].renewed = true;
        authorities[msg.sender].lastCheck = block.timestamp;
        authorities[msg.sender].registered = true;
    }

    function authorityRenewalRequest() external payable {
        require(isAuthority(msg.sender));
        require(!isAuthorityValid(msg.sender));
        require(vrfAToR[msg.sender] == 0);

        authorities[msg.sender].renewed = false;

        uint256 requestId = VRFCOORDINATOR.requestRandomWords(
            vrfConfig.keyHash,
            vrfConfig.subscriptionId,
            vrfConfig.requestConfirmations,
            vrfConfig.callbackGasLimit,
            vrfConfig.numWords
        );

        vrfRToA[requestId] = msg.sender;
        vrfAToR[msg.sender] = requestId;
    }

    function authorityRenewalRetrieve() external view returns (string memory) {
        require(isAuthority(msg.sender));
        require(!isAuthorityValid(msg.sender));
        require(!authorities[msg.sender].renewed && vrfAToR[msg.sender] > 0);
        require(bytes(originalRandomStrings[msg.sender]).length > 0);

        return originalRandomStrings[msg.sender];
    }

    function authorityRenewalVerify() external payable {
        require(isAuthority(msg.sender));
        require(!isAuthorityValid(msg.sender));
        require(!authorities[msg.sender].renewed && vrfAToR[msg.sender] > 0);
        require(bytes(originalRandomStrings[msg.sender]).length > 0);
        require(apiAToR[msg.sender] == 0);

        getDNSRecord(msg.sender);
    }

    function authorityRenewalConfirm() external payable {
        require(isAuthority(msg.sender));
        require(!isAuthorityValid(msg.sender));
        require(!authorities[msg.sender].renewed && vrfAToR[msg.sender] > 0);
        require(bytes(originalRandomStrings[msg.sender]).length > 0);
        require(bytes(targetDNSRecords[msg.sender]).length > 0);
        require(
            keccak256(abi.encodePacked(originalRandomStrings[msg.sender])) ==
                keccak256(abi.encodePacked(targetDNSRecords[msg.sender]))
        );

        delete vrfRToA[vrfAToR[msg.sender]];
        delete vrfAToR[msg.sender];
        delete originalRandomStrings[msg.sender];

        delete apiRToA[apiAToR[msg.sender]];
        delete apiAToR[msg.sender];
        delete targetDNSRecords[msg.sender];

        authorities[msg.sender].renewed = true;
        authorities[msg.sender].lastCheck = block.timestamp;
    }

    function getDNSRecord(address _account) private {
        Chainlink.Request memory request = buildChainlinkRequest(
            apiConfig.jobId,
            address(this),
            this.fulfillAPICalls.selector
        );

        // Prepare the target API URL
        request.add(
            "get",
            string(
                abi.encodePacked(
                    "https://dns.google/resolve?name=",
                    authorities[_account].domain,
                    "&type=txt"
                )
            )
        );

        // Prepare the target JSON entry
        request.add("path", "Answer,data");

        bytes32 requestId = sendChainlinkRequest(request, apiConfig.fee);
        apiRToA[requestId] = _account;
        apiAToR[_account] = requestId;
    }

    function fulfillRandomWords(
        uint256 _requestId,
        uint256[] memory _randomWords
    ) internal override {
        originalRandomStrings[vrfRToA[_requestId]] = _randomWords[0].toString();
    }

    function fulfillAPICalls(
        bytes32 _requestId,
        string memory _record
    ) public recordChainlinkFulfillment(_requestId) {
        targetDNSRecords[apiRToA[_requestId]] = _record;
    }

    function withdrawLink() public {
        require(msg.sender == owner);

        LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());

        require(link.transfer(msg.sender, link.balanceOf(address(this))));
    }
}
