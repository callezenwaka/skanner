import 'dotenv/config';
import { Request, Response, NextFunction } from "express";
import Provider from '@truffle/hdwallet-provider';
import Web3 from 'web3';
const Coin = require('../build/contracts/Coin.json');

/**
 * [START LOGIN]
 * @param {object} req Express request context.
 * @param {object} res Express response context.
 * @param {object} next Express next context.
 * @return {object} json items
 * Login
 */
 export const login = async (req: Request, res: Response, next: NextFunction): Promise<any> => {
  try {
    // Todo: get refresh token from the user
    const { credential } = req.body;
    if (!credential) return res.status(400).json("Credential is incorrect!");

    // Todo: create a provider and generate tokens
    const provider = new Provider(`${process.env.PRIVATE_KEY}`, `https://ropsten.infura.io/v3/${process.env.PROJECT_ID}`); 
    const web3 = new Web3(provider);
    const signer = web3.eth.accounts.wallet.add(`0x${credential}`)

    res.status(200).json({ address: signer.address, });
  } catch (error) {
    return res.status(500).json('Internal Server Error!');
  }
}
// [END LOGIN]

/**
 * [START GET BALANCE]
 * @param {object} req Express request context.
 * @param {object} res Express response context.
 * @param {object} next Express next context.
 * @return {object} json items
 * Retrieve items
 */
export const getBalance = async (req: Request, res: Response, next: NextFunction) => {
	try {
		// Todo: create a provider and query for account balance
    const provider = new Provider(`${process.env.PRIVATE_KEY}`, `https://ropsten.infura.io/v3/${process.env.PROJECT_ID}`); 
    const web3 = new Web3(provider);
    const signer = web3.eth.accounts.wallet.add(`0x${req.params.id}`)
    const networkId = await web3.eth.net.getId();
    const coin = new web3.eth.Contract(
      Coin.abi,
      Coin.networks[networkId].address,
      {from: signer.address}
    );
    const balance = await coin.methods.balanceOf(signer.address).call();

		if (typeof balance != 'string') {
			return res.status(200).json('0');
		}
    
		return res.status(200).json(balance);
	} catch (error) {
    console.log(error);
		return res.status(500).json('Internal Server Error!');
	}
}
// [END GET BALANCE]

/**
 * [START TRANSFER]
 * @param {object} req Express request context.
 * @param {object} res Express response context.
 * @param {object} next Express next context.
 * @return {object} json items
 * Retrieve items
 */
export const sendTransfer = async (req: Request, res: Response, next: NextFunction) => {
  try {
    // Todo: create a provider and send transaction
    const { receiver, amount, credential } = req.body;
    const provider = new Provider(`${process.env.PRIVATE_KEY}`, `https://ropsten.infura.io/v3/${process.env.PROJECT_ID}`); 
    const web3 = new Web3(provider);
    const signer = web3.eth.accounts.wallet.add(`0x${credential}`)
    const networkId = await web3.eth.net.getId();
    const coin = new web3.eth.Contract(
      Coin.abi,
      Coin.networks[networkId].address,
      {from: signer.address}
    );

    const response = await coin.methods.transfer(receiver, amount).send({ from: signer.address, gas: 300000, });
    // await response.wait();
    
    // console.log(`Transaction hash: ${response.transactionHash}`);

    return res.status(200).json(response.transactionHash);
  } catch (error) {
    console.log(error);
    return res.status(500).json('Internal Server Error!');
  }
}
// [END TRANSFER]

export default {
	getBalance,
	sendTransfer,
  login,
}