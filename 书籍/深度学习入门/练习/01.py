import numpy as np
import matplotlib.pyplot as plt


# 权重矩阵的列数 = 该层神经元的数量
# 权重矩阵的行数 = 前一层的输出特征
# 矩阵乘法要求：前一个矩阵的列数 = 后一个矩阵的行数

def sigmoid(x):
    """S型函数：把任意数值压缩到0-1之间，像'开关的灵敏度'"""
    return 1 / (1 + np.exp(-x))

def softmax(x):
    """Softmax：把数值变成概率，像'投票计票器'"""
    exp_x = np.exp(x - np.max(x))
    return exp_x / np.sum(exp_x)

# 输入数据 - 披萨的宽度和高度（已标准化）
pizza = np.array([0.6, 0.9])  # [宽度, 高度]

# 第一层：两个美食顾问的分析
w1 = np.array([
    # [顾问A, 顾问B] ← 列代表神经元
    [0.1,  -0.2],  # 宽度特征对两个顾问的影响
    [0.4,   0.3]   # 高度特征对两个顾问的影响
])
b1 = np.array([0.0, 0.1])      # 顾问们的个人偏好

# 第二层：主厨的最终决策
# w2 = [
#       [美式→顾问A权重, 意式→顾问A权重],  
#       [美式→顾问B权重, 意式→顾问B权重]
#      ]
w2 = np.array([[0.2, 0.3],     # 主厨对顾问A意见的重视程度
               [-0.4, 0.1]])   # 主厨对顾问B意见的重视程度
b2 = np.array([0.05, -0.05])   # 主厨的个人偏好

# --- 前向传播 ---
print("=== 披萨分类过程 ===")
print(f"输入: 宽度={pizza[0]:.1f}, 高度={pizza[1]:.1f}")

# 第一层计算
# 顾问A的信号：0.6*0.1 + 0.9*0.4 = 0.06 + 0.36 = 0.42
# 顾问B的信号：0.6*(-0.2) + 0.9*0.3 = -0.12 + 0.27 = 0.15
# 加上偏置 + b1
# advisor_opinions = [
#     0.42 (顾问A信号) + 0.0 (b1[0]) = 0.42,
#     0.15 (顾问B信号) + 0.1 (b1[1]) = 0.25
# ]

advisor_opinions = np.dot(pizza, w1) + b1
print("\n顾问们原始意见:", advisor_opinions)

activated_opinions = sigmoid(advisor_opinions)
print("激活后意见(0-1):", activated_opinions)

# 输出层计算
chef_input = np.dot(activated_opinions, w2) + b2
print("\n主厨接收到的信号:", chef_input)

probabilities = softmax(chef_input)
print("最终概率:", probabilities)

# 结果可视化
# plt.figure(figsize=(10, 3))
# plt.subplot(1, 3, 1)
# plt.bar(['宽度', '高度'], pizza)
# plt.title("输入特征")

# plt.subplot(1, 3, 2)
# plt.bar(['顾问A', '顾问B'], activated_opinions)
# plt.title("隐藏层输出")

# plt.subplot(1, 3, 3)
# plt.bar(['美式(3)', '意式(7)'], probabilities)
# plt.title("预测概率")
# plt.show()


## 输出
# === 披萨分类过程 ===
# 输入: 宽度=0.6, 高度=0.9

# 顾问们原始意见: [0.42 0.25]
# 激活后意见(0-1): [0.60348325 0.5621765 ]

# 主厨接收到的信号: [-0.05417395  0.18726263]
# 最终概率: [0.43993236 0.56006764]