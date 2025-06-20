## 在windows 下生成自签名证书，步骤如下：
 1. 生成RSA私钥（2048位）
```sh
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out lot.cn.key
```

 2. 生成自签名证书
```sh
openssl req -new -x509 -days 365 -key lot.cn.key -out lot.cn.crt -config lot.cn.cnf

```
> OpenSSL 的 genpkey 命令严格要求参数顺序：
> -algorithm RSA：指定算法类型
> -pkeyopt rsa_keygen_bits:2048：设置 RSA 密钥长度为 2048 位
> -out lot.cn.key：指定输出文件名

## 检查私钥是否有效
```sh
openssl rsa -check -in lot.cn.key

```

## 查看证书信息
```sh
openssl x509 -noout -text -in lot.cn.crt
```

## 验证证书中的 SANs

使用以下命令检查证书是否包含正确的域名：
```sh
openssl x509 -noout -text -in lot.cn.crt | Select-String -Pattern "Subject Alternative Name" -Context 0,10
```


## 配置本地域名解析

1. 修改系统的 hosts 文件，将lot.cn指向本地 IP：
2. 用管理员权限打开记事本。
3. 打开C:\Windows\System32\drivers\etc\hosts文件。
 在文件末尾添加如下内容：
```text
127.0.0.1 lot.cn
```

## 信任自签名证书
    为了让系统信任这个自签名证书，你要将其添加到受信任的根证书颁发机构中：
1. 打开证书文件（双击lot.cn.crt）。
2. 点击 “安装证书”。
3. 选择 “本地计算机”，然后点击 “下一步”。
4. 选择 “将所有证书放入下列存储”，再点击 “浏览”。
5. 选择 “受信任的根证书颁发机构”，点击 “确定”。
6. 按照向导完成证书安装。

## 清除浏览器缓存

浏览器可能缓存了旧证书，尝试：
在 Chrome 中访问 chrome://settings/clearBrowserData，选择 “所有时间” 并清除缓存和证书数据。
在 Firefox 中访问 about:preferences#privacy，清除 “缓存的 Web 内容”。

## 卸载已经安装的自签名证书

1. 按下 “Win+R” 组合键打开运行对话框，输入 “certmgr.msc” 并按回车键。
2. 选择证书存储位置：在证书管理器中，根据自签名证书的类型和安装位置，选择相应的文件夹。常见的存储位置有 “个人”“受信任的根证书颁发机构”“受信任的发布者” 等。如果是为本地开发生成的自签名证书，通常可以在 “个人” 文件夹中找到。
3. 选择并卸载证书：在所选的文件夹中找到并展开相应的证书存储，找到要删除的自签名证书。右键单击该证书，选择 “所有任务”，然后在弹出的菜单中选择 “删除”。
4. 确认删除：在弹出的确认对话框中，点击 “是”，确认删除操作。
