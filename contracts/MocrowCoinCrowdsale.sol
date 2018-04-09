pragma solidity ^0.4.15;

import './zeppelin-solidity/contracts/lifecycle/Pausable.sol';

import "./MocrowCoin.sol";
import "./Whitelistable.sol";


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
        uint256 _startTimePreIco) public 
        {
        require(_withdrawalWallet1 != address(0) && _withdrawalWallet2 != address(0) && _withdrawalWallet3 != address(0) && _withdrawalWallet4 != address(0));
        require(_addressForPlatformOperations != address(0) && _addressForFounders != address(0) && _addressForBountyProgram != address(0)); 
        require(_addressForCampaignAllocation != address(0) && _addressForUnsoldTokens != address(0) && _startTimePreIco > now);

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
        token.pause();

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

    /**
    * @dev Fallback function can be used to buy tokens.
    */
    function() public payable {
        if (isPreIco()) {
            sellTokensPreIco();
        } else if (isIco()) {
            sellTokensIco();
        } else {
            revert();
        }
    }

    function transferTokensPreIco(address _walletOwner, uint256 _weiAmount, uint256 _tokensAmount) private {
        tokensRemainingPreIco = tokensRemainingPreIco.sub(_tokensAmount);
        tokensRemainingIco = tokensRemainingIco.sub(_tokensAmount);

        if (investmentsPreIco[_walletOwner] == 0) {
            investorsPreIco.push(_walletOwner);
        }
        investmentsPreIco[_walletOwner] = investmentsPreIco[_walletOwner].add(_weiAmount);

        token.transferFromIco(_walletOwner, _tokensAmount);

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
    public payable 
    {
        require(tokensRemainingPreIco > 0);
        uint256 excessiveFunds = 0;
        uint256 weiAmount = msg.value;

        if (weiAmount > MAXIMAL_INVESTMENT) {
            weiAmount = MAXIMAL_INVESTMENT;
            excessiveFunds = weiAmount.sub(MAXIMAL_INVESTMENT);
        }

        uint256 tokensAmount = weiAmount.div(preIcoTokenRate).mul(10 ** preIcoTokenRateNegativeDecimals);

        if (tokensRemainingPreIco < tokensAmount) {
            uint256 tokensDifferent = tokensAmount.sub(tokensRemainingPreIco);
            excessiveFunds = excessiveFunds.add(tokensDifferent.mul(preIcoTokenRate).div(10 ** preIcoTokenRateNegativeDecimals));

            weiAmount = weiAmount.sub(excessiveFunds);
            tokensAmount = tokensRemainingPreIco;
        }

        withdrawalWalletsTransfer(weiAmount);

        if (excessiveFunds > 0) {
            msg.sender.transfer(excessiveFunds);
        }

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
    public 
    {
        uint256 tokensAmount = _weiAmount.div(preIcoTokenRate).mul(10 ** preIcoTokenRateNegativeDecimals);
        require(tokensRemainingPreIco > tokensAmount);
        transferTokensPreIco(_wallet, _weiAmount, tokensAmount);
    }

    /**
    * @dev Sell tokens during ICO.
    * @dev Sell tokens only for whitelisted wallets if crawdsale is not paused.
    */
    function transferTokensIco(address _walletOwner, uint256 _weiAmount, uint256 _tokensAmount) private {
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
        token.transferFromIco(_walletOwner, tokensAmountWithBonuses);
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
    public payable 
    {
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
    public 
    {
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
        token.transferFromIco(addressForCampaignAllocation, compaignAllocationAndBonusRemainingTokens);
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
        token.transferFromIco(addressForUnsoldTokens, tokensRemainingIco);
        tokensRemainingIco = 0;
        if (compaignAllocationAndBonusRemainingTokens == 0 && tokensRemainingIco == 0) {
            token.transferOwnership(owner);
            whitelist.transferOwnership(owner);
        }
    }
}