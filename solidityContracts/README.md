# 业务合约 main.sol
负责将用户在前端的行为进行处理，与链交互，假设预编译合约已部署在链客户端，地址被分配好。
## 预编译合约对应函数
encrypt,decrypt,compute,key_generator,分别对应相应预编译合约
## 业务函数
主要负责一种功能
generate_key:给用户创建新密钥并存储
encrypt_upload:接受用户的密钥和明文，加密数据发到链上
data_process:处理链上对应ID的数据，进行密态数据计算并返回结果