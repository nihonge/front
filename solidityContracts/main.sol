// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract PrecompiledContractHandler {
    // 定义预编译合约地址
    address constant ENCRYPT_CONTRACT = 0x0000000000000000000000000000000000000001;
    address constant DECRYPT_CONTRACT = 0x0000000000000000000000000000000000000002;
    address constant COMPUTE_CONTRACT = 0x0000000000000000000000000000000000000003;
    address constant KEY_GENERATOR_CONTRACT = 0x0000000000000000000000000000000000000004;

    struct UserData {
        bytes key;
        bytes encryptedData;
    }

    mapping(address => UserData) private userStorage;

    // 生成密钥并存储
    function generate_key() public {
        bytes memory result = generateKey();

        userStorage[msg.sender].key = result;
    }

    // 接受用户的密钥和明文，进行加密并存储到链上
    function encrypt_upload(bytes memory plainData) public {
        require(userStorage[msg.sender].key.length > 0, "No key found. Please generate a key first.");

        bytes memory key = userStorage[msg.sender].key;

        // 加密操作
        bytes memory encryptedData = encrypt(abi.encodePacked(key,plainData));

        // 存储加密数据
        userStorage[msg.sender].encryptedData = encryptedData;
    }

    // 处理链上对应ID的数据，进行密态数据计算并返回结果
    function data_process() public view returns (bytes memory) {
        require(userStorage[msg.sender].encryptedData.length > 0, "No encrypted data found.");

        bytes memory encryptedData = userStorage[msg.sender].encryptedData;

        // 计算操作
        (bool success, bytes memory result) = COMPUTE_CONTRACT.staticcall(encryptedData);
        require(success, "Computation failed");

        return result;
    }

    //-------------------------------------预编译合约函数-----------------------------------
    // 执行加密操作
    function encrypt(bytes memory data) private view returns (bytes memory) {
        (bool success, bytes memory result) = ENCRYPT_CONTRACT.staticcall(data);
        require(success, "Encryption failed");
        return result;
    }

    // 执行解密操作
    function decrypt(bytes memory data) private view returns (bytes memory) {
        (bool success, bytes memory result) = DECRYPT_CONTRACT.staticcall(data);
        require(success, "Decryption failed");
        return result;
    }

    // 执行计算操作
    function compute(bytes memory data) private view returns (bytes memory) {
        (bool success, bytes memory result) = COMPUTE_CONTRACT.staticcall(data);
        require(success, "Computation failed");
        return result;
    }

    // 生成密钥
    function generateKey() private view returns (bytes memory) {
        (bool success, bytes memory result) = KEY_GENERATOR_CONTRACT.staticcall("");
        require(success, "Key generation failed");
        return result;
    }

    // function callSHA256(bytes memory data) public view returns (bytes32 result) {
    //     // 预编译合约地址为 0x2，直接调用即可
    //     (bool success, bytes memory returnData) = address(0x2).staticcall(data);
    //     require(success, "SHA256 precompile failed");

    //     // 返回的数据是一个32字节的哈希值
    //     result = abi.decode(returnData, (bytes32));
    // }

    // function computeSHA256(bytes memory data) public pure returns (bytes32) {
    //     return sha256(data);
    // }
}
