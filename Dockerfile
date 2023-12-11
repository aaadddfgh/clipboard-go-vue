# config.js should be mount to '/root' which at the same dirctory as 'clipboard.out' 

FROM alpine:3.18.4

WORKDIR /root

# go build output files at './release/linux-x64'  
ADD "./release/linux-x64" "/root" 



CMD ["./clipboard.out"]

EXPOSE 3000