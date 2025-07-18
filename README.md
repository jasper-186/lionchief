

# LionChief Controller

Script for controlling a LionChief train via Bluetooth written in go-lang

Reimplimentation of:
https://github.com/Property404/lionchief-controller

PS. Also credit to Property404 for this README which I shamelessly copied and edited to fit

## About

LionChief trains can be controlled via Bluetooth(BLE, not classic) from a smart phone
using the LionChief app. This library uses that functionality to allow you control the train via a generic Bluetooth adapter. 

## Troubleshooting

* Note that the speaker for at least some (if not all) Lionel trains is in the
tender; if it's not hooked up, the train will not make any sounds
* Some values like pitch are signed. Negative values are represented in 2's
compliment

## Usage

Demo usage can be found in the examples directory. Make sure to change the MAC address, or train name depending on script
