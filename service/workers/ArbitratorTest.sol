pragma solidity ^0.5.0;

import "./ConcurrentHashMap.sol";

contract System {
    function getPid() external returns(uint256);
    function revertPid(uint256 pid) external;
    function createDefer(string calldata id, string calldata signature) external;
    function callDefer(string calldata id) external;
}

contract Queue {
    enum DataType {INVALID, ADDRESS, UINT256, BYTES}
    function create(string calldata id, uint256 elemType) external;
    function size(string calldata id) external returns(uint256);
    function pushUint256(string calldata id, uint256 value) external;
    function popUint256(string calldata id) external returns(uint256);
}

contract ArbitratorTest {
    constructor() public {
        ConcurrentHashMap hashmap = ConcurrentHashMap(0x81);
        hashmap.create("hashmap1", int32(ConcurrentHashMap.DataType.UINT256), int32(ConcurrentHashMap.DataType.UINT256));
        hashmap.create("hashmap2", int32(ConcurrentHashMap.DataType.UINT256), int32(ConcurrentHashMap.DataType.UINT256));
        
        System sys = System(0xa1);
        sys.createDefer("defer1", "DeferFunc1(string)");
        sys.createDefer("defer2", "DeferFunc1(string)");
        
        Queue queue = Queue(0x82);
        queue.create("queue1", uint256(Queue.DataType.UINT256));
        queue.create("queue2", uint256(Queue.DataType.UINT256));
    }
    
    function ParallelFunc1(uint256 key, uint256 value) public payable {
        ConcurrentHashMap hashmap = ConcurrentHashMap(0x81);
        hashmap.set("hashmap1", key, value);
    }
    
    function ParallelFunc2(uint256 value) public payable {
        System sys = System(0xa1);
        uint256 pid = sys.getPid();
        
        ConcurrentHashMap hashmap = ConcurrentHashMap(0x81);
        hashmap.set("hashmap1", pid, value);
        
        Queue queue = Queue(0x82);
        queue.pushUint256("queue1", pid);
        
        sys.callDefer("defer1");
    }
    
    function ParallelFunc3(uint256 key, uint256 value) public payable {
        ConcurrentHashMap hashmap = ConcurrentHashMap(0x81);
        hashmap.set("hashmap1", key, value);
        
        Queue queue = Queue(0x82);
        queue.pushUint256("queue2", key);
        
        System sys = System(0xa1);
        sys.callDefer("defer2");
    }
    
    function DeferFunc1(string memory id) public {
        Queue queue = Queue(0x82);
        if (hashCompareWithLengthCheck(id, "defer1")) {
            System sys = System(0xa1);
            uint256 size = queue.size("queue1");
            for (uint256 i = 0; i < size; i++) {
                uint256 pid = queue.popUint256("queue1");
                if (i % 2 != 0) {
                    sys.revertPid(pid);
                }
            }
        } else if (hashCompareWithLengthCheck(id, "defer2")) {
            ConcurrentHashMap hashmap = ConcurrentHashMap(0x81);
            uint256 size = queue.size("queue2");
            for (uint256 i = 0; i < size; i++) {
                uint256 key = queue.popUint256("queue2");
                uint256 pid = hashmap.getUint256("hashmap1", key);
                hashmap.set("hashmap1", pid, 0);
            }            
        }
    }
    
    function hashCompareWithLengthCheck(string memory a, string memory b) internal pure returns (bool) {
        if(bytes(a).length != bytes(b).length) {
            return false;
        } else {
            return keccak256(bytes(a)) == keccak256(bytes(b));
        }
    }
}