# NoteForPhicommN1Armbian

1 create a usb boot with '20.02.0-rc1.037_Aml-s9xxx_bionic_current_5.5.0-rc6_desktop_20200205.img' by 'balenaEtcher-Portable-1.5.115'

2 change  'tb_name=/dtb/meson-gxl-s905d-phicomm-n1-xiangsm.dtb' of 'uEnv.ini'

3 insert usb and power on Phicomm-N1.

4 username:root

5 password:1234

6 run 'install-aml.sh' under 'root' user

7 change 'source.list' to 'mirrors' site

8 change locale and add language support.

//


https://askubuntu.com/questions/657160/automatically-run-a-command-script-at-startup

rc.local has solved this for me, and although I have read is outdated But kept around for compatibility, if somebody has the updated, but JUST AS SIMPLE answer, please share as I will update my server.

sudo nano /etc/rc.local
You can use either calling a script in rc.local or directly run desired commands. Eg:

# By default this script does nothing.

/root/script.sh
(or)
sudo service deluged restart

exit 0


https://zhuanlan.zhihu.com/p/80305764

