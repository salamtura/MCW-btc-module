pragma solidity ^0.4.15;

// File: contracts/zeppelin-solidity/contracts/ownership/Ownable.sol

/**
 * @title Ownable
 * @dev The Ownable contract has an owner address, and provides basic authorization control
 * functions, this simplifies the implementation of "user permissions".
 */
contract Ownable {
  address public owner;


  event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);


  /**
   * @dev The Ownable constructor sets the original `owner` of the contract to the sender
   * account.
   */
  function Ownable() {
    owner = msg.sender;
  }


  /**
   * @dev Throws if called by any account other than the owner.
   */
  modifier onlyOwner() {
    require(msg.sender == owner);
    _;
  }


  /**
   * @dev Allows the current owner to transfer control of the contract to a newOwner.
   * @param newOwner The address to transfer ownership to.
   */
  function transferOwnership(address newOwner) onlyOwner public {
    require(newOwner != address(0));
    OwnershipTransferred(owner, newOwner);
    owner = newOwner;
  }

}

// File: contracts/zeppelin-solidity/contracts/math/SafeMath.sol

/**
 * @title SafeMath
 * @dev Math operations with safety checks that throw on error
 */
library SafeMath {
  function mul(uint256 a, uint256 b) internal constant returns (uint256) {
    uint256 c = a * b;
    assert(a == 0 || c / a == b);
    return c;
  }

  function div(uint256 a, uint256 b) internal constant returns (uint256) {
    // assert(b > 0); // Solidity automatically throws when dividing by 0
    uint256 c = a / b;
    // assert(a == b * c + a % b); // There is no case in which this doesn't hold
    return c;
  }

  function sub(uint256 a, uint256 b) internal constant returns (uint256) {
    assert(b <= a);
    return a - b;
  }

  function add(uint256 a, uint256 b) internal constant returns (uint256) {
    uint256 c = a + b;
    assert(c >= a);
    return c;
  }
}

// File: contracts/zeppelin-solidity/contracts/token/ERC20Basic.sol

/**
 * @title ERC20Basic
 * @dev Simpler version of ERC20 interface
 * @dev see https://github.com/ethereum/EIPs/issues/179
 */
contract ERC20Basic {
  uint256 public totalSupply;
  function balanceOf(address who) public constant returns (uint256);
  function transfer(address to, uint256 value) public returns (bool);
  event Transfer(address indexed from, address indexed to, uint256 value);
}

// File: contracts/zeppelin-solidity/contracts/token/BasicToken.sol

/**
 * @title Basic token
 * @dev Basic version of StandardToken, with no allowances.
 */
contract BasicToken is ERC20Basic {
  using SafeMath for uint256;

  mapping(address => uint256) balances;

  /**
  * @dev transfer token for a specified address
  * @param _to The address to transfer to.
  * @param _value The amount to be transferred.
  */
  function transfer(address _to, uint256 _value) public returns (bool) {
    require(_to != address(0));

    // SafeMath.sub will throw if there is not enough balance.
    balances[msg.sender] = balances[msg.sender].sub(_value);
    balances[_to] = balances[_to].add(_value);
    Transfer(msg.sender, _to, _value);
    return true;
  }

  /**
  * @dev Gets the balance of the specified address.
  * @param _owner The address to query the the balance of.
  * @return An uint256 representing the amount owned by the passed address.
  */
  function balanceOf(address _owner) public constant returns (uint256 balance) {
    return balances[_owner];
  }

}

// File: contracts/zeppelin-solidity/contracts/token/ERC20.sol

/**
 * @title ERC20 interface
 * @dev see https://github.com/ethereum/EIPs/issues/20
 */
contract ERC20 is ERC20Basic {
  function allowance(address owner, address spender) public constant returns (uint256);
  function transferFrom(address from, address to, uint256 value) public returns (bool);
  function approve(address spender, uint256 value) public returns (bool);
  event Approval(address indexed owner, address indexed spender, uint256 value);
}

// File: contracts/zeppelin-solidity/contracts/token/StandardToken.sol

/**
 * @title Standard ERC20 token
 *
 * @dev Implementation of the basic standard token.
 * @dev https://github.com/ethereum/EIPs/issues/20
 * @dev Based on code by FirstBlood: https://github.com/Firstbloodio/token/blob/master/smart_contract/FirstBloodToken.sol
 */
