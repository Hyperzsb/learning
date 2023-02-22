// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/utils/Strings.sol";
import "./CCSToken.sol";

contract CCS is CCSToken {
    using Strings for address;
    using Strings for uint256;

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

    /**
     * @notice This part is for appearance-related features, including
     *  - Contract URI
     *  - Slot URI
     *  - Token URI
     */

    /**
     * @dev Returns the URI for this contract, which provides metadata about the contract itself.
     * @return A string representing the URI for this contract.
     */
    function contractURI()
        public
        view
        virtual
        override
        returns (string memory)
    {
        return
            string(
                abi.encodePacked(
                    '<svg fill="none" viewBox="0 0 600 600" width="600" height="600" xmlns="http://www.w3.org/2000/svg">'
                    '  <foreignObject width="100%" height="100%">'
                    '    <div xmlns="http://www.w3.org/1999/xhtml">'
                    "      <style>"
                    "        .container {"
                    "          width: 600px;"
                    "          height: 600px;"
                    "          display: flex;"
                    "          flex-direction: column;"
                    "          justify-content: center;"
                    "          align-items: center;"
                    "          background-color: white;"
                    "          color: black;"
                    "          text-align: center;"
                    "        }"
                    "      </style>"
                    '      <div class="container">'
                    "        <h1>Community Credit System (CCS)</h1>"
                    "        <p>This contract is used to provide credit and reputation services to users in the community based on ERC3525 Semi-Fungible Token standard.</p>"
                    "        <p><b>Address: </b>",
                    address(this).toHexString(),
                    "</p>"
                    "        <p><b>Owner: </b>",
                    address(owner).toHexString(),
                    "</p>"
                    "      </div>"
                    "    </div>"
                    "  </foreignObject>"
                    "</svg>"
                )
            );
    }

    /**
     * @dev Returns the URI for the specified slot, which provides metadata about the slot.
     * @param _slot The ID of the slot.
     * @return A string representing the URI for the specified slot.
     */
    function slotURI(
        uint256 _slot
    ) public view virtual override returns (string memory) {
        return
            string(
                abi.encodePacked(
                    '<svg fill="none" viewBox="0 0 600 600" width="600" height="600" xmlns="http://www.w3.org/2000/svg">'
                    '  <foreignObject width="100%" height="100%">'
                    '    <div xmlns="http://www.w3.org/1999/xhtml">'
                    "      <style>"
                    "        .container {"
                    "          width: 600px;"
                    "          height: 600px;"
                    "          display: flex;"
                    "          flex-direction: column;"
                    "          justify-content: center;"
                    "          align-items: center;"
                    "          background-color: white;"
                    "          color: black;"
                    "          text-align: center;"
                    "        }"
                    "      </style>"
                    '      <div class="container">'
                    "        <h1>CCS - Slot</h1>"
                    "        <p><b>ID: </b>",
                    _slot.toString(),
                    "</p>"
                    "        <p><b>Name: </b>",
                    slotInfo(_slot),
                    "</p>"
                    "      </div>"
                    "    </div>"
                    "  </foreignObject>"
                    "</svg>"
                )
            );
    }

    /**
     * @dev Returns the URI for the specified token, which provides metadata about the token.
     * @param _tokenId The ID of the token.
     * @return A string representing the URI for the specified token.
     */
    function tokenURI(
        uint256 _tokenId
    ) public view virtual override returns (string memory) {
        return
            string(
                abi.encodePacked(
                    '<svg fill="none" viewBox="0 0 600 600" width="600" height="600" xmlns="http://www.w3.org/2000/svg">'
                    '  <foreignObject width="100%" height="100%">'
                    '    <div xmlns="http://www.w3.org/1999/xhtml">'
                    "      <style>"
                    "        .container {"
                    "          width: 600px;"
                    "          height: 600px;"
                    "          display: flex;"
                    "          flex-direction: column;"
                    "          justify-content: center;"
                    "          align-items: center;"
                    "          background-color: white;"
                    "          color: black;"
                    "          text-align: center;"
                    "        }"
                    "      </style>"
                    '      <div class="container">'
                    "        <h1>CCS - Token</h1>"
                    "        <p><b>ID: </b>",
                    _tokenId.toString(),
                    "</p>"
                    "        <p><b>Owner: </b>",
                    ownerOf(_tokenId).toHexString(),
                    "</p>"
                    "        <p><b>From: </b>",
                    tokenFrom[_tokenId].toHexString(),
                    "</p>"
                    "        <p><b>Credit: </b>",
                    balanceOf(_tokenId).toString(),
                    "</p>"
                    "      </div>"
                    "    </div>"
                    "  </foreignObject>"
                    "</svg>"
                )
            );
    }
}
