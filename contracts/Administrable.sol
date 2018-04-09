pragma solidity ^0.4.15;

import './zeppelin-solidity/contracts/ownership/Ownable.sol';


/**
 * @title Administrable
 * @dev The Administrable contract has an owner and administrators addresses
 */
contract Administrable is Ownable {
    mapping(address => bool) internal administrators;
    uint256 public administratorsLength = 0;
    /**
     * @dev The Administrable constructor sets the original `owner` of the contract to the sender
     * account and 3 admins.
     */
    function Administrable() {
        owner = msg.sender;
    }

    /**
     * @dev Throws if called by any account other than the owner or administrator.
     */
    modifier onlyAdministratorOrOwner() {
        require(msg.sender == owner || administrators[msg.sender]);
        _;
    }

    function addAdministrator(address _admin) onlyOwner public {
        require(administratorsLength <= 3);
        require(!administrators[_admin]);
        require(_admin != address(0) && _admin != owner);
        administrators[_admin] = true;
        administratorsLength++;
    }

    function removeAdministrator(address _admin) onlyOwner public {
        require(_admin != address(0));
        require(administrators[_admin]);
        administrators[_admin] = false;
        administratorsLength--;
    }

    function isAdministrator(address _admin) public constant returns (bool) {
        return administrators[_admin];
    }
}
