

# Setting Permissions for Stepper Controller

The following works for my Pololu Tic T834 and will likely work for all Pololu serial controllers.

```
echo 'SUBSYSTEM=="usb", ATTRS{idVendor}=="1ffb", MODE:="0666"' > /etc/udev/rules.d/99-stepper.rules
udevadm control --reload
```