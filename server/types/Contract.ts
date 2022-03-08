type AbiItem = any;

interface Contract {
  abi: AbiItem[];
  address: string;
}

export default Contract;

// (property) Eth.Contract: new (jsonInterface: AbiItem | AbiItem[], address?: string | undefined, options?: ContractOptions | undefined) => Contract