pragma solidity ^0.4.15;

import './zeppelin-solidity/contracts/token/StandardToken.sol';
import './zeppelin-solidity/contracts/token/BurnableToken.sol';
import './zeppelin-solidity/contracts/lifecycle/Pausable.sol';

import './FreezableToken.sol';


interface tokenRecipient {
    function receiveApproval(
        address _from,
        uint256 _value,
        address _token,
        bytes _extraData)
    public;
}


contract MocrowCoin is StandardToken, BurnableToken, FreezableToken, Pausable {
    string constant public name = "MOCROW";
    string constant public symbol = "MCW";
    uint256 constant public decimals = 18;

    uint256 constant public RESERVED_TOKENS_FOR_FOUNDERS = 23500000 * (10 ** decimals);
    uint256 constant public RESERVED_TOKENS_FOR_BOUNTY_PROGRAM = 9480000 * (10 ** decimals);
    uint256 constant public RESERVED_TOKENS_FOR_PLATFORM_OPERATIONS = 70588235 * (10 ** decimals);

    uint256 constant public TOTAL_SUPPLY_VALUE = 235294118 * (10 ** decimals);

    address private addressIco;

    modifier onlyIco() {
        require(msg.sender == addressIco);
        _;
    }

    /**
    * @dev Create MocrowCoin contract with reserves.
    * @param _foundersReserve The address of founders reserve.
    * @param _bountyProgramReserve The address of bounty program reserve.
    * @param _platformOperationsReserve The address of platform operations reserve.
    */
    function MocrowCoin(address _foundersReserve, address _bountyProgramReserve, address _platformOperationsReserve) public {
        require(
            _platformOperationsReserve != address(0) && 
            _foundersReserve != address(0) && _bountyProgramReserve != address(0)
        );

        addressIco = msg.sender;

        totalSupply = TOTAL_SUPPLY_VALUE;

        // balances[msg.sender] = TOTAL_SUPPLY_VALUE - RESERVED_TOKENS_FOR_PLATFORM_OPERATIONS - RESERVED_TOKENS_FOR_FOUNDERS - RESERVED_TOKENS_FOR_BOUNTY_PROGRAM;
        balances[msg.sender] = TOTAL_SUPPLY_VALUE.sub(RESERVED_TOKENS_FOR_PLATFORM_OPERATIONS.add(RESERVED_TOKENS_FOR_FOUNDERS).add(RESERVED_TOKENS_FOR_BOUNTY_PROGRAM));

        balances[_platformOperationsReserve] = RESERVED_TOKENS_FOR_PLATFORM_OPERATIONS;
        balances[_foundersReserve] = RESERVED_TOKENS_FOR_FOUNDERS;
        balances[_bountyProgramReserve] = RESERVED_TOKENS_FOR_BOUNTY_PROGRAM;
    }

    /**
    * @dev Transfer token for a specified address with pause and freeze features for owner.
    * @dev Only applies when the transfer is allowed by the owner.
    * @param _to The address to transfer to.
    * @param _value The amount to be transferred.
    */
    function transfer(address _to, uint256 _value) whenNotPaused public returns (bool) {
        require(!isFrozen(msg.sender));
        super.transfer(_to, _value);
    }

    /**
    * @dev Transfer tokens from one address to another with pause and freeze features for owner.
    * @dev Only applies when the transfer is allowed by the owner.
    * @param _from address The address which you want to send tokens from
    * @param _to address The address which you want to transfer to
    * @param _value uint256 the amount of tokens to be transferred
    */
    function transferFrom(address _from, address _to, uint256 _value) whenNotPaused public returns (bool) {
        require(!isFrozen(msg.sender));
        require(!isFrozen(_from));
        super.transferFrom(_from, _to, _value);
    }

    /**
    * @dev Transfer tokens from ICO address to another address.
    * @param _to The address to transfer to.
    * @param _value The amount to be transferred.
    */
    function transferFromIco(address _to, uint256 _value) onlyIco public returns (bool) {
        super.transfer(_to, _value);
    }

    /**
    * Set allowance for other address and notify
    *
    * Allows `_spender` to spend no more than `_value` tokens on your behalf, and then ping the contract about it
    *
    * @param _spender The address authorized to spend
    * @param _value the max amount they can spend
    * @param _extraData some extra information to send to the approved contract
    */
    function approveAndCall(address _spender, uint256 _value, bytes _extraData) public returns (bool success) {
        tokenRecipient spender = tokenRecipient(_spender);
        if (approve(_spender, _value)) {
            spender.receiveApproval(
                msg.sender,
                _value, this,
                _extraData);
            return true;
        }
    }

}
