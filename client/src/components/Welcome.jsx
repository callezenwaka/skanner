import React, { useContext } from "react";
import { AiFillPlayCircle } from "react-icons/ai";
import { TransactionContext } from "../context/TransactionContext";
import { Loader } from ".";

const Input = ({ placeholder, name, type, value, handleChange }) => (
  <input
    placeholder={placeholder}
    type={type}
    step="0.0001"
    value={value}
    onChange={(e) => handleChange(e, name)}
    className="my-2 w-full rounded-sm p-2 outline-none bg-transparent text-black border-none text-sm white-glassmorphism"
    style={{outline: "2px solid cadetblue"}}
  />
);

const Welcome = () => {
  const { address, formData, isLoading, connectWallet, sendTransfer, handleChange, handleCredential } = useContext(TransactionContext);

  const handleSubmit = (e) => {
    const { receiver, amount } = formData;
    e.preventDefault();
    if (!receiver || !amount) return;
    sendTransfer();
  };

  return (
    <div className="flex w-full justify-center items-center">
      <div className="flex sm:w-96 w-full mf:flex-row flex-col items-start justify-between md:p-20 py-12 px-4">
        <div className="flex flex-col flex-1 items-center justify-start w-full mf:mt-0 mt-10">
          {address?
            (<div className="p-5 sm:w-96 w-full flex flex-col justify-start items-center blue-glassmorphism">
              <Input placeholder="Receiver address" name="receiver" type="text" handleChange={handleChange} />
              <Input placeholder="Amount (PC)" name="amount" type="number" handleChange={handleChange} />
              <div className="h-[1px] w-full bg-gray-400 my-2" />
              {isLoading
                ? <Loader />
                : (
                  <button
                    type="button"
                    onClick={handleSubmit}
                    className="text-white w-full mt-2 border-[1px] p-2 border-[#3d4f7c] bg-[#5f9ea0] hover:bg-[#97cbcd] rounded-full cursor-pointer"
                  >
                    Transfer Coin
                  </button>
                )}
            </div>)
            : 
            (<div className="sm:w-96 w-full flex flex-col justify-center items-center w-full">
              <input 
                placeholder="Private key"
                name="sender"
                type="text"
                onChange={(e) => handleCredential(e)} 
                className="my-2 w-full rounded-sm p-2 outline bg-transparent text-black border-none text-sm white-glassmorphism" 
                style={{outline: "2px solid lavenderblush"}}
              />
              <button
                type="button"
                onClick={connectWallet}
                className="flex flex-row justify-center items-center my-5 bg-[#5f9ea0] p-3 rounded-full cursor-pointer hover:bg-[#97cbcd]"
              >
                <AiFillPlayCircle className="text-white mr-2" />
                <p className="text-white text-base font-semibold">
                  Connect Wallet
                </p>
              </button>
            </div>)}
        </div>
      </div>
    </div>
  );
};

export default Welcome;