contract StandardToken is ERC20, BasicToken {

  mapping (address => mapping (address => uint256)) allowed;


  /**
   * @dev Transfer tokens from one address to another
   * @param _from address The address which you want to send tokens from
   * @param _to address The address which you want to transfer to
   * @param _value uint256 the amount of tokens to be transferred
   */
  function transferFrom(address _from, address _to, uint256 _value) public returns (bool) {
    require(_to != address(0));

    uint256 _allowance = allowed[_from][msg.sender];

    // Check is not needed because sub(_allowance, _value) will already throw if this condition is not met
    // require (_value <= _allowance);

    balances[_from] = balances[_from].sub(_value);
    balances[_to] = balances[_to].add(_value);
    allowed[_from][msg.sender] = _allowance.sub(_value);
    Transfer(_from, _to, _value);
    return true;
  }

  /**
   * @dev Approve the passed address to spend the specified amount of tokens on behalf of msg.sender.
   *
   * Beware that changing an allowance with this method brings the risk that someone may use both the old
   * and the new allowance by unfortunate transaction ordering. One possible solution to mitigate this
   * race condition is to first reduce the spender's allowance to 0 and set the desired value afterwards:
   * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
   * @param _spender The address which will spend the funds.
   * @param _value The amount of tokens to be spent.
   */
  function approve(address _spender, uint256 _value) public returns (bool) {
    allowed[msg.sender][_spender] = _value;
    Approval(msg.sender, _spender, _value);
    return true;
  }

  /**
   * @dev Function to check the amount of tokens that an owner allowed to a spender.
   * @param _owner address The address which owns the funds.
   * @param _spender address The address which will spend the funds.
   * @return A uint256 specifying the amount of tokens still available for the spender.
   */
  function allowance(address _owner, address _spender) public constant returns (uint256 remaining) {
    return allowed[_owner][_spender];
  }

  /**
   * approve should be called when allowed[_spender] == 0. To increment
   * allowed value is better to use this function to avoid 2 calls (and wait until
   * the first transaction is mined)
   * From MonolithDAO Token.sol
   */
  function increaseApproval (address _spender, uint _addedValue)
    returns (bool success) {
    allowed[msg.sender][_spender] = allowed[msg.sender][_spender].add(_addedValue);
    Approval(msg.sender, _spender, allowed[msg.sender][_spender]);
    return true;
  }

  function decreaseApproval (address _spender, uint _subtractedValue)
    returns (bool success) {
    uint oldValue = allowed[msg.sender][_spender];
    if (_subtractedValue > oldValue) {
      allowed[msg.sender][_spender] = 0;
    } else {
      allowed[msg.sender][_spender] = oldValue.sub(_subtractedValue);
    }
    Approval(msg.sender, _spender, allowed[msg.sender][_spender]);
    return true;
  }

}

// File: contracts/BurnableToken.sol

/**
* @title Customized Burnable Token
* @dev Token that can be irreversibly burned (destroyed).
*/
contract BurnableToken is StandardToken, Ownable {

    event Burn(address indexed burner, uint256 amount);

    /**
    * @dev Anybody can burn a specific amount of their tokens.
    * @param _amount The amount of token to be burned.
    */
    function burn(uint256 _amount) public {
        require(_amount > 0);
        require(_amount <= balances[msg.sender]);
        // no need to require _amount <= totalSupply, since that would imply the
        // sender's balance is greater than the totalSupply, which *should* be an assertion failure

        address burner = msg.sender;
        balances[burner] = balances[burner].sub(_amount);
        totalSupply = totalSupply.sub(_amount);
        Transfer(burner, address(0), _amount);
        Burn(burner, _amount);
    }
}

// File: contracts/FreezableToken.sol

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
    function freezeAccount(address _wallet) onlyOwner public {
        require(_wallet != address(0));
        frozenList[_wallet] = true;
        FrozenFunds(_wallet, true);
    }

    /**
    * @dev Owner can unfreeze the token balance for chosen token holder.
    * @param _wallet The address of token holder whose tokens to be unfrozen.
    */
    function unfreezeAccount(address _wallet) onlyOwner public {
        require(_wallet != address(0));
        frozenList[_wallet] = false;
        FrozenFunds(_wallet, false);
    }

    /**
    * @dev Check the specified token holder whether his/her token balance is frozen.
    * @param _wallet The address of token holder to check.
    */ 
    function isFrozen(address _wallet) constant public returns (bool) {
        return frozenList[_wallet];
    }

}

