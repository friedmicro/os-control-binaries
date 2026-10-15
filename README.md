# os-remote-control

## Context on why this tooling exists
This is fairly bespoke and esoteric to my very specific needs. You likely don't need or want this, but I've been asked questions about it enough that I feel I need to upload the code for it. And perhaps someone will get some use from it\modify it into something they use.

For context my primary system is a thin client running on a box with a passed through GPU, USB PCIe card, and network card. It also sits behind a double NAT and serves as the jump point as well as the CI\build server for personal projects of mine across my network. Lastly I make heavy use of very latency sensitive software (a target window of 2-3ms for encoding\decoding\transmitting) for game streaming (Sunshine\Moonlight).

My remote box includes a 3 OS install: Windows, Linux (Arch), Proxmox. I regularly bootup whatever I need on this box and either use it directly (such as with VR) or I stream it back (games, software compilation, AI workloads).

Unfortunately the tooling required for remoting orchestrating a 3 OS system is pretty much nonexistant so if you want something like this you pretty much just have to build it yourself.

## Why?
The thin client is running as a VM because realistically on a system that is used for productivity applications, programming, software builds, and streaming games back to it I don't need to allocate an entire box. Also if I want to keep latency low as well as high quality on the game streaming image I need to have a reasonably powerful GPU for the decoding as the resolution I use being 1440p at 120hz is very bandwidth\computationally heavy.

The double NAT is mostly a security measure and this machine being the jump server is largely because there are very few VMs on the host node. Additionally if you want to do something like development over ssh when traveling you don't want a jump server in there as your dev tooling may not properly handle it.

The CI\CD on here is largely because I test builds on the box and make heavy use of containers so in the interest of minimizing docker image downloads (or just general OCI images) I would rather build on the same box. This saves me bandwidth and saves me time. Additionally, considering the number of people who are actively using my CI is very low I'm not likely to hit issues here and complex binaries for things like a desktop application (or a highly complex microservice) can be easily compiled within my k8s cluster. CI builds are usually under 1 minute compute on the thin client so why bother with a seperate dev box?

As for application\game streaming my partner and I's PCs are both located in our VR place spaces which require a large amount of space, but there spaces do not have comfortable room for desks nor would we want to cut into the space available. We have a den where our thin clients are located, but that room traps heat exceptionally well so just deploying additional high performance machines simply is not an option. By the nature of our situation we cannot run all of our software on our local systems which is why we both make heavy use of application streaming. This is also why my application launcher Athena was created (https://github.com/friedmicro/athena).

As for the remote box having 3 OS installs. The Windows bare metal is because sometimes you want to run a game which simply cannot run in a compatibility layer and requires bare metal access due to an anticheat or interaction with hardware. The Linux (Arch) bare metal is the default OS that the system boots into and it's mostly used for remotely streaming games as well as running computationally heavy builds during development tasks. The proxmox install is to allow me to cleanly deploy VMs on the box, hardware passthrough with a single GPU is significantly simpler when you are not dealing with the OS also trying to run a desktop environment on and you likely do not want to mix the required kernel parameters\module blacklists so keeping it on its own OS install is necesssary.

The complexity exists because it is truly the best solution given the circumstances. If abstracted in software it becomes a nonissue which is why this code is here.

## How it works
The project will build a set of binaries and copy them to `/usr/local/bin`. The reason seperate binaries are built is there is a lot of muscle memory and the binaries are used in some specific tooling\scripts that I would like to avoid updating if I can.

If starting the remote box a webhook request will be sent to turn on a smart plug (Wake On Lan is not supported properly on that motherboard and even if it was it wouldn't be super reliable). That PC is set to always power on when power is restored. We then attempt to ssh into the box to check /etc/os-release file. Assuming this is on Arch then we know the remote box started as expected. The next steps depend upon if we are staying in this OS or if we are booting into another. If we are staying the binary exits here.

If we are booting into Windows a script is ran on the Arch OS install that runs the efibootmgr command to set the respective OS to start and we then reboot the box. Additionally we will create a file called `/tmp/windows-bare-metal` so additional tooling of mine can be aware we are in Windows (currently this is a daemon to control my monitors\screen based upon what a physical KVM is connected to). And lastly we also run a webhook in homeassistant that toggles on all my base stations in the room. Exiting the Windows OS instance simply deletes the file, reboots back into Arch, powers off the OS, and then flips the smart plug.

If we are booting into Proxmox then we run a efibootmgr command to set proxmox as the next boot and then reboot the box into it. We then poll until it actually stays Debian (though Proxmox is also checked in the event this changes down the road). At this point any VMs you want started such as a Windows VM or an AI VM will be started depending on what binary was ran. To power down we spin down the VM, reboot back into Arch, and same process as if we were powering down from Windows.

If we are starting an application stream in addition to starting the Linux box we also spin down the CI (so that new work will not get run while the box is unreachable) and then we switch our network interface online from a virtual nic to a passed through nic that skips the double NAT. The close-nic-port and open-nic-port are automatically granted permissions to run on my system so this is done without interupting any applications using them.

## Note on hardcoded values
Any webhooks in here only work locally so I have no issue with them being in here because they're not exploitable. If you want to replicate this setup or take inspiration from it you'll see a group of configurable constants at the top of `internal/osctl/osctl.go`.

## Build

```sh
make            # build all binaries into ./bin/
make test       # run unit tests
make vet        # run go vet
make install    # build + copy binaries to $(PREFIX)/bin (default /usr/local/bin)
make uninstall  # remove the binaries from $(PREFIX)/bin
```
