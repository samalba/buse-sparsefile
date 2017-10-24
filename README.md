# Buse-sparsefile

[Buse-go](https://github.com/samalba/buse-go) driver for using a sparsefile.

## How to use?

Load the nbd module and build the code...

```
sudo modprobe nbd
go build
```

You need the rights on /dev/nbd0 in this example.

```
./buse-sparsefile /dev/nbd0 file.sparse 1024
```

This creates a 1GB sparsefile (the physical size should be 1 block only).
You can now format the device...

```
mkfs.ext4 /dev/nbd0
```

## Development

```
vagrant plugin install vagrant-vbguest
vagrant up
```
