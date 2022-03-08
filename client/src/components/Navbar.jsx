import React, { useContext } from "react";
import { TransactionContext } from "../context/TransactionContext";
import { shortenAddress } from "../utils/shortenAddress";
import ethereum from "../../images/ethereum.svg";

const Navbar = () => {

  const { address, balance, handleLogout } = useContext(TransactionContext);

  return (
    <nav className="w-full flex md:justify-center justify-between items-center p-4 pr-1 pl-1" style={{justifyContent: "space-between"}}>
      <div className="md:flex-[0.5] flex-initial justify-center items-center">
        <img src={ethereum} alt="Ethereum" title="Crypto Xchange" className="w-24 cursor-pointer logo" />
      </div>
      <ul className="text-white md:flex list-none flex-row justify-between items-center flex-initial">
        {address ?
          (<div className="flex flex-row items-center text-black">
            <p>{shortenAddress(address)}({ balance } PC)</p>
            <button 
              type="button"
              onClick={ handleLogout }
              className="bg-[#5f9ea0] py-1 px-3 mx-4 rounded-full cursor-pointer hover:bg-[#97cbcd] text-white text-base font-semibold"
            >
              Logout
            </button>
          </div>)
          :
          (<button 
              type="button"
              className="bg-[#5f9ea0] py-1 px-3 mx-4 rounded-full cursor-pointer hover:bg-[#97cbcd] text-white text-base font-semibold"
            >
            Connect
          </button>)
        }
      </ul>
    </nav>
  );
};

export default Navbar;