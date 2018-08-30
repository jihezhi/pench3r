// 这个直接参考了标准答案(尴尬了,没想到直接搜索到原题,下次可需要注意了)

package main

import (	
	"fmt"
)

// 建立from和to channel之间的关系以及方向
func cc(from,to chan int) {
    to <- 1 + <- from
}

func main() {
    // 实际应用中都是使用nobuffer-channel，也提倡使用，很少使用buffer-channel
    var dstorg chan int = make(chan int)
    // 保存用于出结果的channel，相当于尾部
    var src chan int = dstorg
    var dst chan int = dstorg
	  for i := 0;i<10000;i++ {
        // 创建channel用完不需要保存，因为通过goroutine进行管道之间管理的维持
        // 最终只需要这一串channel的头和尾即可
	      src = make(chan int)
		    go cc(src, dst)
		    dst = src
	  }
    // 在头部的channel中放入数据
    go func(c chan int) {c <- 4}(src)
    // 直接就在尾部的channel中获取到经过多成管道处理的结果
	  fmt.Println(<-dstorg)
}
