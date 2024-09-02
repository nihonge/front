// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract PrecompiledContractHandler {
    // 定义预编译合约地址
    address constant TEST_CONTRACT = 0x0000000000000000000000000000000000000021;
    address constant COMPUTE_CONTRACT = 0x0000000000000000000000000000000000000022;

    function callCompute(bytes memory input) public view returns (bytes memory) {
        (bool success, bytes memory output) = COMPUTE_CONTRACT.staticcall(input);
        require(success, "Call to precompiled contract failed");
        return output;
    }
}
