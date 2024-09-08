## Overview
This program manages files and directories on Linux systems, allowing users to create hardlinks between files. It provides the ability to detect missing files based on their inode number across multiple directories. The application can be used in both simulation mode and real execution mode.

## Installation
The program requires Go to be installed on your system. If you don't have it yet, please install it first by following the official instructions at [Go website](https://go.dev/). Once installed, clone this repository using Git:

```bash
git clone https://github.com/sunrize/d2dhl.git
cd d2dhl
```

Then, build the binary:

```bash
go build
```

After building, you will find the executable file `d2dhl` in the root directory of the project.

## Usage
To use the program, specify the required arguments as follows:

```bash
./d2dhl -src_dirs="path/to/source1,path/to/source2" -dst_dirs="path/to/destination1,path/to/destination2" -dest="path/to/main/destination"
```

Here are some usage examples:

1. Detect missing files and create hard links without console output:
   ```bash
   ./d2dhl -src_dirs="path/to/source1,path/to/source2" -dst_dirs="path/to/destination1,path/to/destination2" -dest="path/to/main/destination"
   ```

2. Detect missing files and create hard links with console output:
   ```bash
   ./d2dhl -src_dirs="path/to/source1,path/to/source2" -dst_dirs="path/to/destination1,path/to/destination2" -dest="path/to/main/destination" -output
   ```

3. Dry run to see what would happen without performing any action:
   ```bash
   ./d2dhl -src_dirs="path/to/source1,path/to/source2" -dst_dirs="path/to/destination1,path/to/destination2" -dest="path/to/main/destination" -dry
   ```

Each parameter must be separated by commas.

## Command Line Arguments

| Argument        | Description                      |
|-----------------|----------------------------------|
| `-src_dirs`     | Comma-separated list of source directories                                |
| `-dst_dirs`     | Comma-separated list of destination directories                            |
| `-dest`         | Main destination directory for links                                       |
| `-output`       | Whether to output collected inodes                       |
| `-dry`          | Whether to perform actions or just simulate them                           |
