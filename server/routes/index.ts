'use strict';

// import packages and dependencies
import coin from "../controllers/index";
import express from "express";
const router = express();

router.post('/', coin.sendTransfer);

router.post('/login', coin.login);

router.get('/:id', coin.getBalance);
 
export default router;