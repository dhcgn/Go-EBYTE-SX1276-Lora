# Go EBYTE SX1276 Lora

![](docs/overview.jpg)

I orderd a "EBYTE Lora UART SX1276 868 MHz 100 mW E32-868T20D Radio Transmitter Receiver Module" and want to uses this on my Windows PC and my Rapsberry Pi. For this purpose I created this repository.

![](docs/ttlconverter.PNG)

On Windows I use a "CP2102 USB to TTL converter HW-598 for 3.3V and 5V" to be ablt to plug these modules on a usb hub.

This repository is in an early state, the feature set is limited but will be expand step by step.

- Receiver Module: https://www.ebyte.com/en/product-view-news.aspx?id=132
- Spec: https://www.ebyte.com/en/downpdf.aspx?id=132 (or see [docs/E32-868T20D_Usermanual_EN_v1.8.pdf](docs/E32-868T20D_Usermanual_EN_v1.8.pdf))

## Gettings started

```
PS> go build -o bin.exe; .\bin.exe COM4
```

## Feature Set

- `[X]` Run on Windows
- `[ ]` Run on Linux (Raspberry)
- `[X]` Read settings from the LoRa Receiver Module
- `[ ]` Write settings from the LoRa Receiver Module
- `[X]` Send Data
- `[ ]` Receive Data

## LoRa Receiver Module Parameter

![](docs/parameter.PNG)

## PIN Design

![](docs/pin.jpg)
![](docs/pin_notes.PNG)
