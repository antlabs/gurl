package scheduler

import "github.com/antlabs/cronex"

// Run 使用 cronex 根据表达式 expr 定时调用 fn。
// 调用方通过 stop 通道控制何时结束调度循环。
func Run(expr string, fn func(), stop <-chan struct{}) error {
	c := cronex.New()
	_, err := c.AddFunc(expr, fn)
	if err != nil {
		return err
	}

	c.Start()
	<-stop
	c.Stop()
	return nil
}
