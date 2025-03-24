contract KYC {
    // 客户信息结构体
    struct Customer {
        address customerAddress; // 用户地址
        bool isVerified;        // 验证状态
        address verifier;       // 验证者地址
        uint256 verificationTime; // 验证时间
        uint256 registrationTime; // 注册时间
    }

    // 映射：地址 -> 客户信息
    mapping(address => Customer) public customers;

    // 事件
    event KYCRegistered(address indexed customer, uint256 timestamp);
    event KYCVerified(address indexed customer, address indexed verifier, uint256 timestamp);

    // 注册
    function register() external {
        require(customers[msg.sender].customerAddress == address(0), "Already registered");
        customers[msg.sender] = Customer(msg.sender, false, address(0), 0, block.timestamp);
        emit KYCRegistered(msg.sender, block.timestamp);
    }

    // 验证KYC（链下控制验证者权限）
    function verifyKYC(address _customer) external {
        Customer storage customer = customers[_customer];
        require(customer.customerAddress != address(0), "Customer not found");
        require(!customer.isVerified, "Already verified");
        customer.isVerified = true;
        customer.verifier = msg.sender;
        customer.verificationTime = block.timestamp;
        emit KYCVerified(_customer, msg.sender, block.timestamp);
    }

    // 查询状态
    function getKYCStatus(address _customer) external view returns (bool, uint256, address) {
        Customer memory customer = customers[_customer];
        return (customer.isVerified, customer.verificationTime, customer.verifier);
    }
}