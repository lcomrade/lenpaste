# Build
## Build from source code
On Debian/Ubuntu:
```
sudo apt update
sudo apt -y install git make gcc golang
git clone https://git.lcomrade.su/root/lenpaste.git
cd ./lenpaste/
git checkout vX.X
make
```

You can find the result of the build in the `./dist/` directory.



## Build Docker image
**Why is it necessary?**
An official image may not support your architecture e.g. MIPS, PowerPC, etc.
So you can build your own image to run on an officially unsupported architecture
(of course you have to rebuild it every time you update Lenpaste).

On Debian/Ubuntu:
```
sudo apt update
sudo apt -y install git docker.io
git clone https://git.lcomrade.su/root/lenpaste.git
cd ./lenpaste/
git checkout vX.X
sudo docker build -t localhost/lenpaste:latest ./
```

The `localhost/lenpaste:latest` image should now have appeared on your local machine.
You can use it in `docker-compose.yml` or copy it to another machine.
