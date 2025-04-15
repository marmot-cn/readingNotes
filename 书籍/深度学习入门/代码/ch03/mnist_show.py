# coding: utf-8
import sys, os
sys.path.append(os.pardir)  # 为了导入父目录的文件而进行的设定
import numpy as np
from dataset.mnist import load_mnist
from PIL import Image


def img_show(img):
    # np.uint8(img)：把 img 转换成 uint8 类型（0-255）
    pil_img = Image.fromarray(np.uint8(img))
    pil_img.show()

(x_train, t_train), (x_test, t_test) = load_mnist(flatten=True, normalize=False)

# 获取第一张图片
img = x_train[0]
# 获取第一个标签
label = t_train[0]
print(label)  # 5


# 这时 img 是一个 一维数组，里面存储的是 28×28 个像素点的灰度值
print(img.shape)  # (784,)

# 把一维数据恢复成二维，img 就变成了一个 28×28 的二维数组，每个值对应一个像素点
img = img.reshape(28, 28)  # 把图像的形状变为原来的尺寸
print(img.shape)  # (28, 28)

img_show(img)
