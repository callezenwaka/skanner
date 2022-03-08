const Coin = artifacts.require('./Coin.sol');

contract('Coin', (accounts) => {
  let coin, owner, addr1, addr2;

  before(async () => {
    coin = await Coin.deployed();
    [owner, addr1, addr2, _] = await web3.eth.getAccounts();
  });

  it('deploys successfully', async () => {
    const address = await coin.address;
    assert.notEqual(address, 0x0);
    assert.notEqual(address, '');
    assert.notEqual(address, null);
    assert.notEqual(address, undefined);
  });

  it("Should set the right owner", async () => {
    assert.equal(await coin.owner(), owner);
  });

  it("Should assign the total supply of tokens to the owner", async () => {
    assert.equal(await coin.totalSupply().toString(), await coin.balanceOf(owner).toString());
  });

  it("Should transfer tokens between accounts", async () => {
    await coin.transfer(addr1, 50)
    const addr1Balance = await coin.balanceOf(addr1)
    assert.equal(addr1Balance, 50);
  });

  it("Should fail if not enough tokens", async () => {
    coin.contract.options.from = addr1;
    const initialOwnerBalance = await coin.balanceOf(owner);
    assert(coin.transfer(owner, 1), "Not enough coins");

    const finalOwnerBalance = await coin.balanceOf(owner)
    assert.equal(finalOwnerBalance.toString(), initialOwnerBalance.toString());
  });

  it('Should update balances after transfers', async () => {
    const initialOwnerBalance = await coin.balanceOf(owner);

    await coin.transfer(addr1, 100);
    await coin.transfer(addr2, 50);

    const finalOwnerBalance = await coin.balanceOf(owner);
    assert.equal(finalOwnerBalance, initialOwnerBalance - 150);

    const addr1Balance = await coin.balanceOf(addr1);
    assert.equal(addr1Balance.toString(), 150);

    const addr2Balance = await coin.balanceOf(addr2);
    assert.equal(addr2Balance.toString(), 50);
  });
});