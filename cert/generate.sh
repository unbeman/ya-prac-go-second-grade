openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout RootCA.key -out RootCA.pem -subj "/C=RU/CN=Example-Root-CA"
openssl x509 -outform pem -in RootCA.pem -out RootCA.crt
openssl req -new -nodes -newkey rsa:2048 -keyout server.key -out server.csr -subj "/C=RU/ST=Russia/L=Moscow/O=Example-Certificates/CN=localhost"
openssl x509 -req -sha256 -days 1024 -in server.csr -CA RootCA.pem -CAkey RootCA.key -CAcreateserial -extfile domains.ext -out server.crt

openssl genrsa -out jwt_key.pem 1024