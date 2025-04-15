# coding: utf-8
import sys, os
sys.path.append(os.pardir)  # 为了导入父目录的文件而进行的设定
import numpy as np
import pickle
from dataset.mnist import load_mnist
from common.functions import sigmoid, softmax


def get_data():
    (x_train, t_train), (x_test, t_test) = load_mnist(normalize=True, flatten=True, one_hot_label=False)
    return x_test, t_test


def init_network():
    with open("sample_weight.pkl", 'rb') as f:
        network = pickle.load(f)
    return network

# 输入层接收展平后的图片向量（784维）。
# 隐藏层通过神经元的加权求和和激活函数（这里是sigmoid）计算出一层层的输出（特征提取）。
# 输出层通过softmax激活函数计算每个数字类别（0-9）对应的概率。
def predict(network, x):
    w1, w2, w3 = network['W1'], network['W2'], network['W3']
    b1, b2, b3 = network['b1'], network['b2'], network['b3']

    a1 = np.dot(x, w1) + b1
    print(f"First layer weighted sum (a1): {a1}")  # 打印第一层的加权和a1
    z1 = sigmoid(a1)
    print(f"First layer output after sigmoid (z1): {z1}")  # 打印第一层输出z1
    a2 = np.dot(z1, w2) + b2
    print(f"Second layer weighted sum (a2): {a2}")  # 打印第二层的加权和a2
    z2 = sigmoid(a2)
    print(f"Second layer output after sigmoid (z2): {z2}")  # 打印第二层输出z2
    a3 = np.dot(z2, w3) + b3
    print(f"Third layer weighted sum (a3): {a3}")  # 打印第三层的加权和a3
    y = softmax(a3)
    print(f"Output probabilities after softmax (y): {y}")  # 打印最终的概率输出

    return y


x, t = get_data()
network = init_network()



batch_size = 100 # 批数量
accuracy_cnt = 0

for i in range(0, len(x), batch_size):
    x_batch = x[i:i+batch_size]
    y_batch = predict(network, x_batch)
    p = np.argmax(y_batch, axis=1)
    accuracy_cnt += np.sum(p == t[i:i+batch_size])

print("Accuracy:" + str(float(accuracy_cnt) / len(x)))
