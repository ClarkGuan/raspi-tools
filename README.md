# raspi-tools

这是一个辅助工具集，用于帮助我在树梅派平台上进行 Go 或 Rust 开发。工具集包括：

- `cargo-rpi`

  Rust 构建工具 cargo 的插件，在交叉编译时帮助用户构建和运行树梅派平台的可执行程序。该工具依赖后面将要介绍的 `rpirun`。

- `rpienv`

  该工具用于帮助 Go 编译器找到树梅派的交叉编译环境。

- `rpirun`

  辅助运行交叉编译产物的工具，免于频繁的 scp/rsync 等低效操作，提高调试/测试速度。

## 下载 & 编译安装

```shell
git clone https://github.com/ClarkGuan/raspi-tools.git
cd raspi-tools
go install ./cmd/...
```

前提：安装 Go 环境。

## 使用说明

这些工具依赖树梅派的交叉编译工具链，因此需要提前安装他们：

```shell
sudo apt-get install gcc-aarch64-linux-gnu \
g++-aarch64-linux-gnu \
binutils-aarch64-linux-gnu
```

### 注意事项

- `rpirun` 默认将交叉编译产物拷贝到 `/home/<用户名>/` 目录下，这里没有想过要有什么扩展性，就是写死的路径
- 目前对于带有动态依赖库的可执行文件如何拷贝和运行没有进行适配
- 如果树梅派设备的 ip、user、password 不是默认的，可以使用环境变量 `RASPI_IP`、`RASPI_USER` 和 `RASPI_PASSWORD` 进行设置

### 使用 cargo-rpi 在树梅派设备上运行 Rust helloworld 程序

1. 创建工程

   ```shell
   cargo new hello --bin && cd hello
   ```

2. 确保与树梅派设备在同一个子网中并设置了它的 IP（使用 `RASPI_IP` 环境变量）。当然，如果用户名和密码已经修改过的话也要通过环境变量一并设置。

3. 使用下面命令编译并在设备上运行：

   ```shell
   cargo rpi r --release
      Compiling hello v0.1.0 (/home/xxx/hello)
       Finished release [optimized] target(s) in 0.42s
        Running `rpirun target/aarch64-unknown-linux-musl/release/hello`
   root@192.168.x.x:
   ==================================
   Hello, world!
   ```

### 在树梅派设备上跑 Rust 单元测试

1. 编写测试代码：

    ```rust
    #[cfg(test)]
    mod tests {
        use lazy_static::lazy_static;
    
        lazy_static! {
            static ref G_EXPECTED_NAMES: Vec<&'static str> = vec!["Albert", "Black", "Clark"];
        }
    
        #[test]
        fn test_sort() {
            let mut names = vec!["Clark", "Black", "Albert"];
            names.sort();
            assert_eq!(&names, &*G_EXPECTED_NAMES);
        }
    }
    ```

2. 运行：

    ```shell
    cargo rpi t
       Compiling lazy_static v1.4.0
       Compiling untitled v0.1.0 (/home/xxx/untitled)
        Finished test [unoptimized + debuginfo] target(s) in 0.67s
         Running unittests src/main.rs (target/aarch64-unknown-linux-musl/debug/deps/untitled-6a4cc396201a6d06)
    root@192.168.x.x: 
    ==================================
    
    running 1 test
    test tests::test_sort ... ok
    
    test result: ok. 1 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.00s
    ```

### 使用 rpienv 和 rpirun 在树梅派设备上运行 Go helloworld 程序

1. 创建工程

   ```shell
   mkdir hello && cd hello
   go mod init hello
   ```

2. 创建 main.go，这里使用 cgo 访问 C 的 IO 库函数（测试 cgo 调用 C 函数）：

    ```go
    package main
    
    // #include <stdio.h>
    import "C"
    
    func main() {
        C.puts(C.CString("Hello world!!!"))
        C.fflush(C.stdout)
    }
    ```

3. 编译并运行：

   ```shell
   rpienv go run -exec rpirun .
   root@192.168.x.x:
   ==================================
   Hello world!!!
   ```

### 在树梅派设备上跑 Go 的单元测试或性能测试

1. 创建测试函数：

    ```go
    package main_test
    
    import (
        "reflect"
        "sort"
        "testing"
    )
    
    var names = []string{"Clark", "Blank", "Albert"}
    
    func TestSort(t *testing.T) {
        sort.Sort(sort.StringSlice(names))
        if !reflect.DeepEqual(names, []string{"Albert", "Blank", "Clark"}) {
            t.Fail()
        }
    }
    ```

2. 运行下面命令：

    ```shell
    rpienv go test -v -exec rpirun 
    root@192.168.x.x: 
    ==================================
    === RUN   TestSort
    --- PASS: TestSort (0.00s)
    PASS
    ok  	hello	1.296s
    ```
   
