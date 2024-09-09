# Raspberry Pi Configuration

This page contains the information found in the following tutorial to setup kiosk mode on a Raspberry Pi.

[Tutorial](https://die-antwort.eu/techblog/2017-12-setup-raspberry-pi-for-kiosk-mode/)

## Getting Starting

1. Update all preinstalled packages

```
sudo apt-get update
sudo apt-get upgrade
```

2. Edit config file under `/boot/config.txt` to set display rotation

```txt
# Add the following line at end of file
display_hdmi_rotate=3

# Comment out following line
dtoverlay=vc4-fkms-v3d
```

## Install X Server and Window Manager

`sudo apt-get install --no-install-recommends xserver-xorg x11-xserver-utils xinit openbox`

## Install Chromium

`sudo apt-get install --no-install-recommends chromium-browser`

## Openbox Configuration

Edit `/etc/xdg/openbox/autostart`

```
# Disable any form of screen saver / screen blanking / power management
xset s off
xset s noblank
xset -dpms

# Start Chromium in kiosk mode
sed -i 's/"exited_cleanly":false/"exited_cleanly":true/' ~/.config/chromium/'Local State'
sed -i 's/"exited_cleanly":false/"exited_cleanly":true/; s/"exit_type":"[^"]\+"/"exit_type":"Normal"/' ~/.config/chromium/Default/Preferences
chromium-browser --disable-infobars --kiosk --display=:0 --kiosk --incognito --window-position=0,0 --enable-features=WebUIDarkMode --force-dark-mode https://ghost-hologram-mirror.fly.dev/display
```

### Start X Server

`sudo startx -- -nocursor`

## Automatically Start X Server

Edit `.bash_profile` and append the following line.

```
cd ~

sudo nano .bash_profile

[[ -z $DISPLAY && $XDG_VTNR -eq 1 ]] && startx -- -nocursor

source ~/.bash_profile
```
