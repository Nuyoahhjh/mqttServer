openssl req -sha512 -new -x509 -key test.key -out test.crt -days 3650 -subj /CN=www.zerowait.cn

openssl x509 -in test.crt -inform DER -outform PEM -out test.pem

CA根证书的生成步骤
//生成私钥
openssl genrsa -out test.key 2048
4.4 第二步 生成CA证书请求（.csr）
openssl req -new -key test.key -out test.csr
4.5 第三步 自签名得到根证书（.crt）
openssl x509 -req -days 365 -in test.csr -signkey test.key -out test.crt

 4.6 第四步 生成pem格式证书
cat test.crt test.key > test.pem