// File: contracts/zeppelin-solidity/contracts/lifecycle/Pausable.sol

/**
 * @title Pausable
 * @dev Base contract which allows children to implement an emergency stop mechanism.
 */
contract Pausable is Ownable {
  event Pause();
  event Unpause();

  bool public paused = false;


  /**
   * @dev Modifier to make a function callable only when the contract is not paused.
   */
  modifier whenNotPaused() {
    require(!paused);
    _;
  }

  /**
   * @dev Modifier to make a function callable only when the contract is paused.
   */
  modifier whenPaused() {
    require(paused);
    _;
  }

  /**
   * @dev called by the owner to pause, triggers stopped state
   */
  function pause() onlyOwner whenNotPaused public {
    paused = true;
    Pause();
  }

  /**
   * @dev called by the owner to unpause, returns to normal state
   */
  function unpause() onlyOwner whenPaused public {
    paused = false;
    Unpause();
  }
}

// File: contracts/zeppelin-solidity/contracts/token/MintableToken.sol

/**
 * @title Mintable token
 * @dev Simple ERC20 Token example, with mintable token creation
 * @dev Issue: * https://github.com/OpenZeppelin/zeppelin-solidity/issues/120
 * Based on code by TokenMarketNet: https://github.com/TokenMarketNet/ico/blob/master/contracts/MintableToken.sol
 */

contract MintableToken is StandardToken, Ownable {
  event Mint(address indexed to, uint256 amount);
  event MintFinished();

  bool public mintingFinished = false;


  modifier canMint() {
    require(!mintingFinished);
    _;
  }

  /**
   * @dev Function to mint tokens
   * @param _to The address that will receive the minted tokens.
   * @param _amount The amount of tokens to mint.
   * @return A boolean that indicates if the operation was successful.
   */
  function mint(address _to, uint256 _amount) onlyOwner canMint public returns (bool) {
    totalSupply = totalSupply.add(_amount);
    balances[_to] = balances[_to].add(_amount);
    Mint(_to, _amount);
    Transfer(0x0, _to, _amount);
    return true;
  }

  /**
   * @dev Function to stop minting new tokens.
   * @return True if the operation was successful.
   */
  function finishMinting() onlyOwner public returns (bool) {
    mintingFinished = true;
    MintFinished();
    return true;
  }
}

// File: contracts/MocrowCoin.sol

interface tokenRecipient {
    function receiveApproval(
        address _from,
        uint256 _value,
        address _token,
        bytes _extraData)
    public;
}


