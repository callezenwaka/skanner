import React, { useEffect, useState } from 'react';
import axios from 'axios';

export const TransactionContext = React.createContext();

export const TransactionsProvider = ({ children }) => {
  const [formData, setformData] = useState({
    receiver: '',
    amount: '',
  })
  const [credential, setCredential] = useState('');
  const [address, setAddress] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [balance, setBalance] = useState(false);

  const request = async (url, method, data) => {
    const res = await axios({
      method: method,
      url: `${url}`,
      data: data,
      headers: {
        'content-type': 'application/json',
        Accept: 'application/json',
        // Authorization: `Bearer ${token}`,
      },
    })
    return res.data
  }

  const handleChange = (e, name) => {
    setformData((prevState) => ({ ...prevState, [name]: e.target.value }))
  }

  const handleCredential = (e) => {
    setCredential(e.target.value)
  }

  const handleLogout = async () => {
    localStorage.removeItem('address');
    localStorage.removeItem('credential');

    setAddress('');
    setCredential('');
    window.location.reload();
  }

  const getBalance = async () => {
    try {
      // TODO: Get balance
      const address = localStorage.getItem('address')
      const credential = localStorage.getItem('credential');
      if (!address) return
      const data = await request(`${import.meta.env.VITE_BASE_URL}/coin/${credential}`, 'get', '');

      setBalance(data);
    } catch (error) {
      console.log(error);
    }
  }

  const checkIfWalletIsConnect = async () => {
    try {
      // TODO: Check wallet connection
      const address = localStorage.getItem('address');
      if (!address) return;
      setAddress(address);

      getBalance();
    } catch (error) {
      console.log(error);
    }
  }

  const connectWallet = async () => {
    try {
      // TODO: Connect wallet
      if (!credential) return
      localStorage.setItem('credential', credential);
      const data = await request(`${import.meta.env.VITE_BASE_URL}/coin/login`, 'post', { credential: credential })

      if (typeof data != 'object') return
      localStorage.setItem('address', data.address);

      setAddress(data.address)
      window.location.reload()
    } catch (error) {
      console.log(error)
    }
  }

  const sendTransfer = async () => {
    try {
      // TODO: Send transfer
      const address = localStorage.getItem('address');
      const credential = localStorage.getItem('credential');
      if (!address) return;
      const { receiver, amount } = formData;
      setIsLoading(true);
      await request( `${import.meta.env.VITE_BASE_URL}/coin`, 'post', {
        receiver: receiver,
        amount: amount,
        credential: credential,
      });
      setIsLoading(false);
      getBalance();
      window.location.reload();
    } catch (error) {
      console.log(error);
    }
  }

  useEffect(() => {
    checkIfWalletIsConnect();
  }, []);

  return (
    <TransactionContext.Provider
      value={{
        address,
        isLoading,
        formData,
        balance,
        connectWallet,
        sendTransfer,
        handleCredential,
        handleChange,
        handleLogout,
      }}
    >
      {children}
    </TransactionContext.Provider>
  )
}
