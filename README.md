# ssrf302
# 原理
![原因](https://github.com/HToTH/ssrf302/blob/master/images/reason.png)
服务器后端，控制很严格的时候，只能通过搭建新的服务器，利用location来进一步请求。
# 使用
#### 1.存在漏洞的请求包，把漏洞的地方用[ssrf_data]代替。（注意content-length尽量比原来的多几十个字符）
![例子](https://github.com/HToTH/ssrf302/blob/master/images/request.png)
#### 2.url 为漏洞的请求地址
#### 3.lserver 自己搭建的服务器的地址
#### 4.filename 抓包的包地址
![介绍](https://github.com/HToTH/ssrf302/blob/master/images/instruction.png)
![使用](https://github.com/HToTH/ssrf302/blob/master/images/use.png)
# 案例
![案例](https://github.com/HToTH/ssrf302/blob/master/images/case.png)