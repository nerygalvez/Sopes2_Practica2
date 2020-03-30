obj-m += memo_201403525.o
obj-m += cpu_201403525.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

ejecutar_modulos:
	dmesg -C
	insmod memo_201403525.ko
	rmmod memo_201403525.ko
	insmod cpu_201403525.ko
	rmmod cpu_201403525.ko
	dmesg

