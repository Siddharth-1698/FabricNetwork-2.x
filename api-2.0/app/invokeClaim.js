const {
    Gateway,
    Wallets,
    TxEventHandler,
    GatewayOptions,
    DefaultEventHandlerStrategies,
    TxEventHandlerFactory,
  } = require("fabric-network");
  const fs = require("fs");
  const EventStrategies = require("fabric-network/lib/impl/event/defaulteventhandlerstrategies");
  const path = require("path");
  const log4js = require("log4js");
  const logger = log4js.getLogger("BasicNetwork");
  const util = require("util");
  
  const helper = require("./helper");
  const { blockListener, contractListener } = require("./Listeners");
  
  const invokeClaimTransaction = async (
    channelName,
    chaincodeName,
    fcn,
    args,
    username,
    org_name,
    transientData
  ) => {
    try {
      const ccp = await helper.getCCP(org_name);
      console.log(
        "==================",
        channelName,
        chaincodeName,
        fcn,
        args,
        username,
        org_name
      );
  
      const walletPath = await helper.getWalletPath(org_name);
      const wallet = await Wallets.newFileSystemWallet(walletPath);
      console.log(`Wallet path: ${walletPath}`);
  
      let identity = await wallet.get(username);
      if (!identity) {
        console.log(
          `An identity for the user ${username} does not exist in the wallet, so registering user`
        );
        await helper.getRegisteredUser(username, org_name, true);
        identity = await wallet.get(username);
        console.log("Run the registerUser.js application before retrying");
        return;
      }
  
      const connectOptions = {
        wallet,
        identity: username,
        discovery: { enabled: true, asLocalhost: true },
        // eventHandlerOptions: EventStrategies.NONE
      };
  
      const gateway = new Gateway();
      await gateway.connect(ccp, {
        wallet,
        identity: username,
        discovery: { enabled: true, asLocalhost: true },
      });
  
      const network = await gateway.getNetwork(channelName);
      const contract = network.getContract(chaincodeName);
  
      // Important: Please dont set listener here, I just showed how to set it. If we are doing here, it will set on every invoke call.
      // Instead create separate function and call it once server started, it will keep listening.
      // await contract.addContractListener(contractListener);
      // await network.addBlockListener(blockListener);
  
      // Multiple smartcontract in one chaincode
      let result;
      let message;
      console.log("I am in helper function, Here is the arguments: ", args[0]);
      new_arg =JSON.parse(args[0]) 
      let claimId = "claim/"+new_arg.fhir_id
      console.log(claimId,new_arg.status)

      result = await contract.submitTransaction("UpdateClaimStatus",claimId,new_arg.status );
      console.log(result);
      result = { txid: result.toString() };
   
      result = { txid: result.toString() };
  
      
  
      await gateway.disconnect();
  
      result = JSON.parse(result.toString());
  
      let response = {
        message: message,
        result,
      };
  
      return response;
    } catch (error) {
      console.log(`Getting error: ${error}`);
      return error.message;
    }
  };
  
  exports.invokeClaimTransaction = invokeClaimTransaction;
  