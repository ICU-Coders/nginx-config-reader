
 <p align="center" >
   <img src="https://raw.githubusercontent.com/ICU-Coders/IconLib/master/icon.jpg" alt="ICU-Coders" title="ICU-Coders">
 </p>

# nginx-config-reader(ncr)

nginx configuration file reader tool, used to quickly preview all nginx configuration files under the current service

![MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat)

Results Show
![example](./example.jpg)

### How To Use
The following two options are available：

1. Download the compiled executable file [Release](https://github.com/ICU-Coders/nginx-config-reader/releases)

```shell
chmod o+x ./ncr-linux
./ncr-linux
```

2. Local operation
```shell
go run ./main.go
```
### Optional parameters
```
-i input Specify the address of the parsed file 
-l log Output the location of the nginx configuration log
-s sort server/listen/uri/dir, default `listen`
-m match Matching character filter
-h help Help content
```

## MIT License

Copyright (c) 2022 ICU-Coders

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.