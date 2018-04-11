pragma solidity ^0.4.19;

import "./zeppelin-solidity/contracts/ownership/Ownable.sol";


/**
* @title Freezable Token
* @dev Token that can be freezed for chosen token holder.
*/
contract FreezableToken is Ownable {

    mapping (address => bool) public frozenList;

    event FrozenFunds(address indexed wallet, bool frozen);

    /**
    * @dev Owner can freeze the token balance for chosen token holder.
    * @param _wallet The address of token holder whose tokens to be frozen.
    */
    function freezeAccount(address _wallet) public onlyOwner {
        require(_wallet != address(0));
        frozenList[_wallet] = true;
        FrozenFunds(_wallet, true);
    }

    /**
    * @dev Owner can unfreeze the token balance for chosen token holder.
    * @param _wallet The address of token holder whose tokens to be unfrozen.
    */
    function unfreezeAccount(address _wallet) public onlyOwner {
        require(_wallet != address(0));
        frozenList[_wallet] = false;
        FrozenFunds(_wallet, false);
    }

    /**
    * @dev Check the specified token holder whether his/her token balance is frozen.
    * @param _wallet The address of token holder to check.
    */ 
    function isFrozen(address _wallet) public view returns (bool) {
        return frozenList[_wallet];
    }

}