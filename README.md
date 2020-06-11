

## Approximate Bill of Materials
- 1x #3140 Tic 36v4 USB Multi-Interface High-Power Stepper Motor Controller (Connectors Soldered)
- 1x #1478 Stepper Motor: Bipolar, 200 Steps/Rev, 57×76mm, 3.2V, 2.8 A/Phase
- 1x #1462 Wall Power Adapter: 5VDC, 5A, 5.5×2.1mm Barrel Jack, Center-Positive 
- 1x #1993 Pololu Universal Aluminum Mounting Hub for 1/4″ (6.35mm) Shaft, #4-40 Holes (2-Pack)
- 1x #2258 Steel L-Bracket for NEMA 23 Stepper Motors 
- 1x Pi Zero Wireless

### Install ticcmd (this also adds udev rules for the stepper controller)

```bash
wget https://www.pololu.com/file/0J1349/pololu-tic-1.8.0-linux-rpi.tar.xz
tar -xvf pololu-tic-1.8.0-linux-rpi.tar.xz
cd pololu-tic-1.8.0-linux-rpi/
sudo ./install.sh
udevadm control --reload
```

### Install antenna-switch
```bash
# pull the code
git clone https://github.com/tzneal/antenna-switch.git

# build it
cd antenna-switch/
go build github.com/tzneal/antenna-switch/cmd/antenna-switch

# install the binary
sudo mv antenna-switch /usr/bin 

# setcap so it can bind to privileged ports (<1024)
sudo setcap cap_net_bind_service=+epi /usr/bin/antenna-switch
 
# install the service
cd tools
./install-service.sh
```

### Configuration

```bash
vi ~/.antenna-switch.json 
sudo systemctl restart antenna-switch
```
