pragma solidity ^0.4.19;

import "./Whitelist.sol";
import "./Administrable.sol";


contract Whitelistable is Administrable {
    Whitelist public whitelist;

    modifier whenWhitelisted(address _sender) {
        require(isWhitelisted(_sender));
        _;
    }

    /**
    * @dev Constructor for Whitelistable contract.
    */
    function Whitelistable() public {
        whitelist = new Whitelist();
    }

    /**
    * @dev Add wallet to whitelist.
    * @dev Accept request from the owner or administrator.
    * @param _wallet The address of wallet to add.
    */
    function addWalletToWhitelist(address _wallet) public onlyAdministratorOrOwner {
        whitelist.addWallet(_wallet);
    }

    /**
    * @dev Remove wallet from whitelist.
    * @dev Accept request from the owner or administrator.
    * @param _wallet The address of whitelisted wallet to remove.
    */
    function removeWalletFromWhitelist(address _wallet) public onlyAdministratorOrOwner {
        whitelist.removeWallet(_wallet);
    }

    /**
    * @dev Check the specified wallet whether it is in the whitelist.
    * @param _wallet The address of wallet to check.
    */
    function isWhitelisted(address _wallet) public view returns (bool) {
        return whitelist.isWhitelisted(_wallet);
    }
}
