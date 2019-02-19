![SuBFraMe - Logo](https://www.fuchstim.de/subframe/res/img/logo/logo-1-nobg_medium.png)
# How to Contribute to the SuBFraMe Project
Here, You will find a guide on how to contribute to the SuBFraMe Project. These Information might be subject to change, so it might be useful to check back every once and again!

If you wish to contribute anything not neccessarily related to programming (e.g. Fixing Typos, Documentation, Testing, Building a Website, etc.), please do not hesitate to seek contact! You can do so via [EMail](mailto:subframe@gp.ftim.eu) or via [our Discord-Server](https://discord.gg/HwTebxs)

#### Table of Contents
1. [Setting up Your Environment](#setting-up-your-environment)
2. [Folder Structures](#folder-structures)

#### Setting Up Your Environment
SuBFraMe is written in Go. Therefore, you need to install and set up the [latest Version of the Go-Compiler](https://golang.org/dl/), your local Go-Workspace ([see here for more information](https://golang.org/doc/install), and a supported IDE or Editor of your choice. 

SuBFraMe requires mattn's go-sqlite3-Library. A guide on how to install this library can be found [here](http://mattn.github.io/go-sqlite3/). The Library relies on cgo and requires a working gcc installation. For Linux, you can easily install it from most repos. For Windows, [TDM-GCC](http://tdm-gcc.tdragon.net/download) works, so far without any problems.

That's pretty much it! You should now be able to locally compile and run SuBFraMe. Please don't hesitate to report any Issues or uncertainties!

**If you want to actively support and contribute to the SuBFraMe - Project,** please consider joining [our Discord](https://discord.gg/HwTebxs). This is not required, but makes communication easier and helps to resolve questions, uncertainties or problems. Discord is free to use, can be used completely in-browser and is substancially faster than #Slack.

#### Folder Structures
Dividing Code into multiple files and categorizing them by function makes especially extensive projects much more clear. 
Therefore, the SuBFraMe code is divided and labeled accordingly. Please try to adapt and stick to the basic principles already used in the existing code.
Everything StorageNode- or CoordinatorNode-Related goes into the `server` directory. Inside lie all application components, categorized by function.
Data Structs, which can be used by both the server applications and possible future client applications, should reside in the `structs` folder.

TL;DR: Please keep the code you write as organized as possible!
