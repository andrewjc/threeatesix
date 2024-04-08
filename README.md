# ThreeAteSix - A 386 Emulator

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.20-blue.svg)](https://golang.org/dl/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/andrewjc/threeatesix/actions)
[![Coverage Status](https://img.shields.io/badge/coverage-80%25-green.svg)](https://codecov.io/gh/andrewjc/threeatesix)
[![GitHub contributors](https://img.shields.io/github/contributors/andrewjc/threeatesix.svg)](https://github.com/andrewjc/threeatesix/graphs/contributors)

ThreeAteSix is a 386 emulator written in Go, aiming to accurately emulate the functionality of a 386 CPU and associated hardware components. The project is a work in progress and is intended for educational purposes. 

## Table of Contents

- [Getting Started](#getting-started)
- [Features](#features)
- [Project Structure](#project-structure)
- [CPU Emulation Progress](#cpu-emulation-progress)
- [Configuration](#configuration)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Getting Started

To get started with ThreeAteSix, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/andrewjc/threeatesix.git
   ```

2. Build the project:
   ```
   go build
   ```

3. Prepare a BIOS image file named `bios.bin` and place it in the project directory.

4. Run the emulator:
   ```
   ./threeatesix
   ```
   
5. Use the `monitor` package to interact with the emulator and monitor its state:
    ```
   devices/monitor/hardware_monitor.go:
   device.logCpuInstructions = true
    ```
   Setting logCpuInstructions to true will display all instructions executed by the CPU.

## Sample Output
```
2024/04/08 11:05:11 PS/2 Keyboard connected
2024/04/08 11:05:11 PRIMARY PROCESSOR entered REAL MODE
2024/04/08 11:05:11 MATH CO PROCESSOR entered REAL MODE
2024/04/08 11:05:11 BIOS POST: 0x01 - Register test about to start
2024/04/08 11:05:11 BIOS POST: 0x02 - Register test has passed
2024/04/08 11:05:11 BIOS POST: 0x03 - ROM BIOS checksum test passed
2024/04/08 11:05:11 BIOS POST: 0x04 - Passed keyboard controller test with and w
ithout mouse
2024/04/08 11:05:11 BIOS POST: 0x05 - Chipset initialized... DMA and interrupt c
ontroller disabled
2024/04/08 11:05:11 CMOS RAM WRITE: 0x00, 0x00
2024/04/08 11:05:11 BIOS POST: 0x06 - Video system disabled and the system timer
 checks OK
2024/04/08 11:05:11 BIOS POST: 0x07 - 8265 programmable interval timer initializ
ed
2024/04/08 11:05:11 BIOS POST: 0x08 - Delta counter channel 2 initialized
2024/04/08 11:05:12 BIOS POST: 0x09 - Delta counter channel 1 initialized
2024/04/08 11:05:12 BIOS POST: 0x0a - Delta counter channel 0 initialized
2024/04/08 11:05:12 PS2 Port Output Data: [0x0055]                       
2024/04/08 11:05:12 BIOS POST: 0x0b - Refresh started                    
2024/04/08 11:05:12 BIOS POST: 0x0c - System timer started               
2024/04/08 11:05:12 BIOS POST: 0x0d - Refresh check OK                   
2024/04/08 11:05:12 CMOS RAM WRITE: 0x0d, 0x0c                           
2024/04/08 11:05:12 CMOS RAM WRITE: 0x00, 0x0c                           
2024/04/08 11:05:12 BIOS POST: 0x0e - Refresh period check OK            
2024/04/08 11:05:12 CMOS RAM WRITE: 0x0c, 0xdc
2024/04/08 11:05:12 BIOS POST: 0x0f - System timer check OK
2024/04/08 11:05:12 BIOS POST: 0x10 - Ready to start 64KB base memory test      
2024/04/08 11:05:12 CMOS RAM WRITE: 0xdc, 0x00
2024/04/08 11:05:12 BIOS POST: 0x11 - Address line test OK

```

## Project Structure

The project is structured into the following main packages:

- \`common\`: Contains common constants, functions, and utilities used throughout the project.
- \`devices\`: Contains the emulated hardware devices.
    - \`bus\`: Implements the device bus for communication between devices.
    - \`cga\`: Emulates the Motorola 6845 CGA video controller.
    - \`hid/kb\`: Emulates the PS/2 keyboard.
    - \`intel8086\`: Emulates the Intel 80386 CPU and 80287 Math Coprocessor.
    - \`intel82335\`: Emulates the Intel 82335 High Integration Interface Device.
    - \`intel8259a\`: Emulates the Intel 8259A Programmable Interrupt Controller.
    - \`intel82C54\`: Emulates the Intel 82C54 Programmable Interval Timer.
    - \`io\`: Implements the I/O Port Access Controller.
    - \`memmap\`: Implements the Memory Access Controller.
    - \`monitor\`: Provides debugging and monitoring capabilities.
    - \`ps2\`: Emulates the PS/2 Controller.
- \`pc\`: Defines the main \`PersonalComputer\` struct representing the emulated PC.
 
## Features

- Emulation of Intel 80386 CPU
- Emulation of Intel 80287 Math Coprocessor
- Emulation of Intel 82335 High Integration Interface Device
- Emulation of Intel 8259A Programmable Interrupt Controller
- Emulation of Intel 82C54 Programmable Interval Timer
- Emulation of PS/2 Controller and Keyboard
- Emulation of CGA (Motorola 6845) Video Controller
- Memory and I/O Port Access Controller
- Configurable RAM and ROM sizes
- Loadable BIOS image
- Debugging and monitoring capabilities

## CPU Emulation Progress

The ThreeAteSix emulator currently implements the following features of the Intel 80386 CPU:

- Real mode execution
- Protected mode execution
- Segmentation and paging
- Interrupt handling
- Basic arithmetic and logical instructions
- Memory access instructions
- Control flow instructions (jumps, calls, returns)
- Flags and condition codes

The emulator aims to provide accurate emulation of the 80386 CPU, including its registers, instruction set, and memory model. While significant progress has been made, some advanced features and corner cases may not yet be fully implemented.

## Configuration

The emulator can be configured by modifying the constants in the `pc` package:

- `BiosFilename`: Specifies the name of the BIOS image file to be loaded.
- `MaxRAMBytes`: Specifies the amount of RAM installed in the virtual machine.

## Testing

The project includes unit tests for various components. To run the tests, use the following command:

```
go test ./...
```

## Contributing


Contributions to ThreeAteSix are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

The development of ThreeAteSix was inspired by the desire to learn more about computer architecture and emulation. Special thanks to the authors of the reference materials and resources used during the development process.
