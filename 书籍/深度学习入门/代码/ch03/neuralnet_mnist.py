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

# accuracy_cnt = 0
# for i in range(len(x)):
#     y = predict(network, x[i])
#     p= np.argmax(y) # 获取概率最高的元素的索引
#     if p == t[i]:
#         accuracy_cnt += 1
#
# print("Accuracy:" + str(float(accuracy_cnt) / len(x)))


# x 代表输入数据，是测试集中的图片数据。每张图片是28x28像素的手写数字图像，在代码中被展平为一个784维的向量（28 * 28 = 784），所以 x 的形状是 (num_samples, 784)，其中 num_samples 是测试集中的样本数。每个样本是一张28x28的图片，经过展平处理，变成一个784维的向量。


# 只处理第一张图片
x_batch = x[0:1]  # 获取第一张图片（只取第一个样本）
y_batch = predict(network, x_batch)  # 使用神经网络进行预测
p = np.argmax(y_batch, axis=1)  # 找到预测结果中最大值的索引，表示预测的类别
print(f"Predicted label: {p[0]}")  # 输出预测标签
print(f"result label: {t[0]}")  # 输出预测标签