contract MocrowCoin is BurnableToken, FreezableToken, Pausable {
    string constant public name = "MOCROW";
    string constant public symbol = "MCW";
    uint256 constant public decimals = 18;

    uint256 constant public RESERVED_TOKENS_FOR_FOUNDERS = 23500000 * (10 ** decimals);
    uint256 constant public RESERVED_TOKENS_FOR_BOUNTY_PROGRAM = 9480000 * (10 ** decimals);
    uint256 constant public RESERVED_TOKENS_FOR_PLATFORM_OPERATIONS = 70588235 * (10 ** decimals);

    uint256 constant public TOTAL_SUPPLY_VALUE = 235294118 * (10 ** decimals);

    /**
    * @dev Create MocrowCoin contract with reserves.
    * @param _foundersReserve The address of founders reserve.
    * @param _bountyProgramReserve The address of bounty program reserve.
    * @param _platformOperationsReserve The address of platform operations reserve.
    */
    function MocrowCoin(address _foundersReserve, address _bountyProgramReserve, address _platformOperationsReserve) public {
        require(
            _platformOperationsReserve != address(0) &&
            _foundersReserve != address(0) &&
            _bountyProgramReserve != address(0)
        );
        totalSupply = TOTAL_SUPPLY_VALUE;
        balances[msg.sender] = TOTAL_SUPPLY_VALUE - RESERVED_TOKENS_FOR_PLATFORM_OPERATIONS - RESERVED_TOKENS_FOR_FOUNDERS - RESERVED_TOKENS_FOR_BOUNTY_PROGRAM;
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

// File: contracts/Administrable.sol

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

// File: contracts/Whitelist.sol

/**
 * @title Whitelist contract
 * @dev Whitelist for wallets.
*/
contract Whitelist is Ownable {
    mapping(address => bool) whitelist;

    uint256 public whitelistLength = 0;

    function Whitelist() {
        owner = msg.sender;
    }
    /**
    * @dev Add wallet to whitelist.
    * @dev Accept request from the owner only.
    * @param _wallet The address of wallet to add.
    */
    function addWallet(address _wallet) onlyOwner public {
        require(_wallet != address(0));
        require(!isWhitelisted(_wallet));
        whitelist[_wallet] = true;
        whitelistLength++;
    }

    /**
    * @dev Remove wallet from whitelist.
    * @dev Accept request from the owner only.
    * @param _wallet The address of whitelisted wallet to remove.
    */  
    function removeWallet(address _wallet) onlyOwner public {
        require(_wallet != address(0));
        require(isWhitelisted(_wallet));
        whitelist[_wallet] = false;
        whitelistLength--;
    }

    /**
    * @dev Check the specified wallet whether it is in the whitelist.
    * @param _wallet The address of wallet to check.
    */
    function isWhitelisted(address _wallet) constant public returns (bool) {
        return whitelist[_wallet];
    }

}

// File: contracts/Whitelistable.sol

contract Whitelistable is Administrable {
    Whitelist public whitelist;

    modifier whenWhitelisted(address _sender) {
        require(isWhitelisted(_sender));
        _;
    }

    /**
    * @dev Constructor for Whitelistable contract.
    */
    function Whitelistable() {
        whitelist = new Whitelist();
    }

    /**
    * @dev Add wallet to whitelist.
    * @dev Accept request from the owner or administrator.
    * @param _wallet The address of wallet to add.
    */
    function addWalletToWhitelist(address _wallet) onlyAdministratorOrOwner {
        whitelist.addWallet(_wallet);
    }

    /**
    * @dev Remove wallet from whitelist.
    * @dev Accept request from the owner or administrator.
    * @param _wallet The address of whitelisted wallet to remove.
    */
    function removeWalletFromWhitelist(address _wallet) onlyAdministratorOrOwner {
        whitelist.removeWallet(_wallet);
    }

    /**
    * @dev Check the specified wallet whether it is in the whitelist.
    * @param _wallet The address of wallet to check.
    */
    function isWhitelisted(address _wallet) constant returns (bool) {
        return whitelist.isWhitelisted(_wallet);
    }
}

// File: contracts/MocrowCoinCrowdsale.sol

contract MocrowCoinCrowdsale is Whitelistable, Pausable {
    using SafeMath for uint256;

    uint256 constant public DECIMALS = 18;

    uint256 constant public HARDCAP_TOKENS_PRE_ICO = 36000000 * (10 ** DECIMALS);
    uint256 constant public HARDCAP_TOKENS_ICO = 84000000 * (10 ** DECIMALS);
    uint256 constant public COMPAIGN_ALLOCATION_AND_BONUSES_TOKENS = 11725883 * (10 ** DECIMALS);

    uint256 constant public TOKEN_RATE_PRE_ICO = 17934;
    uint256 constant public TOKEN_RATE_ICO = 35868;

    uint256 constant public MINIMAL_INVESTMENT = 0.1 ether;
    uint256 constant public MAXIMAL_INVESTMENT = 5 ether;

    uint256 constant public MINIMAL_TEN_PERCENT_BONUS_BY_VALUE = 3.5 ether;
    uint256 constant public MINIMAL_FIVE_PERCENT_BONUS_BY_VALUE = 2.5 ether;

    mapping(address => uint256) public investmentsPreIco;
    address[] private investorsPreIco;

    mapping(address => uint256) public investmentsIco;
    address[] private investorsIco;

    uint256 public preIcoDurationDays = 10;
    uint256 public startTimePreIco;
    uint256 public endTimePreIco;

    uint256 public daysDelayAfterPreIco = 10;

    uint256 public icoDurationDays = 60;
    uint256 public startTimeIco;
    uint256 public endTimeIco;

    uint256 public icoTenPercentBonusEnded;
    uint256 public icoFivePercentBonusEnded;

    address private withdrawalWallet1;
    address private withdrawalWallet2;
    address private withdrawalWallet3;
    address private withdrawalWallet4;

    uint256 public withdrawalWallet1Percent = 50;
    uint256 public withdrawalWallet2Percent = 20;
    uint256 public withdrawalWallet3Percent = 15;
    uint256 public withdrawalWallet4Percent = 15;

    address private addressForCampaignAllocation;
    address private addressForUnsoldTokens;

    uint256 public preIcoTokenRateNegativeDecimals = 8;
    uint256 public preIcoTokenRate = TOKEN_RATE_PRE_ICO;
    uint256 public lastDayChangePreIcoTokenRate = 0;
    uint256 public tokensRemainingPreIco = HARDCAP_TOKENS_PRE_ICO;

    uint256 public icoTokenRateNegativeDecimals = 8;
    uint256 public icoTokenRate = TOKEN_RATE_ICO;
    uint256 public lastDayChangeIcoTokenRate = 0;
    uint256 public tokensRemainingIco = HARDCAP_TOKENS_PRE_ICO + HARDCAP_TOKENS_ICO;

    uint256 public compaignAllocationAndBonusRemainingTokens = COMPAIGN_ALLOCATION_AND_BONUSES_TOKENS;

    MocrowCoin public token;

    function isPreIco() constant public returns(bool) {
        return startTimePreIco < now && now < endTimePreIco;
    }

    function isIco() constant public returns(bool) {
        return startTimeIco < now && now < endTimeIco;
    }

    modifier beforePreIcoSalePeriod() {
        require(now < startTimePreIco);
        _;
    }

    modifier beforeIcoSalePeriod() {
        require(now < startTimeIco);
        _;
    }

    modifier preIcoSalePeriod () {
        require(isPreIco());
        _;
    }

    modifier icoSalePeriod() {
        require(isIco());
        _;
    }

    modifier afterIcoSalePeriod() {
        require(endTimeIco < now);
        _;
    }

    modifier minimalInvestment(uint256 _weiAmount) {
        require(_weiAmount > MINIMAL_INVESTMENT);
        _;
    }

    /**
    * @dev Constructor for MocrowCoinCrowdsale contract.
    * @dev Set the owner who can manage administrators, whitelist and token.
    * @param _withdrawalWallet1 The first withdrawal wallet address.
    * @param _withdrawalWallet2 The second withdrawal wallet address.
    * @param _withdrawalWallet3 The third withdrawal wallet address.
    * @param _withdrawalWallet4 The fourth withdrawal wallet address.
    * @param _addressForPlatformOperations The address to which reserved tokens for platform operations will be transferred.
    * @param _addressForFounders The address to which reserved tokens for founders will be transferred.
    * @param _addressForBountyProgram The address to which reserved tokens for bounty program will be transferred.
    * @param _addressForCampaignAllocation The address to which remaining tokens for campaign allocation will be transferred.
    * @param _addressForUnsoldTokens The address to which unsold tokens will be transferred.
    * @param _startTimePreIco The start time of the pre-ICO and crowdsale in general.
    */
    function MocrowCoinCrowdsale(
        address _withdrawalWallet1,
        address _withdrawalWallet2,
        address _withdrawalWallet3,
        address _withdrawalWallet4,
        address _addressForPlatformOperations,
        address _addressForFounders,
        address _addressForBountyProgram,
        address _addressForCampaignAllocation,
        address _addressForUnsoldTokens,
        uint256 _startTimePreIco
    ) public {
        require(_withdrawalWallet1 != address(0) &&
        _withdrawalWallet2 != address(0) &&
        _withdrawalWallet3 != address(0) &&
        _withdrawalWallet4 != address(0) &&
        _addressForPlatformOperations != address(0) &&
        _addressForFounders != address(0) &&
        _addressForBountyProgram != address(0) &&
        _addressForCampaignAllocation != address(0) &&
        _addressForUnsoldTokens != address(0) &&
        _startTimePreIco > now
        );

        startTimePreIco = _startTimePreIco;
        endTimePreIco = startTimePreIco + (preIcoDurationDays * 1 days);

        startTimeIco = endTimePreIco + (daysDelayAfterPreIco * 1 days);
        endTimeIco = startTimeIco + (icoDurationDays * 1 days);

        icoTenPercentBonusEnded = startTimeIco + (2 days);
        icoFivePercentBonusEnded = icoTenPercentBonusEnded + (3 days);

        withdrawalWallet1 = _withdrawalWallet1;
        withdrawalWallet2 = _withdrawalWallet2;
        withdrawalWallet3 = _withdrawalWallet3;
        withdrawalWallet4 = _withdrawalWallet4;

        token = new MocrowCoin(_addressForFounders, _addressForBountyProgram, _addressForPlatformOperations);
        addressForCampaignAllocation = _addressForCampaignAllocation;
        addressForUnsoldTokens = _addressForUnsoldTokens;
    }


    /**
    * @dev Change pre-ICO start time.
    * @dev Only administrator or owner can change pre-ICO start time and only before pre-ICO period.
    * @dev The end time must be less than start time of ICO.
    * @param _startTimePreIco The start time which must be more than now time.
    */
    function changePreIcoStartTime(uint256 _startTimePreIco) onlyAdministratorOrOwner beforePreIcoSalePeriod public {
        require(now < _startTimePreIco);
        uint256 _endTimePreIco = _startTimePreIco + (preIcoDurationDays * 1 days);
        require(_endTimePreIco < startTimeIco);

        startTimePreIco = _startTimePreIco;
        endTimePreIco = _endTimePreIco;
    }

    /**
    * @dev Change ICO start time.
    * @dev Only administrator or owner can change ICO start time and only before ICO period.
    * @dev The end time must be less than start time of ICO.
    * @param _startTimeIco The start time which must be more than end time of the pre-ICO and more than now time.
    */
    function changeIcoStartTime(uint256 _startTimeIco) onlyAdministratorOrOwner beforeIcoSalePeriod public {
        require(_startTimeIco > now && _startTimeIco > endTimePreIco);

        startTimeIco = _startTimeIco;
        endTimeIco = startTimeIco + (icoDurationDays * 1 days);
    }

    /**
    * @dev Change pre-ICO token rate.
    * @dev Only administrator or owner can change pre-ICO token rate and only once per day.
    * @param _preIcoTokenRate Pre-ICO rate of the token.
    * @param _negativeDecimals Number of decimals after comma.
    */
    function changePreIcoTokenRate(uint256 _preIcoTokenRate, uint256 _negativeDecimals) onlyAdministratorOrOwner public {
        uint256 dayNumber = now / (1 days);
        require(dayNumber != lastDayChangePreIcoTokenRate);

        preIcoTokenRate = _preIcoTokenRate;
        preIcoTokenRateNegativeDecimals = _negativeDecimals;
        lastDayChangePreIcoTokenRate = dayNumber;
    }

    /**
    * @dev Change ICO token rate.
    * @dev Only administrator or owner can change pre-ICO token rate and only once per day.
    * @param _icoTokenRate ICO rate of the token.
    * @param _negativeDecimals Number of decimals after comma.
    */
    function changeIcoTokenRate(uint256 _icoTokenRate, uint256 _negativeDecimals) onlyAdministratorOrOwner public {
        uint256 dayNumber = now / (1 days);
        require(dayNumber != lastDayChangeIcoTokenRate);

        icoTokenRate = _icoTokenRate;
        icoTokenRateNegativeDecimals = _negativeDecimals;
        lastDayChangeIcoTokenRate = dayNumber;
    }

    /**
    * @dev Called by the owner or administrator to pause, triggers stopped state
    */
    function pause() onlyAdministratorOrOwner whenNotPaused public {
        paused = true;
        Pause();
    }

    /**
    * @dev Called by the owner or administrator to unpause, returns to normal state
    */
    function unpause() onlyAdministratorOrOwner whenPaused public {
        paused = false;
        Unpause();
    }

    /**
    * @dev Transfer to withdrawal wallets with considering of percentage.
    */
    function withdrawalWalletsTransfer(uint256 value) private {
        uint256 withdrawalWallet1Value = withdrawalWallet1Percent.mul(value).div(100);
        uint256 withdrawalWallet2Value = withdrawalWallet2Percent.mul(value).div(100);
        uint256 withdrawalWallet3Value = withdrawalWallet3Percent.mul(value).div(100);
        uint256 withdrawalWallet4Value = value.sub(withdrawalWallet1Value.add(withdrawalWallet2Value).add(withdrawalWallet3Value));
        withdrawalWallet1.transfer(withdrawalWallet1Value);
        withdrawalWallet2.transfer(withdrawalWallet2Value);
        withdrawalWallet3.transfer(withdrawalWallet3Value);
        withdrawalWallet4.transfer(withdrawalWallet4Value);
    }

    function transferTokensPreIco(address _walletOwner, uint256 _weiAmount, uint256 _tokensAmount) private {
        tokensRemainingPreIco = tokensRemainingPreIco.sub(_tokensAmount);
        tokensRemainingIco = tokensRemainingIco.sub(_tokensAmount);

        if (investmentsPreIco[_walletOwner] == 0) {
            investorsPreIco.push(_walletOwner);
        }
        investmentsPreIco[_walletOwner] = investmentsPreIco[_walletOwner].add(_weiAmount);

        token.transfer(_walletOwner, _tokensAmount);

    }

    /**
    * @dev Sell tokens during pre-ICO.
    * @dev Sell tokens only for whitelisted wallets if crawdsale is not paused.
    */
    function sellTokensPreIco()
    preIcoSalePeriod
    whenWhitelisted(msg.sender)
    whenNotPaused
    minimalInvestment(msg.value)
    public payable {
        require(tokensRemainingPreIco > 0);
        uint256 weiAmount = msg.value;
        uint256 tokensAmount = weiAmount.div(preIcoTokenRate).mul(10 ** preIcoTokenRateNegativeDecimals);
        if (tokensRemainingPreIco < tokensAmount) {
            uint256 tokensDifferent = tokensAmount.sub(tokensRemainingPreIco);
            uint256 excessiveFunds = tokensDifferent.mul(preIcoTokenRate).div(10 ** preIcoTokenRateNegativeDecimals);

            weiAmount = weiAmount.sub(excessiveFunds);
            tokensAmount = tokensRemainingPreIco;
            msg.sender.transfer(excessiveFunds);
        }
        withdrawalWalletsTransfer(weiAmount);
        transferTokensPreIco(msg.sender, weiAmount, tokensAmount);
    }

    /**
    * @dev Sell tokens during pre-ICO for BTC.
    * @dev Only administrator or owner can sell tokens only for whitelisted wallets if crawdsale is not paused.
    */
    function sellTokensForBTCPreIco(address _wallet, uint256 _weiAmount)
    onlyAdministratorOrOwner
    preIcoSalePeriod
    whenWhitelisted(_wallet)
    whenNotPaused
    minimalInvestment(_weiAmount)
    public {
        uint256 tokensAmount = _weiAmount.div(preIcoTokenRate).mul(10 ** preIcoTokenRateNegativeDecimals);
        require(tokensRemainingPreIco > tokensAmount);
        transferTokensPreIco(_wallet, _weiAmount, tokensAmount);
    }

    /**
    * @dev Sell tokens during ICO.
    * @dev Sell tokens only for whitelisted wallets if crawdsale is not paused.
    */
    function transferTokensIco(address _walletOwner, uint256 _weiAmount, uint256 _tokensAmount) public payable {
        uint256 bonusTokens = 0;

        if (compaignAllocationAndBonusRemainingTokens > 0) {
            uint bonus = 0;
            if (now < icoTenPercentBonusEnded) {
                bonus = bonus + 10;
            } else if (now < icoFivePercentBonusEnded) {
                bonus = bonus + 5;
            }

            if (_weiAmount >= MINIMAL_TEN_PERCENT_BONUS_BY_VALUE) {
                bonus = bonus + 10;
            } else if (_weiAmount >= MINIMAL_FIVE_PERCENT_BONUS_BY_VALUE) {
                bonus = bonus + 5;
            }

            bonusTokens = _tokensAmount.mul(bonus).div(100);

            if (compaignAllocationAndBonusRemainingTokens < bonusTokens) {
                bonusTokens = compaignAllocationAndBonusRemainingTokens;
            }
            compaignAllocationAndBonusRemainingTokens = compaignAllocationAndBonusRemainingTokens.sub(bonusTokens);
        }

        tokensRemainingIco = tokensRemainingIco.sub(_tokensAmount);

        uint256 tokensAmountWithBonuses = _tokensAmount.add(bonusTokens);

        if (investmentsIco[_walletOwner] == 0) {
            investorsIco.push(_walletOwner);
        }

        investmentsIco[_walletOwner] = investmentsIco[_walletOwner].add(_weiAmount);
        token.transfer(_walletOwner, tokensAmountWithBonuses);
    }


    /**
    * @dev Sell tokens during ICO.
    * @dev Sell tokens only for whitelisted wallets if crawdsale is not paused.
    */
    function sellTokensIco()
    icoSalePeriod
    whenWhitelisted(msg.sender)
    whenNotPaused
    minimalInvestment(msg.value)
    public payable {
        require(tokensRemainingIco > 0);
        uint256 excessiveFunds = 0;
        uint256 weiAmount = msg.value;

        if (weiAmount > MAXIMAL_INVESTMENT) {
            weiAmount = MAXIMAL_INVESTMENT;
            excessiveFunds = weiAmount.sub(MAXIMAL_INVESTMENT);
        }

        uint256 tokensAmount = weiAmount.div(icoTokenRate).mul(10 ** icoTokenRateNegativeDecimals);

        if (tokensRemainingIco < tokensAmount) {
            uint256 tokensDifferent = tokensAmount.sub(tokensRemainingIco);
            excessiveFunds = excessiveFunds.add(tokensDifferent.mul(icoTokenRate).div(10 ** icoTokenRateNegativeDecimals));

            weiAmount = weiAmount.sub(excessiveFunds);
            tokensAmount = tokensRemainingIco;
        }
        withdrawalWalletsTransfer(weiAmount);
        if (excessiveFunds > 0) {
            msg.sender.transfer(excessiveFunds);
        }
        transferTokensIco(msg.sender, weiAmount, tokensAmount);
    }

    /**
    * @dev Sell tokens during ICO for BTC.
    * @dev Only administrator or owner can sell tokens only for whitelisted wallets if crawdsale is not paused.
    */
    function sellTokensForBTCIco(address _wallet, uint256 _weiAmount)
    onlyAdministratorOrOwner
    icoSalePeriod
    whenWhitelisted(_wallet)
    whenNotPaused
    minimalInvestment(_weiAmount)
    public {
        uint256 tokensAmount = _weiAmount.div(icoTokenRate).mul(10 ** icoTokenRateNegativeDecimals);
        require(tokensRemainingIco > tokensAmount);
        transferTokensIco(_wallet, _weiAmount, tokensAmount);
    }

    /**
    * @dev Transfer remaining compaign allocation and bonus tokens.
    * @dev Transfer tokens only for administrators or owner and only after ICO period.
    */
    function transferRemainingCompaignAllocationAndBonusTokens() onlyAdministratorOrOwner afterIcoSalePeriod public {
        require(compaignAllocationAndBonusRemainingTokens > 0);
        token.transfer(addressForCampaignAllocation, compaignAllocationAndBonusRemainingTokens);
        compaignAllocationAndBonusRemainingTokens = 0;
        if (compaignAllocationAndBonusRemainingTokens == 0 && tokensRemainingIco == 0) {
            token.transferOwnership(owner);
            whitelist.transferOwnership(owner);
        }
    }

    /**
    * @dev Transfer unsold tokens.
    * @dev Transfer tokens only for administrators or owner and only after ICO period.
    */
    function transferUnsoldTokens() onlyAdministratorOrOwner afterIcoSalePeriod public {
        require(tokensRemainingIco > 0);
        token.transfer(addressForUnsoldTokens, tokensRemainingIco);
        tokensRemainingIco = 0;
        if (compaignAllocationAndBonusRemainingTokens == 0 && tokensRemainingIco == 0) {
            token.transferOwnership(owner);
            whitelist.transferOwnership(owner);
        }
    }
}
